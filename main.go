package main

import (
	api "Spotify-Visualizer/api/v1"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

const address = "localhost:4000"

func main() {
	router := chi.NewRouter()
	v1Handler := api.NewApiRouter()
	router.Mount("/v1", v1Handler)

	fmt.Printf("Started http server on %s\n", address)
	err := http.ListenAndServe(address, router)
	fmt.Println(err)

}