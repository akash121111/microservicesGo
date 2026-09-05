package handler

import (
	"errors"
	"order/model"
	"order/service"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) GetAllMyOrderWithItem(ctx *gin.Context) {
	userIDValue, _ := ctx.Get("userID")
	userIDString, ok := userIDValue.(string)
	if !ok {
		ctx.Error(errors.New("userID is not a string"))
		return
	}

	reqCtx := ctx.Request.Context()
	order := h.orderService.GetMyAllOrderWithItem(reqCtx, userIDString)

	ctx.JSON(200, gin.H{
		"sucess":  true,
		"message": "fetch all my order sucessfully",
		"data":    order,
	})
}

func (h *OrderHandler) GetOrderById(ctx *gin.Context) {
	id := ctx.Param("id")
	reqCtx := ctx.Request.Context()
	// userIDValue, _ := ctx.Get("userID")
	// userIDString, ok := userIDValue.(string)
	// if !ok {
	// 	ctx.Error(errors.New("userID is not a string"))
	// 	return
	// }
	// userID, err := uuid.Parse(userIDString)
	// if err != nil {
	// 	ctx.Error(err)
	// 	return
	// }

	// roleValue, _ := ctx.Get("role")
	// roleString, ok := roleValue.(string)
	// if !ok {
	// 	ctx.Error(errors.New("role is not a string"))
	// 	return
	// }

	order, err := h.orderService.GetOrderById(reqCtx, id)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{
		"sucess":  true,
		"message": "fetch order by id sucessfully",
		"data":    order,
	})
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var orderItem model.OrderRequestBody
	err := ctx.ShouldBindJSON(&orderItem)
	if err != nil {
		ctx.Error(err)
		return
	}
	idempotencyKey := ctx.GetHeader("Idempotency-Key")

	userIDValue, _ := ctx.Get("userID")
	userIDString, ok := userIDValue.(string)
	if !ok {
		ctx.Error(errors.New("userID is not a string"))
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		ctx.Error(err)
		return
	}

	createdOrder, err := h.orderService.CreateOrder(ctx.Request.Context(), idempotencyKey, userID, orderItem)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Order Created Succesfully",
		"data":    createdOrder,
	})
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	orderId := ctx.Param("orderId")
	// userIDValue, _ := ctx.Get("userID")
	// userIDString, ok := userIDValue.(string)
	// if !ok {
	// 	ctx.Error(errors.New("userID is not a string"))
	// 	return
	// }
	// userID, err := uuid.Parse(userIDString)
	// if err != nil {
	// 	ctx.Error(err)
	// 	return
	// }

	// roleValue, _ := ctx.Get("role")
	// roleString, ok := roleValue.(string)
	// if !ok {
	// 	ctx.Error(errors.New("role is not a string"))
	// 	return
	// }
	err := h.orderService.CancelOrder(ctx.Request.Context(), orderId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Order Cancel Succesfully",
		"data":    orderId,
	})

}

func (h *OrderHandler) UpdateOrderStatus(ctx *gin.Context) {
	orderId := ctx.Param("orderId")
	status := ctx.Param("status")
	err := h.orderService.UpdateOrderStatus(ctx.Request.Context(), orderId, status)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Order status update Succesfully",
		"data":    orderId,
	})

}
