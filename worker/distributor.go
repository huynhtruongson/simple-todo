package worker

import "github.com/hibiken/asynq"

// type TaskDistributor interface {
// 	DistributeTaskSend
// }

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewTaskDistributor(opt asynq.RedisClientOpt) *RedisTaskDistributor {
	return &RedisTaskDistributor{
		client: asynq.NewClient(opt),
	}
}
