package middleware

import (
	appError "address-book-server/error"
	"address-book-server/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.Error(appError.Unauthorized(
				"Invalid token",
				nil,
			))
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Error(appError.Forbidden(
				"Invalid authorization header",
				nil,
			))
			ctx.Abort()
			return
		}

		token, _ := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error){
			return utils.JwtSecret(), nil
		})



		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set("user_id", uint64(claims["user_id"].(float64)))
			ctx.Set("email", string(claims["email"].(string)))
			ctx.Next()
			return
		}

		ctx.Error(appError.Forbidden(
			"Invalid token",
			nil,
		))

		ctx.Abort()
	}
}