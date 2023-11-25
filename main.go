package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/huynhtruongson/simple-todo/interceptor"
	"github.com/huynhtruongson/simple-todo/pb"
	auth_service "github.com/huynhtruongson/simple-todo/services/auth"
	auth_port "github.com/huynhtruongson/simple-todo/services/auth/port"
	auth_repo "github.com/huynhtruongson/simple-todo/services/auth/repository"
	task_service "github.com/huynhtruongson/simple-todo/services/task"
	task_port "github.com/huynhtruongson/simple-todo/services/task/port"
	task_repo "github.com/huynhtruongson/simple-todo/services/task/repository"
	user_service "github.com/huynhtruongson/simple-todo/services/user"
	user_port "github.com/huynhtruongson/simple-todo/services/user/port"
	user_repo "github.com/huynhtruongson/simple-todo/services/user/repository"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration", err)
	}
	db, err := pgx.Connect(context.Background(), config.DBAddress)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		log.Fatal("cannot initialize token maker", err)
	}
	r := gin.Default()

	userRepo := user_repo.NewUserRepo()
	sessionRepo := auth_repo.NewSessionRepo()
	taskRepo := task_repo.NewTaskRepo()

	authService := auth_service.NewAuthService(db, tokenMaker, userRepo, sessionRepo)
	taskService := task_service.NewTaskService(db, taskRepo, userRepo)
	userService := user_service.NewUserService(db, userRepo)

	// go runGatewayServer(config, authService)
	go runGRPCServer(config, tokenMaker, userService, authService, taskService)
	runGinServer(r, config, tokenMaker, authService, userService, taskService)
}

func runGinServer(
	r *gin.Engine,
	config utils.Config,
	tokenMaker token.TokenMaker,
	authSv *auth_service.AuthService,
	userSv *user_service.UserService,
	taskSv *task_service.TaskService,
) {
	authAPI := auth_port.NewAuthAPIService(authSv)
	userAPI := user_port.NewUserAPIService(userSv)
	taskAPI := task_port.NewTaskAPIService(taskSv)
	v1 := r.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authAPI.Login)
			auth.POST("/renew-token", authAPI.RenewToken)
		}

		user := v1.Group("/user")
		// user.Use(interceptor.AuthMiddleware(tokenMaker))
		{
			user.POST("/create", userAPI.CreateUser)
		}

		task := v1.Group("/task")
		task.Use(interceptor.AuthMiddleware(tokenMaker))
		{
			task.GET("/list", taskAPI.ListTask)
			task.POST("/create", taskAPI.CreateTask)
			task.PUT("/update/:id", taskAPI.UpdateTask)
			task.PUT("/delete/:id", taskAPI.DeleteTask)
		}
	}

	r.Run(config.ApiServerPort)
}

func runGRPCServer(
	config utils.Config,
	tokenMaker token.TokenMaker,
	userSv *user_service.UserService,
	authSv *auth_service.AuthService,
	taskSv *task_service.TaskService,
) {
	authInterceptor := interceptor.NewAuthInterceptor(tokenMaker)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor(authInterceptor)))
	userGRPCService := user_port.NewUserGRPCService(userSv)
	pb.RegisterUserServiceServer(grpcServer, userGRPCService)
	authGRPCService := auth_port.NewAuthGRPCService(authSv)
	pb.RegisterAuthServiceServer(grpcServer, authGRPCService)
	taskGRPCService := task_port.NewTaskGRPCService(taskSv)
	pb.RegisterTaskServiceServer(grpcServer, taskGRPCService)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerPort)
	if err != nil {
		log.Fatal("cannot create listener", err)
	}
	log.Printf("Listening and serving GRPC on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot serve grpc server", err)
	}
}

func runGatewayServer(
	config utils.Config,
	authSv *auth_service.AuthService,
) {
	grpcMux := runtime.NewServeMux()
	authGRPCService := auth_port.NewAuthGRPCService(authSv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterAuthServiceHandlerServer(ctx, grpcMux, authGRPCService)
	if err != nil {
		log.Fatal("can not register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.ApiServerPort)
	if err != nil {
		log.Fatal("cannot create listener", err)
	}
	log.Printf("Listening and serving HTTP Gateway on %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot serve http gateway server", err)
	}
}
