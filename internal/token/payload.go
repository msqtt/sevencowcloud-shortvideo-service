package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Payload inside a token.
type Payload struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	IssuedAt  time.Time `json:"issued_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

// NewPayload creates a payload.
func NewPayload(nickname string, duration time.Duration) (*Payload, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot new uuid: %w", err)
	}
	return &Payload{
		ID: u,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}
