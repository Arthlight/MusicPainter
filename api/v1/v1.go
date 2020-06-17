package v1

import (
	"Spotify-Visualizer/handler"
	"github.com/go-chi/chi"
	"net/http"
)

func NewApiRouter() http.Handler {
	router := chi.NewRouter()

	router.HandleFunc("/ws", handler.Socket)

	return router
}