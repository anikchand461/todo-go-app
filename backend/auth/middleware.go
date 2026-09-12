package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization token",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}

		claims, err := ValidateToken(parts[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (int, error) {
	value, exists := c.Get("user_id")

	if !exists {
		return 0, errors.New("user_id not found")
	}

	userID, ok := value.(int)

	if !ok {
		return 0, errors.New("invalid user_id")
	}

	return userID, nil
}
