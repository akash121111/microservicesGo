package model

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EventType  string    `gorm:"not null"`
	RoutingKey string    `gorm:"not null"`
	Payload    []byte    `gorm:"type:jsonb;not null"`
	Published  bool      `gorm:"not null;index"`

	CreatedAt   time.Time
	PublishedAt *time.Time
}
