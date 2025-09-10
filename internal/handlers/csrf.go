package handlers

import (
	"net/http"

	"github.com/Joshua-takyi/expense/server/internal/helpers"
	"github.com/gin-gonic/gin"
)

func CSRFHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrfToken, err := helpers.GenerateCsrfToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to generate csrf token",
				"error":   err.Error(),
			})
			return
		}

		c.SetCookie("csrf_token", csrfToken, 3600*24*7, "/", "localhost", true, false)

		c.JSON(http.StatusOK, gin.H{
			"message":    "CSRF token provided",
			"csrf_token": csrfToken,
		})
	}
}
