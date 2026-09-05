package handler

import (
	"product/apperror"
	"product/model"
	"product/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productHandler *service.ProductService
}

func NewProductHandler(productHandler *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productHandler: productHandler,
	}
}

func (h *ProductHandler) GetAllProduct(ctx *gin.Context) {
	data := h.productHandler.GetAllProduct()
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "All Product Fetched",
		"data":    data,
	})
}

func (h *ProductHandler) GetProductByID(ctx *gin.Context) {

	id := ctx.Param("id")
	data, proErr := h.productHandler.GetProductByID(ctx.Request.Context(), string(id))
	if proErr != nil {
		ctx.Error(apperror.ErrProductNotFound)
		return
	}
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "fetch Product By Id",
		"data":    data,
	})
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var product model.ProductModel
	err := ctx.ShouldBindJSON(&product)
	if err != nil {
		ctx.Error(err)
		return
	}
	data, err := h.productHandler.CreateProduct(ctx.Request.Context(), product)
	if err != nil {
		ctx.Error(apperror.ErrProductAlreadyExist)
		return
	}
	ctx.JSON(201, gin.H{
		"success": true,
		"message": "product Created Successfully",
		"data":    data,
	})
}

func (h *ProductHandler) PatchProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var product model.ProductPatchRequest
	productErr := ctx.ShouldBindJSON(&product)
	if productErr != nil {
		ctx.Error(productErr)
	}
	data, proErr := h.productHandler.PatchProduct(ctx.Request.Context(), id, product)
	if proErr != nil {
		ctx.Error(apperror.ErrProductNotFound)
		return
	}
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "product patch update By Id",
		"data":    data,
	})
}

func (h *ProductHandler) ReserveProductStocks(ctx *gin.Context) {
	var orderItem []model.OrderEventItem
	Err := ctx.ShouldBindJSON(&orderItem)
	if Err != nil {
		ctx.Error(Err)
	}
	Err = h.productHandler.ReserveProductStock(ctx.Request.Context(), orderItem)
	if Err != nil {
		ctx.Error(Err)
		return
	}
	ctx.Status(204)
}

func (h *ProductHandler) RelaeseProductStocks(ctx *gin.Context) {
	var orderItem model.ReleaseStockRequest
	Err := ctx.ShouldBindJSON(&orderItem)
	if Err != nil {
		ctx.Error(Err)
	}
	Err = h.productHandler.ReleaseProductStock(ctx.Request.Context(), orderItem.Items)
	if Err != nil {
		ctx.Error(Err)
		return
	}
	ctx.Status(204)
}

func (h *ProductHandler) DeleteProductByID(ctx *gin.Context) {
	id := ctx.Param("id")
	_, proErr := h.productHandler.DeleteProductByID(ctx.Request.Context(), id)
	if proErr != nil {
		ctx.Error(apperror.ErrProductNotFound)
		return
	}
	ctx.Status(204)
}
