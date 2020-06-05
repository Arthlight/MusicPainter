package main

import (
	"Spotify-Visualizer/spotify"
	"fmt"
)

//const port = ":4000"

func main() {
/*	router := chi.NewRouter()
	v1Handler := api.NewApiRouter()
	router.Mount("/v1", v1Handler)

	fmt.Printf("Started http server on %s\n", port)
	err := http.ListenAndServe(port, router)
	fmt.Println(err)*/
	//fmt.Println(spotify.GetAccessToken(""))
	str, ok := spotify.GetCurrentTrackID()
	fmt.Println(str, ok)
}