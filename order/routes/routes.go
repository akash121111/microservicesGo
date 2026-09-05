package routes

import (
	"order/handler"
	"order/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, orderHandler *handler.OrderHandler) {

	publicRoutes := r.Group("")

	publicRoutes.GET("/order/:id", orderHandler.GetOrderById)
	publicRoutes.PATCH("/order/cancel/:orderId", orderHandler.CancelOrder)
	publicRoutes.PATCH("/order/:orderId/status/:status", orderHandler.UpdateOrderStatus)

	privateRoutes := r.Group("")
	privateRoutes.Use(middleware.UserMiddleware())
	privateRoutes.POST("/order", orderHandler.CreateOrder)
	privateRoutes.GET("/order/my", orderHandler.GetAllMyOrderWithItem)
}
