package generator

import (
	"fmt"
	"passvault/pkg/log"
	"passvault/pkg/types"
	"time"

	"github.com/google/uuid"
)

type PayloadGenerator struct {
	logManager log.LogManagerInterface
}

func NewPayloadGenerator(logManager log.LogManagerInterface) *PayloadGenerator {
	payloadGenerator := &PayloadGenerator{logManager: logManager}

	return payloadGenerator
}

func (pg *PayloadGenerator) GeneratePayload(duration time.Duration) (*types.Payload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		pg.logManager.LogError(fmt.Sprintf(uuidGenerationFailMessage, err))
		return nil, err
	}

	payload := &types.Payload{
		Uuid:      uuid,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	pg.logManager.LogDebug(successfulPayloadGenerationMessage)

	return payload, nil
}
