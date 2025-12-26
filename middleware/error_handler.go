package middleware

import (
	"net/http"

	appError "address-book-server/error"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err

		if ae, ok := err.(*appError.AppError); ok {

			response := gin.H{
				"status": "fail",
				"error":   ae.Code,
				"message": ae.Message,
			}

			if len(ae.Details) > 0 {
				response["details"] = ae.Details
			}

			ctx.JSON(ae.StatusCode, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":   "INTERNAL_ERROR",
			"message": "Something went wrong",
		})
	}
}