package database

import (
	"context"
	"time"

	"github.com/team-api-gateway/api-gw-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *Db) CustomizeApi(id string, custom domain.CustomizableAPI) (domain.API, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	custom.ApiId = id
	cursor := db.Collection("customized").FindOne(ctx, bson.M{"_ref": id})
	if cursor.Err() != nil {
		_, err := db.Collection("customized").InsertOne(ctx, custom)
		if err != nil {
			return domain.API{}, err
		}
	} else {
		_, err := db.Collection("customized").ReplaceOne(ctx, bson.M{"_ref": custom.ApiId}, custom)
		if err != nil {
			return domain.API{}, err
		}
	}
	return db.GetApi(custom.ApiId)
}
func (db *Db) DeleteCustomization(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := db.Collection("customized").DeleteOne(ctx, bson.M{"_ref": id})
	if err != nil {
		return err
	}
	return nil
}
