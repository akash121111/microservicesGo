package client

import (
	"context"
	"fmt"
	"order/model"
	productpb "proto/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn   *grpc.ClientConn
	client productpb.ProductServiceClient
}

func NewProductClient(address string) (*ProductClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create product grpc client: %w", err)
	}

	return &ProductClient{
		conn:   conn,
		client: productpb.NewProductServiceClient(conn),
	}, nil
}

func (c *ProductClient) GetProduct(
	ctx context.Context,
	id string,
) (*model.ProductModel, error) {
	resp, err := c.client.GetProduct(ctx, &productpb.GetProductRequest{ProductId: id})
	if err != nil {
		return nil, fmt.Errorf("get product via grpc: %w", err)
	}

	return &model.ProductModel{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
		Stock:       resp.Stock,
		Status:      resp.Status,
	}, nil
}

func (c *ProductClient) BookStock(
	ctx context.Context,
	productID string,
	stock int,
) error {
	product, err := c.GetProduct(ctx, productID)
	if err != nil {
		return err
	}
	resp, err := c.client.ReserveStock(ctx)

	return fmt.Errorf("get product via grpc:")
}
