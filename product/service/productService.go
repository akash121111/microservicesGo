package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"product/apperror"
	"product/middleware"
	"product/model"
	"product/repository"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// type ProductRepository interface {
// 	CreateProduct(productData model.Product) model.ProductS
// 	ProductNameExist(productName string) bool
// 	GetAllProduct() []model.Product
// 	GetProductByID(id string) (*model.Product, error)
// 	PatchProduct(id string, product model.ProductPatchRequest) (*model.Product, error)
// 	DeleteProductByID(id string) bool
// 	WithTx(tx *gorm.DB) *ProductRepository
// }

type ProductService struct {
	productRepo *repository.ProductRepository
	redis       *redis.Client
	db          *gorm.DB
	logger      *slog.Logger
}

func NewProductService(productRepo *repository.ProductRepository, redis *redis.Client, db *gorm.DB, logger *slog.Logger) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		redis:       redis,
		db:          db,
		logger:      logger,
	}
}

func (s *ProductService) GetAllProduct() []model.ProductModel {
	return s.productRepo.GetAllProduct()
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*model.ProductModel, error) {
	s.logger.Info(
		"Get product By Id",
		"productId", id,
		"correlationId", middleware.GetCorrelationID(ctx),
	)
	// 1. Check Redis
	key := "product:" + id
	cashe, err := s.redis.Get(ctx, key).Result()

	// 2. If cache hit → return
	if err == nil {
		log.Println("hit")
		var product model.ProductModel
		err := json.Unmarshal([]byte(cashe), &product)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}
	// 3. Cache miss → repository
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrProductNotFound
	}
	// 4. Store result in Redis
	data, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}
	err = s.redis.Set(ctx, key, data, 10*time.Minute).Err()
	// 5. Return
	if err != nil {
		log.Println("redis set failed", err)
	}

	return product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, productData model.ProductModel) (model.ProductModel, error) {
	exist := s.productRepo.ProductNameExist(ctx, productData.Name)
	if exist {
		return model.ProductModel{}, apperror.ErrProductAlreadyExist
	}
	if productData.Status == "" {
		productData.Status = string(model.ProductActive)
	}

	return s.productRepo.CreateProduct(ctx, productData), nil
}

func (s *ProductService) PatchProduct(ctx context.Context, id string, product model.ProductPatchRequest) (*model.ProductModel, error) {

	result, err := s.productRepo.PatchProduct(ctx, id, product)

	if err != nil {
		return nil, apperror.ErrProductNotFound
	}
	key := "product:" + id
	if casheErr := s.redis.Del(ctx, key).Err(); casheErr != nil {
		log.Println("failed to delete cashe")
	}

	return result, nil
}
func (s *ProductService) ReserveProductStock(ctx context.Context, items []model.OrderEventItem) error {

	s.logger.Info(
		"Reserve stock",
		"correlationId", middleware.GetCorrelationID(ctx),
	)
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := s.productRepo.WithTx(tx)
		for _, item := range items {
			err := repo.ReserveProductStock(ctx, item.ProductID.String(), item.Quantity)
			if err != nil {
				return err
			}
			key := "product:" + item.ProductID.String()
			if casheErr := s.redis.Del(ctx, key).Err(); casheErr != nil {
				log.Println("failed to delete cashe")
			}

		}

		return nil
	})
}
func (s *ProductService) ReleaseProductStock(ctx context.Context, items []model.OrderEventItem) error {
	fmt.Println(items)
	for _, item := range items {
		err := s.productRepo.ReleaseProductStock(ctx, item.ProductID.String(), item.Quantity)
		if err != nil {
			return err
		}
		key := "product:" + item.ProductID.String()
		if casheErr := s.redis.Del(ctx, key).Err(); casheErr != nil {
			log.Println("failed to delete cashe")
		}

	}
	return nil
}
func (s *ProductService) DeleteProductByID(ctx context.Context, id string) (bool, error) {
	result := s.productRepo.DeleteProductByID(ctx, id)

	if !result {
		return false, apperror.ErrProductNotFound
	}
	return true, nil
}
