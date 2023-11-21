package auth_service

import (
	"context"

	"github.com/huynhtruongson/simple-todo/lib"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetUsersByUsername(ctx context.Context, db lib.QueryExecer, username string) ([]user_entity.User, error)
}
type SessionRepo interface {
	CreateSession(ctx context.Context, db lib.QueryExecer, session auth_entity.Session) error
	GetSessionByIds(ctx context.Context, db lib.QueryExecer, ids uuid.UUIDs) ([]auth_entity.Session, error)
}

type AuthService struct {
	DB         lib.DB
	TokenMaker token.TokenMaker
	UserRepo
	SessionRepo
}

func NewAuthService(db lib.DB, tokenMaker token.TokenMaker, userRepo UserRepo, sessionRepo SessionRepo) *AuthService {
	return &AuthService{
		DB:          db,
		TokenMaker:  tokenMaker,
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}
