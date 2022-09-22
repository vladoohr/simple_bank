package token

import "time"

type Maker interface {
	// CreateToken creates new token for a given username and password
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken verifies the token and returns the token payload
	VerifyToken(token string) (*Payload, error)
}
