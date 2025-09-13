package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Joshua-takyi/expense/server/internal/constants"
	"github.com/Joshua-takyi/expense/server/internal/helpers"
	"github.com/Joshua-takyi/expense/server/internal/models"
)

func AddTransaction(r models.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var tx models.Transaction
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(400, gin.H{"error": constants.ErrInvalidInput, "message": " invalid input"})
			return
		}

		userClaims, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": constants.ErrUnauthorized, "message": "unauthorized"})
			return
		}

		claims, ok := userClaims.(*helpers.UseClaims)
		if !ok {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "internal server error"})
			return

		}

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "invalid user ID"})
			return
		}
		if err := r.AddTransaction(ctx, &tx, userID); err != nil {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "failed to add transaction"})
			return
		}
		c.JSON(201, gin.H{"message": "transaction added successfully", "transaction": tx})
	}
}

func QueryTransactions(r models.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "10")
		offset := c.DefaultQuery("offset", "0")
		ctx := c.Request.Context()
		userClaims, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": constants.ErrUnauthorized, "message": "unauthorized"})
			return
		}

		claims, ok := userClaims.(*helpers.UseClaims)
		if !ok {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "internal server error"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "invalid user ID"})
			return
		}

		query := c.Query("search")
		category := c.QueryArray("category")
		order := c.Query("order")

		transactions, err := r.GetTransactionByQuery(ctx, userID, query, category, order, helpers.ParseInt(limit), helpers.ParseInt(offset))
		if err != nil {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "failed to search transactions"})
			return
		}
		c.JSON(200, gin.H{"data": transactions})
	}
}

func ListUserTransactions(r models.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "10")
		offset := c.DefaultQuery("offset", "0")
		ctx := c.Request.Context()
		userClaims, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": constants.ErrUnauthorized, "message": "unauthorized"})
			return
		}

		claims, ok := userClaims.(*helpers.UseClaims)
		if !ok {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "internal server error"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "invalid user ID"})
			return
		}

		transactions, err := r.ListUserTransactions(ctx, userID, helpers.ParseInt(limit), helpers.ParseInt(offset))
		if err != nil {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "failed to list transactions"})
			return
		}
		c.JSON(200, gin.H{"data": transactions})
	}
}

func RemoveTransaction(r models.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		txID := c.Param("id")
		if txID == "" {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "transaction ID is required"})
			return
		}
		id, err := primitive.ObjectIDFromHex(txID)
		if err != nil {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "invalid transaction ID"})
			return
		}

		userClaims, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": constants.ErrUnauthorized, "message": "unauthorized"})
			return
		}

		claims, ok := userClaims.(*helpers.UseClaims)
		if !ok {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "internal server error"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(400, gin.H{"error": constants.ErrBadRequest, "message": "invalid user ID"})
			return
		}

		tx, err := r.GetTransactionDetails(ctx, id)
		if err != nil {
			c.JSON(404, gin.H{"error": constants.ErrNoDocuments, "message": "transaction not found"})
			return
		}

		if tx.UserId != userID {
			c.JSON(403, gin.H{"error": constants.ErrForbidden, "message": "forbidden"})
			return
		}

		if err := r.RemoveTransaction(ctx, id); err != nil {
			c.JSON(500, gin.H{"error": constants.ErrInternalServer, "message": "failed to remove transaction"})
			return
		}
		c.JSON(200, gin.H{"message": "transaction removed successfully"})
	}
}
