package model

import (
	"github.com/google/uuid"
)

type ProductModel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Stock       int64   `json:"stock" binding:"required"`
	Price       float64 `json:"price"  binding:"required"`
	Status      string  `json:"status"`
}

type OrderEventItem struct {
	ProductID uuid.UUID `json:"productID"`
	Quantity  int64     `json:"quantity"`
}

type ReleaseStockRequest struct {
	Items []OrderEventItem `json:"items"`
}
