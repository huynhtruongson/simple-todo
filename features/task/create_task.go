package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	common_resp "github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/features/common"
	"github.com/huynhtruongson/simple-todo/field"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	task_repo "github.com/huynhtruongson/simple-todo/services/task/repository"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/stretchr/testify/assert"
)

func (s *suite) createTask(ctx context.Context, payloadType string) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	randStr := utils.RandomString(10)
	taskPayload := task_entity.Task{
		Title:       field.NewString(fmt.Sprintf("title %s", randStr)),
		Status:      task_entity.TaskStatusDoing,
		Description: field.NewString(fmt.Sprintf("description %s", randStr)),
	}
	switch payloadType {
	case "empty title":
		taskPayload.Title = field.NewNullString()
	case "invalid status":
		taskPayload.Status = 3
	}

	jsonPayload, err := json.Marshal(taskPayload)
	if err != nil {
		return ctx, fmt.Errorf("json marshal task payload failed,%w\n", err)
	}
	acToken, ok := stepState.Get(common.TokenKey).(string)
	if !ok {
		return ctx, fmt.Errorf("retrieve token string failed")
	}
	statusCode, resBody, err := common.MakeHttpRequest(fmt.Sprintf("%s/%s", s.BaseURL, "v1/task/create"), http.MethodPost, jsonPayload, acToken)
	if err != nil {
		return ctx, err
	}

	if payloadType == "valid information" {
		err = common.AssertExpectedAndActual(assert.Equal, http.StatusOK, statusCode, "response status code mismatch")
		if err != nil {
			return ctx, err
		}
		var respStruct common_resp.SuccessResponse
		err = json.Unmarshal(resBody, &respStruct)
		if err != nil {
			return ctx, fmt.Errorf("json unmarshal response data failed,%w\n", err)
		}
		taskPayload.TaskID = field.NewInt(int(respStruct.Data.(float64)))
		payload, ok := stepState.Get(common.TokenPayloadKey).(token.TokenPayload)
		if !ok {
			return ctx, fmt.Errorf("retrieve token payload failed")
		}
		taskPayload.UserID = field.NewInt(payload.UserID)
		stepState.Set(common.RequestPayloadKey, taskPayload)
		return common.StepStateToContext(ctx, stepState), nil
	} else {
		var errorStruct common_resp.AppError
		err = json.Unmarshal(resBody, &errorStruct)
		if err != nil {
			return ctx, fmt.Errorf("json unmarshal response data failed,%w\n", err)
		}
		stepState.Set(common.ResponseKey, errorStruct)
		return common.StepStateToContext(ctx, stepState), nil
	}
}

func (s *suite) assertCreatedTask(ctx context.Context) (context.Context, error) {
	stepState := common.StepStateFromContext(ctx)
	payload, ok := stepState.Get(common.RequestPayloadKey).(task_entity.Task)
	if !ok {
		return ctx, fmt.Errorf("retrieve task payload failed")
	}
	taskRepo := task_repo.NewTaskRepo()
	tasksDB, err := taskRepo.GetTasksByIds(ctx, s.DB, payload.UserID.Int(), []int{payload.TaskID.Int()})
	if err != nil {
		return ctx, fmt.Errorf("get task by ids error,%w\n", err)
	}

	if len(tasksDB) != 1 {
		return ctx, fmt.Errorf("task not found in db,user_id:%d\n", payload.UserID.Int())
	}
	if err := common.AssertExpectedAndActual(assert.Equal, payload.Title.String(), tasksDB[0].Title.String()); err != nil {
		return ctx, err
	}
	if err := common.AssertExpectedAndActual(assert.Equal, payload.Status, tasksDB[0].Status); err != nil {
		return ctx, err
	}
	if err := common.AssertExpectedAndActual(assert.Equal, payload.UserID.Int(), tasksDB[0].UserID.Int()); err != nil {
		return ctx, err
	}
	if err := common.AssertExpectedAndActual(assert.Equal, payload.Description.String(), tasksDB[0].Description.String()); err != nil {
		return ctx, err
	}
	return ctx, nil
}
