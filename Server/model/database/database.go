package database

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var conn Connection

func InitDatabase() {
	err := godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("Failed: Loading environments variables")
		return
	}

	db_uri := os.Getenv("DATABASE_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db_uri))
	if err != nil {
		fmt.Println("Failed: Connecting to Database")
		return
	}

	db := client.Database("test")

	conn = Connection{
		Client: client,
		DB:     db,
	}
}

func GetCollection(name string) *mongo.Collection {
	return conn.DB.Collection(name)
}
