package handlers

import (
	"github.com/Joshua-takyi/expense/server/internal/helpers"
	"github.com/gin-gonic/gin"
)

func CsrfHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrfToken, err := helpers.GenerateCsrfToken()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate CSRF token"})
			return
		}

		// set csrf token in cookie
		c.SetCookie("csrf_token", csrfToken, 3600, "/", "", false, true)

		c.JSON(200, gin.H{"message": "CSRF token set", "csrf_token": csrfToken})
	}
}
