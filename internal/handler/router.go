package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/team-api-gateway/api-gw-backend/internal/domain"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(cors.Handler)
	router.Use(middleware.Recoverer)
	return router
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func RespondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *domain.HttpError:
		RespondWithJSON(w, e.Status, MakeError(e.Status, "Http Error", e.Error()))
	default:
		RespondWithJSON(w, 500, MakeError(500, "Unknown Error", e.Error()))
	}
}
func H(f func(w http.ResponseWriter, r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := f(w, r)
		if err != nil {
			RespondWithError(w, err)
		} else {
			if result == nil {
				RespondWithJSON(w, 200, MakeError(200, "Success", "Operation was successful"))
			} else {
				RespondWithJSON(w, 200, result)
			}
		}
	}
}

type ErrorObject struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func MakeError(status int, message string, data string) ErrorObject {
	return ErrorObject{Status: status, Message: message, Data: data}
}
