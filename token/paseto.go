package token

import (
	"time"

	"github.com/huynhtruongson/simple-todo/common"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	// Must be 32 bytes due to chacha20poly1305 algorithm
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, common.InvalidKeySizeErrorMessage
	}
	tokenMaker := &PasetoMaker{
		paseto:       *paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return tokenMaker, nil
}

func (m *PasetoMaker) CreateToken(userID int, duration time.Duration, tokenType TokenType) (string, TokenPayload, error) {
	payload, err := NewTokenPayload(userID, duration, tokenType)
	if err != nil {
		return "", payload, err
	}

	token, err := m.paseto.Encrypt(m.symmetricKey, payload, nil)
	if err != nil {
		return "", payload, err
	}
	return token, payload, nil
}
func (m *PasetoMaker) VerifyToken(token string) (TokenPayload, error) {
	payload := TokenPayload{}

	if err := m.paseto.Decrypt(token, m.symmetricKey, &payload, nil); err != nil {
		return payload, err
	}
	if err := payload.Valid(); err != nil {
		return payload, err
	}
	return payload, nil
}
