package database

import (
	"context"
	"time"

	"github.com/team-api-gateway/api-gw-backend/internal/domain"
)

type Log struct {
	domain.CustomizableAPI `bson:",inline"`
	Timestamp              time.Time
}

func (db *Db) WriteLog(id string, customization domain.CustomizableAPI) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	customization.ApiId = id
	l := Log{
		CustomizableAPI: customization,
		Timestamp:       time.Now(),
	}
	_, err := db.Collection("log").InsertOne(ctx, l)
	return err
}
