package token

import (
	"time"

	"github.com/google/uuid"
	"github.com/huynhtruongson/simple-todo/common"
)

type TokenMaker interface {
	CreateToken(userID int, duration time.Duration, tokenType TokenType) (string, TokenPayload, error)
	VerifyToken(token string) (TokenPayload, error)
}

type TokenType int32

const (
	AccessToken TokenType = iota
	RefreshToken
)

type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int       `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expired_at"`
	Type      TokenType `json:"token_type"`
}

const (
	AccessTokenDuration  = time.Minute * 15
	RefreshTokenDuration = time.Hour * 24
)

func NewTokenPayload(userID int, duration time.Duration, tokenType TokenType) (TokenPayload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return TokenPayload{}, err
	}
	now := time.Now().UTC()

	payload := TokenPayload{
		ID:        uuid,
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
		Type:      tokenType,
	}
	return payload, nil
}

func (p TokenPayload) Valid() error {
	now := time.Now().UTC()
	if now.After(p.ExpiresAt.UTC()) {
		return common.ExpiredTokenErrorMessage
	}
	return nil
}
