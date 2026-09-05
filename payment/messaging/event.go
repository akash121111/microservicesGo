package messaging

import "github.com/google/uuid"

type PaymentProcessCommand struct {
	Event       string    `json:"event"`
	EventID     uuid.UUID `json:"eventId"`
	OrderID     uuid.UUID `json:"orderId"`
	UserID      uuid.UUID `json:"userId"`
	TotalAmount float64   `json:"totalAmount"`
}

type PaymentSuccessEvent struct {
	Event   string    `json:"event"`
	EventID uuid.UUID `json:"eventId"`
	OrderID uuid.UUID `json:"orderId"`
	UserID  uuid.UUID `json:"userId"`
	Amount  float64   `json:"amount"`
}

type PaymentFailedEvent struct {
	Event   string    `json:"event"`
	EventID uuid.UUID `json:"eventId"`
	OrderID uuid.UUID `json:"orderId"`
	UserID  uuid.UUID `json:"userId"`
	Reason  string    `json:"reason"`
}
