package database

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/team-api-gateway/api-gw-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *Db) GetApis() ([]domain.API, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := db.Collection("apis").Find(ctx, bson.M{})
	if err != nil {
		return []domain.API{}, err
	}
	var apis []domain.API
	if err = cursor.All(ctx, &apis); err != nil {
		return []domain.API{}, err
	}
	return apis, nil
}
func (db *Db) GetApi(id string) (domain.CustomizableAPI, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor := db.Collection("apis").FindOne(ctx, bson.M{"_id": id})
	if cursor.Err() != nil {
		return domain.CustomizableAPI{}, cursor.Err()
	}
	result := domain.API{}
	cursor.Decode(&result)

	for key, path := range result.Spec.Paths {
		if path.Delete != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodDelete, Selected: false})
		}
		if path.Get != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodGet, Selected: false})
		}
		if path.Head != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodHead, Selected: false})
		}
		if path.Options != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodOptions, Selected: false})
		}
		if path.Patch != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodPatch, Selected: false})
		}
		if path.Post != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodPost, Selected: false})
		}
		if path.Put != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodPut, Selected: false})
		}
		if path.Trace != nil {
			result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodTrace, Selected: false})
		}
	}

	cursor = db.Collection("customized").FindOne(ctx, bson.M{"_ref": id})
	if cursor.Err() == nil {
		customizedPart := domain.API{}
		cursor.Decode(&customizedPart)

		err := mergo.Merge(&result, customizedPart, mergo.WithOverride)
		if err != nil {
			return domain.CustomizableAPI{}, err
		}
		for key, path := range result.Spec.Paths {
			if path.Delete != nil && !selectionExists(result.Selections, key, http.MethodDelete) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodDelete, Selected: false})
			}
			if path.Get != nil && !selectionExists(result.Selections, key, http.MethodGet) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodGet, Selected: false})
			}
			if path.Head != nil && !selectionExists(result.Selections, key, http.MethodHead) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodHead, Selected: false})
			}
			if path.Options != nil && !selectionExists(result.Selections, key, http.MethodOptions) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodOptions, Selected: false})
			}
			if path.Patch != nil && !selectionExists(result.Selections, key, http.MethodPatch) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodPatch, Selected: false})
			}
			if path.Post != nil && !selectionExists(result.Selections, key, http.MethodPost) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodPost, Selected: false})
			}
			if path.Put != nil && !selectionExists(result.Selections, key, http.MethodPut) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodPut, Selected: false})
			}
			if path.Trace != nil && !selectionExists(result.Selections, key, http.MethodTrace) {
				result.Selections = append(result.Selections, domain.Selection{Path: key, Method: http.MethodTrace, Selected: false})
			}
		}
		bytes, _ := json.Marshal(result)
		var c domain.CustomizableAPI
		json.Unmarshal(bytes, &c)
		return c, nil
	}
	bytes, _ := json.Marshal(result)
	var c domain.CustomizableAPI
	json.Unmarshal(bytes, &c)
	return c, nil
}

func selectionExists(selections []domain.Selection, path string, method string) bool {
	for _, sel := range selections {
		if sel.Path == path && strings.ToLower(method) == strings.ToLower(sel.Method) {
			return true
		}
	}
	return false
}
