package task_service

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/field"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteTaskBiz_DeleteTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name      string
		userID    int
		taskID    int
		mock      func(prop *MockServiceProp)
		expectErr error
	}{
		{
			name:   "should Delete task successfully",
			taskID: 1,
			userID: 1,
			mock: func(prop *MockServiceProp) {
				prop.TaskRepo.EXPECT().GetTasksByIds(ctx, prop.DB, 1, []int{1}).Once().Return([]task_entity.Task{{TaskID: field.NewInt(1)}}, nil)
				prop.DB.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(prop.TX, nil)
				prop.TaskRepo.EXPECT().DeleteTask(ctx, prop.TX, 1).Once().Return(nil)
				prop.TX.EXPECT().Commit(ctx).Once().Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockTaskService(t)
			tt.mock(prop)
			err := sv.DeleteTask(ctx, tt.userID, tt.taskID)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
