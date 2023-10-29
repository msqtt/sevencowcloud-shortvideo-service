package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// CreateToken implements TokenMaker.
func (pm *PasetoMaker) CreateToken(nickname string, duration time.Duration) (string, *Payload, error) {
	p, err := NewPayload(nickname, duration)
	if err != nil {
		return "", nil, err	
	}
	s, err2 := pm.paseto.Encrypt(pm.symmetricKey, p, nil)
	if err2 != nil {
		return "", nil, err2
	}
	return s, nil, nil
}

// ValidToken implements TokenMaker.
func (pm *PasetoMaker) ValidToken(token string) (*Payload, error) {
	var payload Payload
	err := pm.paseto.Decrypt(token, pm.symmetricKey, &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot descrypt token: %w", err)
	}
	if payload.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("token is expired.")
	}
	return &payload, nil
}

var _ TokenMaker = (*PasetoMaker)(nil)

// NewPasetoMaker returns a PasetoMaker using given symmetricKey.
func NewPasetoMaker(symmetricKey []byte) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}, nil
}
