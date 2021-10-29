package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type API struct {
	Id      int64
	OpenAPI []byte
}

type Db struct {
	*gorm.DB
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
	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		return nil, fmt.Errorf("DB_DATABASE is empty")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		return nil, fmt.Errorf("DB_PORT is empty")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Berlin", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(API{})
	if err != nil {
		return nil, err
	}
	return &Db{db}, nil
}
