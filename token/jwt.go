package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/huynhtruongson/simple-todo/common"
)

type JWTMaker struct {
	secretKey string
}
type JWTPayload struct {
	TokenPayload
	jwt.RegisteredClaims
}

func NewJWTMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < 32 {
		return nil, common.InvalidKeySizeErrorMessage
	}
	jwtMaker := &JWTMaker{
		secretKey,
	}
	return jwtMaker, nil
}

func NewJWTPayload(payload TokenPayload) JWTPayload {
	return JWTPayload{
		TokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        payload.ID.String(),
		},
	}
}

func (m *JWTMaker) CreateToken(userID int, duration time.Duration, tokenType TokenType) (string, TokenPayload, error) {
	payload, err := NewTokenPayload(userID, duration, tokenType)
	if err != nil {
		return "", payload, err
	}
	jwtClaim := NewJWTPayload(payload)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	signedToken, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", payload, err
	}
	return signedToken, payload, nil
}

func (m *JWTMaker) VerifyToken(token string) (TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)
	if err != nil {
		return TokenPayload{}, err
	}
	jwtPayload, ok := jwtToken.Claims.(*JWTPayload)
	if !ok {
		return TokenPayload{}, common.InvalidTokenErrorMessage
	}
	return jwtPayload.TokenPayload, nil
}
