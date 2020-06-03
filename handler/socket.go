package handler

import (
	"Spotify-Visualizer/models"
	"Spotify-Visualizer/spotify"
	"fmt"
	"net/http"
)

func Socket(w http.ResponseWriter, r *http.Request) {
	ws, err := models.CreateNewWebsocket(w, r)
	if err != nil {
		fmt.Printf("Error while trying to create new Websocket Connection: %v", err)
	}

	ws.On("refresh_token", func(event *models.Event) {
		frontend := event.Content.(models.Frontend)
		accessToken, err := spotify.GetAccessToken(frontend.RefreshToken)
		if err != nil {
			fmt.Printf("Unable to get Access Token: %v", err)
		}
		spotify.InitializeAccessToken(accessToken)
		go spotify.UpdateAccessTokenAfter(50, frontend.RefreshToken)
		go spotify.LookForCurrentlyPlayingSongWithTimeOut(3)
		go spotify.ComputeNextCoordinatesFromSongInfo(frontend.X, frontend.Y)


	})
}