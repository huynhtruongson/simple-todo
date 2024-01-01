package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskSendVerifyEmailPayload struct {
	Username string `json:"username"`
}

const TaskSendVerifyEmail = "task:send_verify_email"

func (d *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *TaskSendVerifyEmailPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("can not marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	taskInfo, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", taskInfo.Type).Str("queue", taskInfo.Queue).Msg("enqueued task")

	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var taskSendVerifyEmailPayload TaskSendVerifyEmailPayload
	err := json.Unmarshal(task.Payload(), &taskSendVerifyEmailPayload)
	if err != nil {
		return fmt.Errorf("can not unmarshal task payload: %w", asynq.SkipRetry)
	}
	// handle logic in Tx if modify data

	users, err := p.userRepo.GetUsersByUsername(ctx, p.db, taskSendVerifyEmailPayload.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if len(users) != 1 {
		return fmt.Errorf("user not found")
	}
	// TODO: handle logic
	log.Info().Str("type", task.Type()).Str("payload", string(task.Payload())).Str("email", users[0].Email.String()).Msg("processed task")
	return nil
}
