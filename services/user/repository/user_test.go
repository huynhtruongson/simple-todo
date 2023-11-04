package user_repo

import (
	"context"
	"testing"

	mock_db "github.com/sondev/todo-list/mock"
	user_entity "github.com/sondev/todo-list/services/user/entity"
	"github.com/sondev/todo-list/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRepo_CreateUser(t *testing.T) {
	ctx := context.Background()
	db := &mock_db.MockTx{}
	row := &mock_db.MockRow{}
	mockUser := user_entity.User{
		FullName: "fullname",
		Username: "username",
		Password: "123123",
	}
	expectQuery := `INSERT INTO users (fullname,username,password) VALUES ($1,$2,$3) RETURNING user_id`
	tests := []struct {
		name         string
		user         user_entity.User
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call insert query exactly",
			user: mockUser,
			mock: func() {
				db.On("QueryRow", utils.GenerateMockArguments(3, ctx, expectQuery)...).Once().Return(row)
				row.On("Scan", mock.Anything).Once().Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepo{}
			tt.mock()
			userId, err := repo.CreateUser(ctx, db, tt.user)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectUserID, userId)
		})
	}
}

func TestUserRepo_GetUsersByUsername(t *testing.T) {
	ctx := context.Background()
	db := &mock_db.MockTx{}
	rows := &mock_db.MockRows{}
	expectQuery := `SELECT user_id,fullname,username,password FROM users WHERE username = $1 AND deleted_at IS NULL`
	tests := []struct {
		name         string
		username     string
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name:     "should call get query exactly",
			username: "username",
			mock: func() {
				db.On("Query", ctx, expectQuery, "username").Once().Return(rows, nil)
				rows.On("Next").Once().Return(true)
				rows.On("Scan", utils.GenerateMockArguments(4)...).Once().Return(nil)
				rows.On("Next").Once().Return(false)
				rows.On("Close").Once().Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepo{}
			tt.mock()
			_, err := repo.GetUsersByUsername(ctx, db, tt.username)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestUserRepo_GetUsersByUserIds(t *testing.T) {
	ctx := context.Background()
	db := &mock_db.MockTx{}
	rows := &mock_db.MockRows{}
	expectQuery := `SELECT user_id,fullname,username,password FROM users WHERE user_id = ANY($1) AND deleted_at IS NULL`
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
				rows.On("Scan", utils.GenerateMockArguments(4)...).Once().Return(nil)
				rows.On("Next").Once().Return(false)
				rows.On("Close").Once().Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepo{}
			tt.mock()
			_, err := repo.GetUsersByUserIds(ctx, db, tt.ids)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
