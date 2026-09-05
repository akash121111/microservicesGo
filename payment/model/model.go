package model

import (
	"time"

	"github.com/google/uuid"
)

type ProcessedEvent struct {
	EventID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProcessedAt time.Time
}
