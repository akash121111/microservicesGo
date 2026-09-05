package grpc

import (
	"context"
	productpb "proto/product"

	"product/model"
	"product/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	productpb.UnimplementedProductServiceServer
	productService *service.ProductService
}

func NewProductServer(productService *service.ProductService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.ProductResponse, error) {
	product, err := s.productService.GetProductByID(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}
	return &productpb.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Status:      product.Status,
	}, nil

}

func (s *ProductServer) ReserveStock(ctx context.Context, req *productpb.ReleaseStockRequest) (*productpb.StockResponse, error) {
	items := make([]model.OrderEventItem, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		productId, err := uuid.Parse(item.GetProductId())
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"invalid product id: %s",
				item.GetProductId(),
			)
		}
		items = append(items, model.OrderEventItem{
			ProductID: productId,
			Quantity:  item.GetQuantity(),
		})
	}
	if err := s.productService.ReserveProductStock(ctx, items); err != nil {
		return nil, err
	}
	return &productpb.StockResponse{
		Success: true,
		Message: "Stock Reserve Succesfully",
	}, nil
}
