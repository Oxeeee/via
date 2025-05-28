package middleware

import (
	"strings"
	"time"

	"github.com/OxytocinGroup/theca-v3/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware interface {
	JWTMiddleware() gin.HandlerFunc
}

type middleware struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewAuthMiddleware(accessSecret, refreshSecret []byte) AuthMiddleware {
	return &middleware{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (mw *middleware) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.RespondWithError(c, errors.New(errors.CodeUnauthorized, "Unauthorized"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			errors.RespondWithError(c, errors.New(errors.CodeUnauthorized, "Invalid auth header format"))
			c.Abort()
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return mw.accessSecret, nil
		})

		if err != nil || !token.Valid {
			errors.RespondWithError(c, errors.New(errors.CodeUnauthorized, "Invalid or expired token"))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
				errors.RespondWithError(c, errors.New(errors.CodeUnauthorized, "Token expired"))
				c.Abort()
				return
			}
			if userIDFloat, ok := claims["userId"].(float64); ok {
				userID := uint(userIDFloat)
				c.Set("userID", userID)
			} else {
				errors.RespondWithError(c, errors.New(errors.CodeUnauthorized, "Invalid token payload"))
				c.Abort()
				return
			}
		} else {
			errors.RespondWithError(c, errors.New(errors.CodeUnauthorized, "Invalid token payload"))
			c.Abort()
			return
		}

		c.Next()
	}
}
