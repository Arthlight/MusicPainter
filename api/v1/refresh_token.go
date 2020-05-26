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

	fmt.Println(token.Token)
	w.WriteHeader(200)
}