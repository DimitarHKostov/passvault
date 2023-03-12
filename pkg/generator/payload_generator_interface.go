package generator

import (
	"passvault/pkg/types"
	"time"
)

type PayloadGeneratorInterface interface {
	GeneratePayload(duration time.Duration) (*types.Payload, error)
}
