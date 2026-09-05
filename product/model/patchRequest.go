package model

type ProductPatchRequest struct {
	Description *string  `json:"description"`
	Stock       *int64   `json:"stock"`
	Price       *float64 `json:"price"`
	Status      *string  `json:"status"`
}

type ProductUpdateStocksRequest struct {
	Stock *int64 `json:"stock"`
}
