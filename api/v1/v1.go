package v1

import (
	"Spotify-Visualizer/handler"
	"github.com/go-chi/chi"
	"net/http"
)

func NewApiRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/ws", handler.Socket)

	return router
}