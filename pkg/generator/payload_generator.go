package generator

import (
	"passvault/pkg/types"
	"time"

	"github.com/google/uuid"
)

type PayloadGenerator struct{}

func (pg *PayloadGenerator) GeneratePayload(duration time.Duration) (*types.Payload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &types.Payload{
		Uuid:      uuid,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}
