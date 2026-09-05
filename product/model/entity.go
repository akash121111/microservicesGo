package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductActive   ProductStatus = "ACTIVE"
	ProductInactive ProductStatus = "INACTIVE"
)

type Product struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string        `gorm:"type:varchar(200);not null"`
	Description string        `gorm:"type:text"`
	Price       float64       `gorm:"type:decimal(12,2);not null"`
	Stock       int           `gorm:"not null"`
	Status      ProductStatus `gorm:"type:varchar(30);default:'ACTIVE';not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProcessedEvent struct {
	EventID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProcessedAt time.Time
}
