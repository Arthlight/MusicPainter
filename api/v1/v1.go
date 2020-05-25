package v1

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func NewApiRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/refreshToken", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.Body)
	})

	return router
}