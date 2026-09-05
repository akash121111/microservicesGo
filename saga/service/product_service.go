package service

import (
	"context"
	"log"
	"saga/client"
	"saga/model"
)

type SagaProductService struct {
	productClient *client.ProductClient
}

func NewProductSagaService(productClient *client.ProductClient) *SagaProductService {
	return &SagaProductService{
		productClient: productClient,
	}
}

func (p *SagaProductService) ReleaseStock(ctx context.Context, items []model.OrderEventItem) error {
	log.Println(items)
	return p.productClient.ReleaseStock(ctx, items)
}
