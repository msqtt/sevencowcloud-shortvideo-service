package token

import "time"

// TokenMaker is an interface for signing up token.
type TokenMaker interface {
	// CreateToken returns token generated from given nickname, payload and an error, if any.
	CreateToken(nickname string, duration time.Duration) (string, *Payload, error)
	// ValidToken valids the given token and returns the payload and an error, if any.
	ValidToken(token string) (*Payload, error)
}
