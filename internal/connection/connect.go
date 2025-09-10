package connection

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func InitDb() error {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, reading configuration from environment variables")
	}
	uri := os.Getenv("MONGODB_URI")
	password := os.Getenv("MONGODB_PASSWORD")
	fullUri := strings.Replace(uri, "<password>", password, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(fullUri)

	var err error
	Client, err = mongo.Connect(clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err := Client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("✅ MongoDB connected successfully")
	return nil
}

func CloseDb() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := Client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect MongoDB client: %v", err)
		}
		fmt.Println("✅ MongoDB disconnected successfully")
	}
	return nil
}
