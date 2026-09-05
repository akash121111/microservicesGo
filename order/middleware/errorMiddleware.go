package middleware

import (
	"errors"
	"net/http"
	"order/apperror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 {
			return
		}
		err := ctx.Errors.Last().Err

		//validation error
		var validationEroors validator.ValidationErrors

		if errors.As(err, &validationEroors) {
			validateMessage := make(map[string]string)
			for _, fieldErr := range validationEroors {
				switch fieldErr.Tag() {
				case "required":
					validateMessage[fieldErr.Field()] = "is required"
				case "email":
					validateMessage[fieldErr.Field()] = "must be valid email"
				default:
					validateMessage[fieldErr.Field()] = "is Invalid"
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   validateMessage,
			})
			return
		}

		//application Eroor
		switch {
		case errors.Is(err, apperror.ErrInvalidUserID):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  err.Error(),
			})
		case errors.Is(err, apperror.ErrUserAlreadyExist):
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"errors":  err.Error(),
			})
		case errors.Is(err, apperror.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"errors":  err.Error(),
			})
		case errors.Is(err, apperror.ErrInvalidCredentials):
			ctx.JSON(401, gin.H{
				"success": false,
				"errors":  err.Error(),
			})
		case errors.Is(err, apperror.ErrUnauthorized):
			ctx.JSON(401, gin.H{
				"success": false,
				"errors":  err.Error(),
			})
		case errors.Is(err, apperror.ErrForbidden):
			ctx.JSON(403, gin.H{
				"success": false,
				"errors":  err.Error(),
			})
		case errors.Is(err, apperror.ErrProductAlreadyExist):
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"errors":  err.Error(),
			})
		case errors.Is(err, apperror.ErrProductNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"errors":  err.Error(),
			})

		case errors.Is(err, apperror.ErrInsuffiecientProduct):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  err.Error(),
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"errors":  "Internal server Error",
			})

		}

	}
}
