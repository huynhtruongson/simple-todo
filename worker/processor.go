package worker

import (
	"context"

	"github.com/huynhtruongson/simple-todo/lib"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/rs/zerolog/log"

	"github.com/hibiken/asynq"
)

const (
	CriticalQueue = "critical"
	DefaultQueue  = "default"
	LowQueue      = "low"
)

type userRepo interface {
	GetUsersByUsername(context.Context, lib.QueryExecer, string) ([]user_entity.User, error)
}
type RedisTaskProcessor struct {
	server *asynq.Server
	db     lib.DB
	userRepo
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db lib.DB, userRepo userRepo) *RedisTaskProcessor {
	errHandlerFunc := func(ctx context.Context, task *asynq.Task, err error) {
		log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("process task failed")
	}
	return &RedisTaskProcessor{
		server: asynq.NewServer(redisOpt, asynq.Config{
			Queues: map[string]int{
				CriticalQueue: 6,
				DefaultQueue:  3,
				LowQueue:      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(errHandlerFunc),
			Logger:       NewLogger(),
		}),
		db:       db,
		userRepo: userRepo,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)
	// Define more

	return p.server.Start(mux)
}
