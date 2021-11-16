package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	router.Post("/apis/{id}/update", H(h.CustomizeApi))
	router.Post("/apis/{id}/selection", H(h.SetSelectionStatus))

	return router
}

func NewHandler(db *database.Db) *handler {
	return &handler{
		Db: db,
	}
}

func (h *handler) CustomizeApi(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var api domain.CustomizableAPI
	err := json.NewDecoder(r.Body).Decode(&api)
	if err != nil {
		return nil, err
	}
	return h.Db.CustomizeApi(chi.URLParam(r, "id"), api)
}
func (h *handler) SetSelectionStatus(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, &domain.HttpError{Status: 501, Err: fmt.Errorf("not yet implemented")}
}

// swago.tag: API
// swago.response: []api
func (h *handler) GetApis(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	apis, err := h.Db.GetApis()
	if err != nil {
		return nil, err
	}
	type ID struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	}
	ids := []ID{}
	for _, api := range apis {
		ids = append(ids, ID{Id: api.Id, Title: api.Spec.Info.Title})
	}
	return ids, nil
}

// swago.tag: API
// swago.response: api
func (h *handler) GetApi(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	return h.Db.GetApi(id)
}
