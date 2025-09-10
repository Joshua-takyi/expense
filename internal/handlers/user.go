package handlers

import (
	"os"

	"github.com/Joshua-takyi/expense/server/internal/constants"
	"github.com/Joshua-takyi/expense/server/internal/helpers"
	"github.com/Joshua-takyi/expense/server/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterUser(r models.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "invalid request body"})
			return
		}
		if err := r.RegisterUser(ctx, &user); err != nil {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "user registered successfully"})
	}
}

func AuthenticateUser(r models.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := &models.User{}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "invalid request body"})
			return
		}

		user, err := r.AuthenticateUser(ctx, req.Email, req.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": constants.ErrUnauthorized, "message": "invalid email or password"})
			return
		}
		secret := os.Getenv("JWT_SECRET")
		claims := &helpers.UseClaims{
			UserID: user.Id.Hex(),
			Email:  user.Email,
		}

		token, err := helpers.GenerateJWT(claims, secret)
		if err != nil {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "failed to generate token"})
			return
		}

		// generate csrf token
		csrfToken, err := helpers.GenerateCsrfToken()
		if err != nil {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "failed to generate csrf token"})
			return
		}

		// set secure to false for development, true for production
		isProduction := os.Getenv("GIN_MODE") == "release" || os.Getenv("NODE_ENV") == "production"
		c.SetCookie("csrf_token", csrfToken, 3600*24*7, "/", "", isProduction, true)
		c.SetCookie("auth_token", token, 3600*24*7, "/", "", isProduction, true)

		c.JSON(200, gin.H{"token": token, "data": user})

	}
}
