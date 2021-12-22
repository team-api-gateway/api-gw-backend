package database

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
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

func (db *Db) GetResultSpec(id string) (openapi3.T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor := db.Collection("apis").FindOne(ctx, bson.M{"_id": id})
	if cursor.Err() != nil {
		return openapi3.T{}, cursor.Err()
	}
	result := domain.API{}
	cursor.Decode(&result)

	customizedPart := domain.CustomizableAPI{}
	cursor = db.Collection("customized").FindOne(ctx, bson.M{"_ref": id})
	if cursor.Err() == nil {
		cursor.Decode(&customizedPart)
	}

	newPaths := openapi3.Paths{}

	// merge operation descriptions and parameter descriptions
	for key, path := range result.Spec.Paths {
		if path.Get != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Get != nil && customizedPart.Spec.Paths[key].Get.Description != nil {
			result.Spec.Paths[key].Get.Description = *customizedPart.Spec.Paths[key].Get.Description
		}
		if path.Get != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Get != nil && len(customizedPart.Spec.Paths[key].Get.Parameters) > 0 {
			result.Spec.Paths[key].Get.Parameters = mergeParameters(result.Spec.Paths[key].Get.Parameters, customizedPart.Spec.Paths[key].Get.Parameters)
		}
		if path.Delete != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Delete != nil && customizedPart.Spec.Paths[key].Delete.Description != nil {
			result.Spec.Paths[key].Delete.Description = *customizedPart.Spec.Paths[key].Delete.Description
		}
		if path.Delete != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Delete != nil && len(customizedPart.Spec.Paths[key].Delete.Parameters) > 0 {
			result.Spec.Paths[key].Delete.Parameters = mergeParameters(result.Spec.Paths[key].Delete.Parameters, customizedPart.Spec.Paths[key].Delete.Parameters)
		}
		if path.Head != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Head != nil && customizedPart.Spec.Paths[key].Head.Description != nil {
			result.Spec.Paths[key].Head.Description = *customizedPart.Spec.Paths[key].Head.Description
		}
		if path.Head != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Head != nil && len(customizedPart.Spec.Paths[key].Head.Parameters) > 0 {
			result.Spec.Paths[key].Head.Parameters = mergeParameters(result.Spec.Paths[key].Head.Parameters, customizedPart.Spec.Paths[key].Head.Parameters)
		}
		if path.Trace != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Trace != nil && customizedPart.Spec.Paths[key].Trace.Description != nil {
			result.Spec.Paths[key].Trace.Description = *customizedPart.Spec.Paths[key].Trace.Description
		}
		if path.Trace != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Trace != nil && len(customizedPart.Spec.Paths[key].Trace.Parameters) > 0 {
			result.Spec.Paths[key].Trace.Parameters = mergeParameters(result.Spec.Paths[key].Trace.Parameters, customizedPart.Spec.Paths[key].Trace.Parameters)
		}
		if path.Options != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Options != nil && customizedPart.Spec.Paths[key].Options.Description != nil {
			result.Spec.Paths[key].Options.Description = *customizedPart.Spec.Paths[key].Options.Description
		}
		if path.Options != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Options != nil && len(customizedPart.Spec.Paths[key].Options.Parameters) > 0 {
			result.Spec.Paths[key].Options.Parameters = mergeParameters(result.Spec.Paths[key].Options.Parameters, customizedPart.Spec.Paths[key].Options.Parameters)
		}
		if path.Post != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Post != nil && customizedPart.Spec.Paths[key].Post.Description != nil {
			result.Spec.Paths[key].Post.Description = *customizedPart.Spec.Paths[key].Post.Description
		}
		if path.Post != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Post != nil && len(customizedPart.Spec.Paths[key].Post.Parameters) > 0 {
			result.Spec.Paths[key].Post.Parameters = mergeParameters(result.Spec.Paths[key].Post.Parameters, customizedPart.Spec.Paths[key].Post.Parameters)
		}
		if path.Patch != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Patch != nil && customizedPart.Spec.Paths[key].Patch.Description != nil {
			result.Spec.Paths[key].Patch.Description = *customizedPart.Spec.Paths[key].Patch.Description
		}
		if path.Patch != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Patch != nil && len(customizedPart.Spec.Paths[key].Patch.Parameters) > 0 {
			result.Spec.Paths[key].Patch.Parameters = mergeParameters(result.Spec.Paths[key].Patch.Parameters, customizedPart.Spec.Paths[key].Patch.Parameters)
		}
		if path.Put != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Put != nil && customizedPart.Spec.Paths[key].Put.Description != nil {
			result.Spec.Paths[key].Put.Description = *customizedPart.Spec.Paths[key].Put.Description
		}
		if path.Put != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Put != nil && len(customizedPart.Spec.Paths[key].Put.Parameters) > 0 {
			result.Spec.Paths[key].Put.Parameters = mergeParameters(result.Spec.Paths[key].Put.Parameters, customizedPart.Spec.Paths[key].Put.Parameters)
		}

		for _, sel := range customizedPart.Selections {
			if sel.Selected && sel.Path == key {
				switch sel.Method {
				case "Get":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Get = result.Spec.Paths[key].Get
				case "Delete":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Delete = result.Spec.Paths[key].Delete
				case "Head":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Head = result.Spec.Paths[key].Head
				case "Trace":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Trace = result.Spec.Paths[key].Trace
				case "Options":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Options = result.Spec.Paths[key].Options
				case "Post":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Post = result.Spec.Paths[key].Post
				case "Patch":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Patch = result.Spec.Paths[key].Patch
				case "Put":
					if newPaths[key] == nil {
						newPaths[key] = &openapi3.PathItem{}
					}
					newPaths[key].Put = result.Spec.Paths[key].Put
				}
			}
		}
	}
	result.Spec.Paths = newPaths

	return *result.Spec, nil
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

	customizedPart := domain.CustomizableAPI{}
	cursor = db.Collection("customized").FindOne(ctx, bson.M{"_ref": id})
	if cursor.Err() == nil {
		cursor.Decode(&customizedPart)
	}

	// merge operation descriptions and parameter descriptions
	for key, path := range result.Spec.Paths {
		if path.Get != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Get != nil && customizedPart.Spec.Paths[key].Get.Description != nil {
			result.Spec.Paths[key].Get.Description = *customizedPart.Spec.Paths[key].Get.Description
		}
		if path.Get != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Get != nil && len(customizedPart.Spec.Paths[key].Get.Parameters) > 0 {
			result.Spec.Paths[key].Get.Parameters = mergeParameters(result.Spec.Paths[key].Get.Parameters, customizedPart.Spec.Paths[key].Get.Parameters)
		}
		if path.Delete != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Delete != nil && customizedPart.Spec.Paths[key].Delete.Description != nil {
			result.Spec.Paths[key].Delete.Description = *customizedPart.Spec.Paths[key].Delete.Description
		}
		if path.Delete != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Delete != nil && len(customizedPart.Spec.Paths[key].Delete.Parameters) > 0 {
			result.Spec.Paths[key].Delete.Parameters = mergeParameters(result.Spec.Paths[key].Delete.Parameters, customizedPart.Spec.Paths[key].Delete.Parameters)
		}
		if path.Head != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Head != nil && customizedPart.Spec.Paths[key].Head.Description != nil {
			result.Spec.Paths[key].Head.Description = *customizedPart.Spec.Paths[key].Head.Description
		}
		if path.Head != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Head != nil && len(customizedPart.Spec.Paths[key].Head.Parameters) > 0 {
			result.Spec.Paths[key].Head.Parameters = mergeParameters(result.Spec.Paths[key].Head.Parameters, customizedPart.Spec.Paths[key].Head.Parameters)
		}
		if path.Trace != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Trace != nil && customizedPart.Spec.Paths[key].Trace.Description != nil {
			result.Spec.Paths[key].Trace.Description = *customizedPart.Spec.Paths[key].Trace.Description
		}
		if path.Trace != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Trace != nil && len(customizedPart.Spec.Paths[key].Trace.Parameters) > 0 {
			result.Spec.Paths[key].Trace.Parameters = mergeParameters(result.Spec.Paths[key].Trace.Parameters, customizedPart.Spec.Paths[key].Trace.Parameters)
		}
		if path.Options != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Options != nil && customizedPart.Spec.Paths[key].Options.Description != nil {
			result.Spec.Paths[key].Options.Description = *customizedPart.Spec.Paths[key].Options.Description
		}
		if path.Options != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Options != nil && len(customizedPart.Spec.Paths[key].Options.Parameters) > 0 {
			result.Spec.Paths[key].Options.Parameters = mergeParameters(result.Spec.Paths[key].Options.Parameters, customizedPart.Spec.Paths[key].Options.Parameters)
		}
		if path.Post != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Post != nil && customizedPart.Spec.Paths[key].Post.Description != nil {
			result.Spec.Paths[key].Post.Description = *customizedPart.Spec.Paths[key].Post.Description
		}
		if path.Post != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Post != nil && len(customizedPart.Spec.Paths[key].Post.Parameters) > 0 {
			result.Spec.Paths[key].Post.Parameters = mergeParameters(result.Spec.Paths[key].Post.Parameters, customizedPart.Spec.Paths[key].Post.Parameters)
		}
		if path.Patch != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Patch != nil && customizedPart.Spec.Paths[key].Patch.Description != nil {
			result.Spec.Paths[key].Patch.Description = *customizedPart.Spec.Paths[key].Patch.Description
		}
		if path.Patch != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Patch != nil && len(customizedPart.Spec.Paths[key].Patch.Parameters) > 0 {
			result.Spec.Paths[key].Patch.Parameters = mergeParameters(result.Spec.Paths[key].Patch.Parameters, customizedPart.Spec.Paths[key].Patch.Parameters)
		}
		if path.Put != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Put != nil && customizedPart.Spec.Paths[key].Put.Description != nil {
			result.Spec.Paths[key].Put.Description = *customizedPart.Spec.Paths[key].Put.Description
		}
		if path.Put != nil && customizedPart.Spec != nil && customizedPart.Spec.Paths[key] != nil && customizedPart.Spec.Paths[key].Put != nil && len(customizedPart.Spec.Paths[key].Put.Parameters) > 0 {
			result.Spec.Paths[key].Put.Parameters = mergeParameters(result.Spec.Paths[key].Put.Parameters, customizedPart.Spec.Paths[key].Put.Parameters)
		}
	}

	// retrieve selections - first reset all selections for all paths
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
	for _, sel := range customizedPart.Selections {
		for i, s := range result.Selections {
			if s.Path == sel.Path && strings.ToLower(s.Method) == strings.ToLower(sel.Method) {
				result.Selections[i].Selected = sel.Selected
				break
			}
		}
	}

	bytes, _ := json.Marshal(result)
	var c domain.CustomizableAPI
	json.Unmarshal(bytes, &c)
	return c, nil
}

func mergeParameters(existingParams openapi3.Parameters, newParams []domain.Parameter) openapi3.Parameters {
	for i, param := range existingParams {
		for _, p := range newParams {
			if p.Name != "" && p.Name == param.Value.Name {
				existingParams[i].Value.Description = p.Description
				break
			}
		}
	}
	return existingParams
}
