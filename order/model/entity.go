package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPending        = "PENDING"
	OrderStockReserved  = "STOCK_RESERVED"
	OrderPaymentPending = "PAYMENT_PENDING"
	OrderConfirmed      = "CONFIRMED"
	OrderCancelled      = "CANCELLED"
)

type Order struct {
	ID             uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID   `gorm:"not null;uniqueIndex:idx_user_idempotency"`
	Status         OrderStatus `gorm:"type:varchar(30);not null;default:'PENDING'"`
	Total          float64     `gorm:"type:decimal(12,2);not null"`
	IdempotencyKey string      `gorm:"not null;uniqueIndex:idx_user_idempotency"`

	Items []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"type:decimal(12,2);not null"`

	CreatedAt time.Time
}
