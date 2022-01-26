package database

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		return nil, fmt.Errorf("DB_NAME is empty")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		return nil, fmt.Errorf("DB_PORT is empty")
	}
	var dsn string
	if user == "" && password == "" {
		dsn = "mongodb://" + host + ":" + port
	} else {
		if strings.Contains(host, "azure") {
			dsn = "mongodb://" + user + ":" + password + "@" + host + ":" + port + "/?ssl=true&replicaSet=globaldb&retrywrites=false&maxIdleTimeMS=120000&appName=@hska-lab@"
		} else {
			dsn = "mongodb://" + user + ":" + password + "@" + host + ":" + port
		}
	}
	fmt.Println(dsn)
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

	db := client.Database(dbname)
	return &Db{db}, nil
}
