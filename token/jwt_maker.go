package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeyLen = 32

// JWTMaker holds the secret key for signing the token
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLen {
		return nil, fmt.Errorf("Invalid key size, must be at least %d characters", minSecretKeyLen)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates new token for a given username and password
func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(jwtMaker.secretKey))

	return tokenString, payload, err
}

// VerifyToken verifies the token and returns the token payload
func (jwtMaker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(jwtMaker.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if errors.Is(verr.Inner, ErrExpiredToken) {
				return nil, ErrExpiredToken
			}
		}

		return nil, err
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
