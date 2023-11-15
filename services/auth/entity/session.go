package auth_entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID    uuid.UUID
	UserID       int
	RefreshToken string
	UserAgent    string
	ClientIP     string
	IsBlocked    bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func NewSession(id uuid.UUID, userId int, rfToken string, expiresAt time.Time) Session {
	return Session{
		SessionID:    id,
		UserID:       userId,
		RefreshToken: rfToken,
		ExpiresAt:    expiresAt,
	}
}

func (t Session) TableName() string {
	return "sessions"
}

func (s *Session) FieldMap() ([]string, []interface{}) {
	return []string{
			"session_id",
			"user_id",
			"refresh_token",
			"user_agent",
			"client_ip",
			"is_blocked",
			"expires_at",
		}, []interface{}{
			&s.SessionID,
			&s.UserID,
			&s.RefreshToken,
			&s.UserAgent,
			&s.ClientIP,
			&s.IsBlocked,
			&s.ExpiresAt,
		}
}
