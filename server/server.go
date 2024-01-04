package server

import (
	"context"
	"net"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/huynhtruongson/simple-todo/interceptor"
	"github.com/huynhtruongson/simple-todo/lib"
	"github.com/huynhtruongson/simple-todo/pb"
	auth_service "github.com/huynhtruongson/simple-todo/services/auth"
	auth_port "github.com/huynhtruongson/simple-todo/services/auth/port"
	task_service "github.com/huynhtruongson/simple-todo/services/task"
	task_port "github.com/huynhtruongson/simple-todo/services/task/port"
	user_service "github.com/huynhtruongson/simple-todo/services/user"
	user_port "github.com/huynhtruongson/simple-todo/services/user/port"
	user_repo "github.com/huynhtruongson/simple-todo/services/user/repository"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/huynhtruongson/simple-todo/worker"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	db         lib.DB
	config     utils.Config
	tokenMaker token.TokenMaker
	taskSv     *task_service.TaskService
	userSv     *user_service.UserService
	authSv     *auth_service.AuthService
}

func NewServer(opts ...Option) Server {
	server := Server{}
	for _, opt := range opts {
		opt(&server)
	}
	return server
}

func (s Server) RunMigration(path string) {
	migration, err := migrate.New(path, s.config.DBAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize migrate instance")
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migration up")
	}
	log.Info().Msg("run migration successfully!")
}

func (s Server) RunGinServer() {
	r := gin.Default()
	authAPI := auth_port.NewAuthAPIService(s.authSv)
	userAPI := user_port.NewUserAPIService(s.userSv)
	taskAPI := task_port.NewTaskAPIService(s.taskSv)
	v1 := r.Group("/v1")
	v1.Use(interceptor.LoggingMiddleware)
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authAPI.Login)
			auth.POST("/renew-token", authAPI.RenewToken)
		}

		user := v1.Group("/user")
		{
			user.POST("/create", userAPI.CreateUser)
		}

		task := v1.Group("/task")
		task.Use(interceptor.AuthMiddleware(s.tokenMaker))
		{
			task.GET("/list", taskAPI.ListTask)
			task.POST("/create", taskAPI.CreateTask)
			task.PUT("/update/:id", taskAPI.UpdateTask)
			task.PUT("/delete/:id", taskAPI.DeleteTask)
		}
	}

	r.Run(s.config.ApiServerPort)
}

func (s Server) RunGRPCServer() {
	authInterceptor := interceptor.NewAuthInterceptor(s.tokenMaker)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.UnaryServerLoggingInterceptor(),
		interceptor.UnaryServerAuthInterceptor(authInterceptor),
	))
	userGRPCService := user_port.NewUserGRPCService(s.userSv)
	pb.RegisterUserServiceServer(grpcServer, userGRPCService)
	authGRPCService := auth_port.NewAuthGRPCService(s.authSv)
	pb.RegisterAuthServiceServer(grpcServer, authGRPCService)
	taskGRPCService := task_port.NewTaskGRPCService(s.taskSv)
	pb.RegisterTaskServiceServer(grpcServer, taskGRPCService)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", s.config.GRPCServerPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
	log.Info().Msgf("Listening and serving GRPC on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot serve grpc server")
	}
}

func (s Server) RunGatewayServer() {
	grpcMux := runtime.NewServeMux()
	authGRPCService := auth_port.NewAuthGRPCService(s.authSv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterAuthServiceHandlerServer(ctx, grpcMux, authGRPCService)
	if err != nil {
		log.Fatal().Err(err).Msg("can not register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", s.config.ApiServerPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
	log.Info().Msgf("Listening and serving HTTP Gateway on %s", listener.Addr().String())
	handler := interceptor.GatewayLoggingMiddleware(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot serve http gateway server")
	}
}

func (s Server) RunTaskWorker() {
	redisOpt := asynq.RedisClientOpt{
		Addr: s.config.RedisAddress,
	}
	server := worker.NewRedisTaskProcessor(
		redisOpt,
		s.db,
		user_repo.NewUserRepo(),
	)

	err := server.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run task worker")
	}
	log.Info().Msg("starting task worker")
}
