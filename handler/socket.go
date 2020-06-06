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

	// TODO: Vllt. muss ich die Websocket connection weiterreichen in das spotify module, also sowas wie
	// TODO: SetWebsocketOut zb, also das ich nichtmal die gesamte websocket weiterreiche aber zumindest den Out
	// TODO: Channel, damit ich innerhalb des Spotify Modules die neuen coordinates ans frontend senden kann.
	ws.On("refresh_token", func(event *models.Event) {
		frontend := event.Content.(models.Frontend)
		accessToken, err := spotify.GetAccessToken(frontend.RefreshToken)
		if err != nil {
			fmt.Printf("Unable to get Access Token: %v", err)
		}
		spotify.InitializeAccessToken(accessToken)
		spotify.SetXAndY(frontend.X, frontend.Y)
		go spotify.UpdateAccessTokenAfter(50, frontend.RefreshToken)
		go spotify.LookForCurrentlyPlayingSongWithTimeOut(3)
	})
}