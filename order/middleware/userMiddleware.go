package middleware

import (
	"github.com/gin-gonic/gin"
)

func UserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userID := ctx.GetHeader("X-USER-ID")
		if userID == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"error":   "Authentication is Failed",
			})
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("email", ctx.GetHeader("X-USER-EMAIL"))
		ctx.Set("role", ctx.GetHeader("X-USER-ROLE"))

		ctx.Next()
	}
}
