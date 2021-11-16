package sync

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/team-api-gateway/api-gw-backend/internal/database"
	"github.com/team-api-gateway/api-gw-backend/internal/domain"
	"github.com/team-api-gateway/api-gw-backend/internal/hashid"
	"go.mongodb.org/mongo-driver/bson"
)

func Sync(db *database.Db, specs []*openapi3.T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get all apis that exist in db
	cur, err := db.Collection("apis").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	apisFromDB := []domain.API{}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result domain.API
		err := cur.Decode(&result)
		if err != nil {
			return err
		}
		apisFromDB = append(apisFromDB, result)
	}
	if err := cur.Err(); err != nil {
		return err
	}

	// convert spec to domain.API
	apisFromCrawl := []domain.API{}
	for _, spec := range specs {
		id := hashid.GenerateHashId(spec.Info.Title)
		apisFromCrawl = append(apisFromCrawl, domain.API{Spec: spec, Id: id})
	}

	// find all apis to create and to update
	apisToCreate := []domain.API{}
	apisToUpdate := []domain.API{}
	for _, api := range apisFromCrawl {
		notInDb := true
		needsUpdate := false
		for _, a := range apisFromDB {
			if a.Id == api.Id {
				notInDb = false
				if a.Spec.Info.Version != api.Spec.Info.Version {
					needsUpdate = true
				}
			}
		}
		if notInDb {
			apisToCreate = append(apisToCreate, api)
		} else if needsUpdate {
			apisToUpdate = append(apisToUpdate, api)
		}
	}
	// find all apis to delete
	apisToDelete := []domain.API{}
	for _, a := range apisFromDB {
		notInCrawl := true
		for _, api := range apisFromCrawl {
			if a.Id == api.Id {
				notInCrawl = false
			}
		}
		if notInCrawl {
			apisToDelete = append(apisToDelete, a)
		}
	}
	for _, apiToDelete := range apisToDelete {
		fmt.Println("deleting " + apiToDelete.Spec.Info.Title)
		_, err := db.Collection("apis").DeleteOne(ctx, bson.M{"_id": apiToDelete.Id})
		if err != nil {
			return err
		}
		err = db.DeleteCustomization(apiToDelete.Id)
		if err != nil {
			return err
		}
	}
	for _, apiToCreate := range apisToCreate {
		fmt.Println("creating " + apiToCreate.Spec.Info.Title)
		_, err := db.Collection("apis").InsertOne(ctx, apiToCreate)
		if err != nil {
			return err
		}
	}
	for _, apiToUpdate := range apisToUpdate {
		fmt.Println("updating " + apiToUpdate.Spec.Info.Title)
		_, err := db.Collection("apis").ReplaceOne(ctx, bson.M{"_id": apiToUpdate.Id}, apiToUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}
