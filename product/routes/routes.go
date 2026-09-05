package routes

import (
	"product/handler"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, productHandler *handler.ProductHandler) {

	publicRoutes := r.Group("")

	publicRoutes.POST("/product", productHandler.CreateProduct)
	publicRoutes.GET("/product", productHandler.GetAllProduct)
	publicRoutes.GET("/product/:id", productHandler.GetProductByID)
	publicRoutes.PATCH("/product/:id", productHandler.PatchProduct)
	publicRoutes.DELETE("/product/:id", productHandler.DeleteProductByID)
	publicRoutes.PATCH("/product/stock/reserve", productHandler.ReserveProductStocks)
	publicRoutes.PATCH("/product/stock/release", productHandler.RelaeseProductStocks)

}
