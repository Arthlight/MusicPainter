package v1

import (
	"Spotify-Visualizer/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var token models.RefreshToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Wahrscheinlich hier asynchron eine goroutine callen die alle 3 sekunden requests
	// TODO: zur spotify api macht und wenn die goroutine eine response bekommt die bestaetigt
	// TODO: dass der user musik hoert ein event emittet das im frontend "loaded" auf true schaltet
	// TODO: und dann vermutlich eine weitere goroutine callt die dann die computation macht und ans
	// TODO: Frontend die noetigen "draw data" sendet.
	fmt.Println(token.Token)
	w.WriteHeader(200)
}