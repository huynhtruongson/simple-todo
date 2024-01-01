package task_repo

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/field"
	mocks "github.com/huynhtruongson/simple-todo/mocks/lib"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRepo_CreateTask(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	row := mocks.NewRow(t)
	mockTask := task_entity.Task{
		Title:  field.NewString("title"),
		UserID: field.NewInt(1),
		Status: 1,
	}
	expectQuery := `INSERT INTO tasks (user_id,title,status,description) VALUES ($1,$2,$3,$4) RETURNING task_id`
	tests := []struct {
		name         string
		task         task_entity.Task
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call insert query exactly",
			task: mockTask,
			mock: func() {
				db.EXPECT().QueryRow(ctx, expectQuery, utils.GenerateMockArguments(4)...).Once().Return(row)
				row.EXPECT().Scan(mock.Anything).Once().Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			userId, err := repo.CreateTask(ctx, db, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectUserID, userId)
		})
	}
}

func TestUserRepo_UpdateTask(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	mockTask := task_entity.Task{
		TaskID: field.NewInt(1),
		Title:  field.NewString("title"),
		UserID: field.NewInt(1),
		Status: 1,
	}
	expectQuery := `UPDATE tasks SET user_id=$1, title=$2, status=$3, description=$4,updated_at = now() WHERE task_id = $5`
	tests := []struct {
		name         string
		task         task_entity.Task
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call update query exactly",
			task: mockTask,
			mock: func() {
				db.EXPECT().Exec(ctx, expectQuery, utils.GenerateMockArguments(5)...).Once().Return(pgconn.CommandTag{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			err := repo.UpdateTask(ctx, db, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestTaskRepo_GetTasksByIds(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	rows := mocks.NewRows(t)
	expectQuery := `SELECT task_id,user_id,title,status,description FROM tasks WHERE task_id = ANY($1) AND user_id = $2 AND deleted_at IS NULL`
	tests := []struct {
		name         string
		userID       int
		taskIDs      []int
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name:    "should call get query exactly",
			userID:  1,
			taskIDs: []int{1},
			mock: func() {
				db.EXPECT().Query(ctx, expectQuery, []int{1}, 1).Once().Return(rows, nil)
				rows.EXPECT().Next().Once().Return(true)
				rows.EXPECT().Scan(utils.GenerateMockArguments(5)...).Once().Return(nil)
				rows.EXPECT().Next().Once().Return(false)
				rows.EXPECT().Close().Once().Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			_, err := repo.GetTasksByIds(ctx, db, tt.userID, tt.taskIDs)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestTaskRepo_DeleteTask(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	expectQuery := `UPDATE tasks SET deleted_at = now() WHERE task_id = $1`
	tests := []struct {
		name         string
		id           int
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call delete query exactly",
			id:   1,
			mock: func() {
				db.EXPECT().Exec(ctx, expectQuery, 1).Once().Return(pgconn.CommandTag{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			err := repo.DeleteTask(ctx, db, tt.id)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestTaskRepo_CountTask(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	row := mocks.NewRows(t)
	expectQuery := `SELECT count(task_id) FROM tasks WHERE user_id = $1 AND deleted_at IS NULL`
	tests := []struct {
		name         string
		id           int
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call count query exactly",
			id:   1,
			mock: func() {
				db.EXPECT().QueryRow(ctx, expectQuery, 1).Once().Return(row, nil)
				row.EXPECT().Scan(mock.Anything).Once().Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			_, err := repo.CountTask(ctx, db, 1)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestTaskRepo_GetTasksWithFilter(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	rows := mocks.NewRows(t)
	expectQuery := `SELECT task_id,user_id,title,status,description FROM tasks WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at LIMIT $2 OFFSET $3`
	tests := []struct {
		name         string
		id           int
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call count query exactly",
			id:   1,
			mock: func() {
				db.EXPECT().Query(ctx, expectQuery, 1, 10, 0).Once().Return(rows, nil)
				rows.EXPECT().Next().Once().Return(true)
				rows.EXPECT().Scan(utils.GenerateMockArguments(5)...).Once().Return(nil)
				rows.EXPECT().Next().Once().Return(false)
				rows.EXPECT().Close().Once().Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			_, err := repo.GetTasksWithFilter(ctx, db, 1, 10, 0)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
