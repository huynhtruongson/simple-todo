package task_repo

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	mock_db "github.com/sondev/todo-list/mock"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	"github.com/sondev/todo-list/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRepo_CreateTask(t *testing.T) {
	ctx := context.Background()
	db := &mock_db.MockTx{}
	row := &mock_db.MockRow{}
	mockTask := task_entity.Task{
		Title:  "title",
		UserID: 1,
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
				db.On("QueryRow", utils.GenerateMockArguments(4, ctx, expectQuery)...).Once().Return(row)
				row.On("Scan", mock.Anything).Once().Return(nil)
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
	db := &mock_db.MockTx{}
	mockTask := task_entity.Task{
		TaskID: 1,
		Title:  "title",
		UserID: 1,
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
				db.On("Exec", utils.GenerateMockArguments(5, ctx, expectQuery)...).Once().Return(pgconn.CommandTag{}, nil)
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
	db := &mock_db.MockTx{}
	rows := &mock_db.MockRows{}
	expectQuery := `SELECT task_id,user_id,title,status,description FROM tasks WHERE task_id = ANY($1) AND deleted_at IS NULL`
	tests := []struct {
		name         string
		ids          []int
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call get query exactly",
			ids:  []int{1},
			mock: func() {
				db.On("Query", ctx, expectQuery, []int{1}).Once().Return(rows, nil)
				rows.On("Next").Once().Return(true)
				rows.On("Scan", utils.GenerateMockArguments(5)...).Once().Return(nil)
				rows.On("Next").Once().Return(false)
				rows.On("Close").Once().Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			_, err := repo.GetTasksByIds(ctx, db, tt.ids)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestTaskRepo_DeleteTask(t *testing.T) {
	ctx := context.Background()
	db := &mock_db.MockTx{}
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
				db.On("Exec", ctx, expectQuery, 1).Once().Return(pgconn.CommandTag{}, nil)
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
	db := &mock_db.MockTx{}
	row := &mock_db.MockRow{}
	expectQuery := `SELECT count(task_id) FROM tasks WHERE deleted_at IS NULL`
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
				db.On("QueryRow", ctx, expectQuery).Once().Return(row, nil)
				row.On("Scan", mock.Anything).Once().Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			_, err := repo.CountTask(ctx, db)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestTaskRepo_GetTasksWithFilter(t *testing.T) {
	ctx := context.Background()
	db := &mock_db.MockTx{}
	rows := &mock_db.MockRows{}
	expectQuery := `SELECT task_id,user_id,title,status,description FROM tasks WHERE deleted_at IS NULL ORDER BY created_at LIMIT $1 OFFSET $2`
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
				db.On("Query", ctx, expectQuery, 10, 0).Once().Return(rows, nil)
				rows.On("Next").Once().Return(true)
				rows.On("Scan", utils.GenerateMockArguments(5)...).Once().Return(nil)
				rows.On("Next").Once().Return(false)
				rows.On("Close").Once().Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TaskRepo{}
			tt.mock()
			_, err := repo.GetTasksWithFilter(ctx, db, 10, 0)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
