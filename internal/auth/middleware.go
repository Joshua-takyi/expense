package auth

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joshua/expensetracker/internal/helpers"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		// get token from the request header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization token not provided"})
			c.Abort()
			return
		}

		// validate the token
		claims, err := helpers.ValidateToken(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}

		c.Set("user", claims)

		// csrf token protecting mutation and post requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			csrfCookie, err := c.Cookie("csrf_token")
			if err != nil {
				c.AbortWithStatusJSON(403, gin.H{"error": "invalid CSRF token"})
				c.Abort()
				return
			}

			csrfHeader := c.GetHeader("X-CSRF-Token")
			if csrfHeader == "" {
				c.AbortWithStatusJSON(403, gin.H{"error": "invalid CSRF token"})
				c.Abort()
				return
			}

			if csrfCookie != csrfHeader {
				c.AbortWithStatusJSON(403, gin.H{"error": "invalid CSRF token"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
