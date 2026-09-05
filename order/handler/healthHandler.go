package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

func (h *HealthHandler) Ready(ctx *gin.Context) {
	db, err := h.db.DB()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"sucess":  false,
			"message": "not ready DB",
			"errors":  err,
		})
		return
	}
	if err := db.PingContext(ctx.Request.Context()); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"sucess":  false,
			"message": "DB Unavailable",
			"errors":  err,
		})
		return
	}
	if err := h.redis.Ping(ctx.Request.Context()).Err(); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"sucess":  false,
			"message": "Redis Unavailable",
			"errors":  err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": "ready",
	})
}

func (h *HealthHandler) Health(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": "ok",
	})
}
