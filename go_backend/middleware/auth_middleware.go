package middleware

import (
	"net/http"
	"strings"

	"github.com/Dhanrajsinh002/go-url-shortener/auth"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "error": "missing or invalid authorization header" })
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "error": "invalid or expired token" })
			return
		}
		c.Set("username", claims.Username)
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}