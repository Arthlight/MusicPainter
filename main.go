package main

import (
	"fmt"
	"github.com/go-chi/chi"
	api "Spotify-Visualizer/api/v1"
	"net/http"
)

const port = ":4000"

func main() {
	router := chi.NewRouter()
	v1Handler := api.NewApiRouter()
	router.Mount("/v1", v1Handler)

	fmt.Printf("Started http server on %s\n", port)
	fmt.Println(http.ListenAndServe(port, router))
}