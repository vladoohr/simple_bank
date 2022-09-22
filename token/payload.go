package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Define token validation errors
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data for the token
type Payload struct {
	ID       uuid.UUID `json: "id"`
	Username string    `json: "username"`
	IssuedAt time.Time `json: "issued_at"`
	ExpireAt time.Time `json: "expire_at"`
}

// NewPayload returns new payload object
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if token is expired or not
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpireAt) {
		return ErrExpiredToken
	}

	return nil
}
