package repository

import (
	"context"
	"fmt"
	"product/apperror"
	"product/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) WithTx(tx *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: tx,
	}
}

func (r *ProductRepository) GetAllProduct() []model.ProductModel {
	var products []model.Product

	result := r.db.Find(&products)
	if result.Error != nil {
		return []model.ProductModel{}
	}
	return toProducts(products)

}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*model.ProductModel, error) {
	var product model.Product
	result := r.db.WithContext(ctx).Where("id=?", id).First(&product)
	if result.Error != nil {
		return nil, apperror.ErrProductNotFound
	}
	data := toProduct(product)
	return &data, nil
}

func (r *ProductRepository) GetProductForUpdate(ctx context.Context, id string) (*model.ProductModel, error) {
	var product model.Product
	println("Trying to lock product:", id)
	result := r.db.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("id=?", id).First(&product)
	if result.Error != nil {
		return nil, apperror.ErrProductNotFound
	}
	println("Locked product:", id)
	println("Stokes:", product.Stock)
	data := toProduct(product)
	return &data, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, productData model.ProductModel) model.ProductModel {
	product := model.Product{
		Name:        productData.Name,
		Description: productData.Description,
		Stock:       int(productData.Stock),
		Price:       productData.Price,
		Status:      model.ProductStatus(productData.Status),
	}
	result := r.db.WithContext(ctx).Create(&product)

	if result.Error != nil {
		return model.ProductModel{}
	}
	return toProduct(product)

}
func (r *ProductRepository) ProductNameExist(ctx context.Context, productName string) bool {
	var count int64

	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("name=?", productName).Count(&count)
	if result.Error != nil {
		return false
	}
	return count > 0
}

func (r *ProductRepository) PatchProduct(ctx context.Context, id string, product model.ProductPatchRequest) (*model.ProductModel, error) {
	var productEntity model.Product
	result := r.db.WithContext(ctx).Where("id=?", id).First(&productEntity)
	fmt.Println(productEntity)
	if result.Error != nil {
		return nil, apperror.ErrProductNotFound
	}
	if product.Status != nil {
		productEntity.Status = model.ProductStatus(*product.Status)
	}
	if product.Price != nil {
		productEntity.Price = *product.Price
	}
	if product.Stock != nil {
		productEntity.Stock = int(*product.Stock)
	}
	if product.Description != nil {
		productEntity.Description = *product.Description
	}

	result = r.db.WithContext(ctx).Save(&productEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	data := toProduct(productEntity)
	return &data, nil
}

func (r *ProductRepository) ReserveProductStock(
	ctx context.Context,
	id string,
	quantity int64,
) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ? and stock>=?", id, quantity).UpdateColumn("stock", gorm.Expr("stock-?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperror.ErrInsuffiecientProduct
	}
	return nil
}
func (r *ProductRepository) ReleaseProductStock(
	ctx context.Context,
	id string,
	quantity int64,
) error {

	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ? ", id).UpdateColumn("stock", gorm.Expr("stock+?", quantity))

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ProductRepository) DeleteProductByID(ctx context.Context, id string) bool {
	result := r.db.WithContext(ctx).Delete(&model.Product{}, id)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}

func toProduct(product model.Product) model.ProductModel {
	return model.ProductModel{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Stock:       int64(product.Stock),
		Price:       product.Price,
		Status:      string(product.Status),
	}
}

func toProducts(products []model.Product) []model.ProductModel {
	productsData := make([]model.ProductModel, 0, len(products))
	for _, product := range products {
		productsData = append(productsData, toProduct(product))
	}
	return productsData
}
