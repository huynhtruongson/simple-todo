package auth_repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	mocks "github.com/huynhtruongson/simple-todo/mocks/lib"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_CreateSession(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	mockSession := auth_entity.Session{}
	expectQuery := `INSERT INTO sessions (session_id,user_id,refresh_token,user_agent,client_ip,is_blocked,expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	tests := []struct {
		name      string
		session   auth_entity.Session
		mock      func()
		expectErr error
	}{
		{
			name:    "should call insert query exactly",
			session: mockSession,
			mock: func() {
				db.EXPECT().Exec(ctx, expectQuery, utils.GenerateMockArguments(7)...).Once().Return(pgconn.CommandTag{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SessionRepo{}
			tt.mock()
			err := repo.CreateSession(ctx, db, tt.session)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}

func TestTaskRepo_GetSessionByIds(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDB(t)
	rows := mocks.NewRows(t)
	randomID, err := uuid.NewRandom()
	assert.NoError(t, err)
	expectQuery := `SELECT session_id,user_id,refresh_token,user_agent,client_ip,is_blocked,expires_at FROM sessions WHERE session_id = ANY($1)`
	tests := []struct {
		name         string
		ids          uuid.UUIDs
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should call get query exactly",
			ids:  uuid.UUIDs{randomID},
			mock: func() {
				db.EXPECT().Query(ctx, expectQuery, uuid.UUIDs{randomID}).Once().Return(rows, nil)
				rows.EXPECT().Next().Once().Return(true)
				rows.EXPECT().Scan(utils.GenerateMockArguments(7)...).Once().Return(nil)
				rows.EXPECT().Next().Once().Return(false)
				rows.EXPECT().Close().Once().Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SessionRepo{}
			tt.mock()
			_, err := repo.GetSessionByIds(ctx, db, tt.ids)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
