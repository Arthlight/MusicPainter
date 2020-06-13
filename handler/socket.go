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
		frontend := event.Content.(models.FromFrontend)
		accessToken, err := spotify.GetAccessToken(frontend.RefreshToken)
		if err != nil {
			fmt.Printf("Unable to get Access Token: %v", err)
			return // TODO: Maybe enhance this a little bit here and send a notif to the frontend that the refreshtoken
			       // TODO: is expired, but you can't know that here for certain so additional checks would be required for that
		}
		spotify.InitializeAccessToken(accessToken)
		spotify.SetXAndY(frontend.X, frontend.Y)
		spotify.SetPipeline(&ws.Out)
		go spotify.UpdateAccessTokenAfter(50, frontend.RefreshToken)
		go spotify.LookForCurrentlyPlayingSongWithTimeOut(2) // TODO: Does this need to be a goroutine?
	})
}