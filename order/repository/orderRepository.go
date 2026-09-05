package repository

import (
	"context"
	"order/apperror"
	"order/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) WithTx(tx *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: tx,
	}
}

func (r *OrderRepository) GetAllOrder() []model.OrderModel {
	var Orders []model.Order

	result := r.db.Find(&Orders)
	if result.Error != nil {
		return []model.OrderModel{}
	}
	return toOrders(Orders)

}

func (r *OrderRepository) GetMyAllOrderWithItem(ctx context.Context, userId string) []model.OrderDataWithItems {
	var Orders []model.Order

	result := r.db.WithContext(ctx).Preload("Items").Where("user_id=?", userId).Find(&Orders)
	if result.Error != nil {
		return []model.OrderDataWithItems{}
	}
	orderData := make([]model.OrderDataWithItems, 0, len(Orders))

	for _, order := range Orders {
		orderData = append(orderData, model.OrderDataWithItems{
			Order: toOrder(order),
			Items: toOrderItems(order.Items),
		})
	}

	return orderData

}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id string) (*model.OrderDataWithItems, error) {
	var Order model.Order

	result := r.db.WithContext(ctx).Preload("Items").Where("id=?", id).First(&Order)
	if result.Error != nil {
		return nil, apperror.ErrOrderNotFound
	}
	data := model.OrderDataWithItems{
		Order: toOrder(Order),
		Items: toOrderItems(Order.Items),
	}
	return &data, nil
}

func (r *OrderRepository) CreateOrder(ctx context.Context, OrderData model.OrderModel) (*model.OrderModel, error) {
	orderEntity := model.Order{
		UserID:         OrderData.UserID,
		Total:          OrderData.Total,
		Status:         OrderData.Status,
		IdempotencyKey: OrderData.IdempotencyKey,
	}
	result := r.db.WithContext(ctx).Create(&orderEntity)

	if result.Error != nil {
		return nil, result.Error
	}
	order := toOrder(orderEntity)
	return &order, nil

}

func (r *OrderRepository) CreateOrderItem(ctx context.Context, OrderItemData []model.OrderItemModel) (*[]model.OrderItemModel, error) {
	entities := make([]model.OrderItem, 0, len(OrderItemData))
	for _, item := range OrderItemData {
		entities = append(entities, model.OrderItem{
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}
	result := r.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		return nil, result.Error
	}
	order := toOrderItems(entities)
	return &order, nil

}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, id string, status string) (*model.OrderModel, error) {
	var OrderEntity model.Order
	result := r.db.WithContext(ctx).Where("id=?", id).First(&OrderEntity)

	if result.Error != nil {
		return nil, apperror.ErrOrderNotFound
	}

	OrderEntity.Status = model.OrderStatus(status)

	result = r.db.WithContext(ctx).Save(&OrderEntity)
	if result.Error != nil {
		return nil, result.Error
	}

	data := toOrder(OrderEntity)
	return &data, nil
}

func (r *OrderRepository) DeleteOrderByID(ctx context.Context, id string) bool {
	result := r.db.WithContext(ctx).Delete(&model.Order{}, id)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}

func toOrder(Order model.Order) model.OrderModel {
	return model.OrderModel{
		ID:             Order.ID,
		UserID:         Order.UserID,
		Total:          Order.Total,
		Status:         Order.Status,
		IdempotencyKey: Order.IdempotencyKey,
	}
}

func toOrders(Orders []model.Order) []model.OrderModel {
	OrdersData := make([]model.OrderModel, 0, len(Orders))
	for _, Order := range Orders {
		OrdersData = append(OrdersData, toOrder(Order))
	}
	return OrdersData
}

func toOrderItem(Order model.OrderItem) model.OrderItemModel {
	return model.OrderItemModel{
		ID:        Order.ID,
		OrderID:   Order.OrderID,
		ProductID: Order.ProductID,
		Price:     Order.Price,
		Quantity:  int64(Order.Quantity),
	}
}

func toOrderItems(Orders []model.OrderItem) []model.OrderItemModel {
	OrdersData := make([]model.OrderItemModel, 0, len(Orders))
	for _, Order := range Orders {
		OrdersData = append(OrdersData, toOrderItem(Order))
	}
	return OrdersData
}
