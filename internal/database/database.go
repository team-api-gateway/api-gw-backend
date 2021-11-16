package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Db struct {
	*mongo.Database
}

func Connect() (*Db, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return nil, fmt.Errorf("DB_HOST is empty")
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		return nil, fmt.Errorf("DB_USER is empty")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is empty")
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		return nil, fmt.Errorf("DB_NAME is empty")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		return nil, fmt.Errorf("DB_PORT is empty")
	}
	dsn := "mongodb://localhost"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database("api-gateway")
	return &Db{db}, nil
}
