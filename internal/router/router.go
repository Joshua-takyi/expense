package router

import (
	"net/http"

	"github.com/Joshua-takyi/expense/server/internal/auth"
	"github.com/Joshua-takyi/expense/server/internal/handlers"
	"github.com/Joshua-takyi/expense/server/internal/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router(s models.Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	allowedOrigins := []string{
		"http://localhost:3000",
		"https://expense-1-kblg.onrender.com",
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "pong",
		})
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

		protected.GET("/csrf-token", handlers.CSRFHandler())
		protected.POST("/logout", handlers.LogoutUser)
		// protected.GET("/profile", handlers.GetProfile(s))
		// protected.PUT("/profile", handlers.UpdateProfile(s))

		// protected.POST("/categories", handlers.CreateCategory(s))
		// protected.GET("/categories", handlers.GetCategories(s))
		// protected.PUT("/categories/:id", handlers.UpdateCategory(s))
		// protected.DELETE("/categories/:id", handlers.DeleteCategory(s))

		protected.POST("/transactions", handlers.AddTransaction(s))
		protected.GET("/transactions-query/", handlers.QueryTransactions(s))
		protected.GET("/transactions", handlers.ListUserTransactions(s))
		// protected.PUT("/transactions/:id", handlers.UpdateTransaction(s))
		protected.DELETE("/transactions/:id", handlers.RemoveTransaction(s))
	}
	return r

}
