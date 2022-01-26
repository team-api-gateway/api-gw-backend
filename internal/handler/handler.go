package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/team-api-gateway/api-gw-backend/internal/azure"
	"github.com/team-api-gateway/api-gw-backend/internal/database"
	"github.com/team-api-gateway/api-gw-backend/internal/domain"
)

type handler struct {
	Db *database.Db
}

func (h *handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/apis", H(h.GetApis))
	router.Get("/apis/{id}", H(h.GetApi))
	router.Get("/apis/{id}/spec", H(h.GetSpec))
	router.Put("/apis/{id}/spec", H(h.UploadSpec))
	router.Post("/apis/{id}/update", H(h.CustomizeApi))
	router.Post("/set-self-host/{self-host}", H(h.SetSelfHost))

	return router
}
func (h *handler) SetSelfHost(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, os.Setenv("SELF_HOST", "https://"+chi.URLParam(r, "self-host"))
}

func NewHandler(db *database.Db) *handler {
	return &handler{
		Db: db,
	}
}

// @Title Customize the api
// @Description Update the description and the selection state for endpoints in a api
// @Param id path string true "Id of the api"
// @Param object body domain.CustomizableAPI true "The customized part of the api"
// @Success 200 object domain.CustomizableAPI "Customized API"
// @Failure  500  object  ErrorObject  "ErrorResponse"
// @Resource apis
// @Route /apis/{id}/update [post]
func (h *handler) CustomizeApi(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var api domain.CustomizableAPI
	err := json.NewDecoder(r.Body).Decode(&api)
	if err != nil {
		return nil, err
	}
	if api.Username == "" {
		return nil, fmt.Errorf("no username provided")
	}
	if err := h.Db.WriteLog(chi.URLParam(r, "id"), api); err != nil {
		return nil, err
	}
	return h.Db.CustomizeApi(chi.URLParam(r, "id"), api)
}

type ID struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type ArrayOfIds []ID

// @Title Get ids and titles of all apis
// @Description Get the ids and titles of all apis that are stored in the database
// @Success  200  object  ArrayOfIds  "List of IDs"
// @Failure  500  object  ErrorObject  "ErrorResponse"
// @Resource apis
// @Route /apis [get]
func (h *handler) GetApis(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	apis, err := h.Db.GetApis()
	if err != nil {
		return nil, err
	}
	ids := []ID{}
	for _, api := range apis {
		ids = append(ids, ID{Id: api.Id, Title: api.Spec.Info.Title})
	}
	return ids, nil
}

// @Title Get details for one api
// @Description Get the details for one specified api
// @Param  id path string true "Id of the api"
// @Success  200  object  domain.CustomizableAPI  "Customized API"
// @Failure  500  object  ErrorObject  "ErrorResponse"
// @Resource apis
// @Route /apis/{id} [get]
func (h *handler) GetApi(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	return h.Db.GetApi(id)
}

func (h *handler) GetSpec(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	return h.Db.GetResultSpec(id)
}

// @Title Upload a spec
// @Description Upload a spec to the azure api managment gateway
// @Param  id path string true "Id of the api"
// @Success 200 object ErrorObject "Success"
// @Failure  500  object  ErrorObject  "ErrorResponse"
// @Resource apis
// @Route /apis/{id}/spec [put]
func (h *handler) UploadSpec(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	//host := os.Getenv("SELF_HOST")

	spec, err := h.Db.GetResultSpec(id)
	if err != nil {
		return nil, err
	}
	if len(spec.Paths) == 0 {
		fmt.Println("no endpoints selected - removing api")
		return nil, azure.DeleteSpec(id)
	}
	jsonString, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(jsonString))
	content := `{
		"properties": {
		  "format": "openapi+json",
		  "value": "` + strings.ReplaceAll(string(jsonString), `"`, `\"`) + `",
		  "path": "` + strings.ReplaceAll(spec.Info.Title, " ", "-") + `",
		  "subscriptionRequired": false
		}
	  }`
	err = azure.UploadSpec(id, content)
	if err != nil {
		return nil, err
	}
	fmt.Println("upload was successful")
	time.Sleep(15 * time.Second)
	return nil, azure.AddToProduct(id)
}
