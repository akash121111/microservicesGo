package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const correlationIDKey contextKey = "correlationID"

func Correlation() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.GetHeader("X-Correlation-ID")

		if id == "" {
			id = uuid.NewString()
		}

		// Put correlation ID into request context
		reqCtx := context.WithValue(
			ctx.Request.Context(),
			correlationIDKey,
			id,
		)

		// Attach new context to request
		ctx.Request = ctx.Request.WithContext(reqCtx)

		// Return ID to client
		ctx.Header("X-Correlation-ID", id)

		ctx.Next()
	}
}

func GetCorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationIDKey).(string)
	return id
}
