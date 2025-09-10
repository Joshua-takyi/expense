package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua/expensetracker/internal/auth"
	"github.com/joshua/expensetracker/internal/handlers"
	"github.com/joshua/expensetracker/internal/models"
)

func Router(s models.Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1 := r.Group("/api/v1")

	// public routes
	{
		v1.POST("/register", handlers.RegisterUser(s))
		v1.POST("/login", handlers.AuthenticateUser(s))
	}

	// protected routes

	protected := v1.Group("/").Use(auth.Middleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			user, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"user": user})
		})
		// v1.GET("/profile", handlers.GetProfile(s))
		// v1.PUT("/profile", handlers.UpdateProfile(s))

		// v1.POST("/categories", handlers.CreateCategory(s))
		// v1.GET("/categories", handlers.GetCategories(s))
		// v1.PUT("/categories/:id", handlers.UpdateCategory(s))
		// v1.DELETE("/categories/:id", handlers.DeleteCategory(s))

		// v1.POST("/transactions", handlers.CreateTransaction(s))
		// v1.GET("/transactions", handlers.GetTransactions(s))
		// v1.PUT("/transactions/:id", handlers.UpdateTransaction(s))
		// v1.DELETE("/transactions/:id", handlers.DeleteTransaction(s))
	}
	return r

}
