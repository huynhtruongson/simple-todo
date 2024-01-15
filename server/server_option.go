package server

import (
	"context"

	docs "github.com/huynhtruongson/simple-todo/docs"
	auth_service "github.com/huynhtruongson/simple-todo/services/auth"
	auth_repo "github.com/huynhtruongson/simple-todo/services/auth/repository"
	task_service "github.com/huynhtruongson/simple-todo/services/task"
	task_repo "github.com/huynhtruongson/simple-todo/services/task/repository"
	user_service "github.com/huynhtruongson/simple-todo/services/user"
	user_repo "github.com/huynhtruongson/simple-todo/services/user/repository"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/huynhtruongson/simple-todo/worker"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Option func(*Server)

func WithConfig(cfg utils.Config) Option {
	return func(s *Server) {
		s.config = cfg
	}
}

func WithDB() Option {
	return func(s *Server) {
		db, err := pgxpool.New(context.Background(), s.config.DBAddress)

		if err != nil {
			log.Fatal().Err(err).Msg("unable to connect to database")
		}
		s.db = db
	}
}

func WithTokenMaker() Option {
	return func(s *Server) {
		tokenMaker, err := token.NewPasetoMaker(s.config.TokenKey)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot initialize token maker")
		}
		s.tokenMaker = tokenMaker
	}
}

func WithUserService() Option {
	return func(s *Server) {
		redisOpt := asynq.RedisClientOpt{
			Addr: s.config.RedisAddress,
		}
		taskDistributor := worker.NewTaskDistributor(redisOpt)
		userRepo := user_repo.NewUserRepo()
		userSv := user_service.NewUserService(s.db, taskDistributor, userRepo)
		s.userSv = userSv
	}
}

func WithTaskService() Option {
	return func(s *Server) {
		userRepo := user_repo.NewUserRepo()
		taskRepo := task_repo.NewTaskRepo()
		taskSv := task_service.NewTaskService(s.db, taskRepo, userRepo)
		s.taskSv = taskSv
	}
}

func WithAuthService() Option {
	return func(s *Server) {
		userRepo := user_repo.NewUserRepo()
		sessionRepo := auth_repo.NewSessionRepo()
		authService := auth_service.NewAuthService(s.db, s.tokenMaker, userRepo, sessionRepo)
		s.authSv = authService
	}
}

func WithSwaggerDoc() Option {
	return func(s *Server) {
		docs.SwaggerInfo.Title = "Simple Todo API"
		docs.SwaggerInfo.Description = "This is API documentation of Simple Todo server."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = s.config.ApiServerPort
		docs.SwaggerInfo.BasePath = "/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	}
}
