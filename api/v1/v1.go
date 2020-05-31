package v1

import (
	"github.com/go-chi/chi"
	"net/http"
)

func NewApiRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/ws", )

	return router
}