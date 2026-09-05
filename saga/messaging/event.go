package messaging

import "github.com/google/uuid"

type OrderCreatedEvent struct {
	Event       string           `json:"event"`
	EventID     uuid.UUID        `json:"eventId"`
	OrderID     uuid.UUID        `json:"orderId"`
	UserID      uuid.UUID        `json:"userId"`
	TotalAmount float64          `json:"totalAmount"`
	Items       []OrderEventItem `json:"items"`
}
type OrderEventItem struct {
	ProductID uuid.UUID `json:"productID"`
	Quantity  int64     `json:"quantity"`
}

type ReserveStockCommand struct {
	EventID     uuid.UUID        `json:"eventId"`
	Event       string           `json:"event"`
	OrderID     uuid.UUID        `json:"orderId"`
	UserID      uuid.UUID        `json:"userId"`
	Items       []OrderEventItem `json:"items"`
	TotalAmount float64          `json:"totalAmount"`
}

type ReserveStockEvent struct {
	EventID     uuid.UUID        `json:"eventId"`
	Event       string           `json:"event"`
	OrderID     uuid.UUID        `json:"orderId"`
	UserID      uuid.UUID        `json:"userId"`
	Items       []OrderEventItem `json:"items"`
	TotalAmount float64          `json:"totalAmount"`
}

type ReserveStockFailedEvent struct {
	EventID uuid.UUID        `json:"eventId"`
	Event   string           `json:"event"`
	OrderID uuid.UUID        `json:"orderId"`
	Items   []OrderEventItem `json:"items"`
	Reason  string           `json:"reason"`
}

type PaymentProcessCommand struct {
	EventID     uuid.UUID `json:"eventId"`
	Event       string    `json:"event"`
	OrderID     uuid.UUID `json:"orderId"`
	UserID      uuid.UUID `json:"userId"`
	TotalAmount float64   `json:"totalAmount"`
}

type PaymentSuccessEvent struct {
	EventID uuid.UUID `json:"eventId"`
	Event   string    `json:"event"`
	OrderID uuid.UUID `json:"orderId"`
	UserID  uuid.UUID `json:"userId"`
	Amount  float64   `json:"amount"`
}

type PaymentFailedEvent struct {
	EventID uuid.UUID `json:"eventId"`
	Event   string    `json:"event"`
	OrderID uuid.UUID `json:"orderId"`
	UserID  uuid.UUID `json:"userId"`
	Reason  string    `json:"reason"`
}
