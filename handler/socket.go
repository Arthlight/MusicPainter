package handler

import (
	"Spotify-Visualizer/spotify"
	"Spotify-Visualizer/models"
	"fmt"
	"go/token"
	"net/http"
)

func Socket(w http.ResponseWriter, r *http.Request) {
	ws, err := models.CreateNewWebsocket(w, r)
	if err != nil {
		fmt.Printf("Error while trying to create new Websocket Connection: %v", err)
	}

	ws.On("refresh_token", func(event *models.Event) {
		frontend := event.Content.(models.Frontend)
		accessToken, err:= spotify.GetAccessToken(frontend.RefreshToken)

	})
}