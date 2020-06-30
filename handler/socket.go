package handler

import (
	"Spotify-Visualizer/models"
	"Spotify-Visualizer/spotify"
	"encoding/json"
	"fmt"
	"net/http"
)

func Socket(w http.ResponseWriter, r *http.Request) {
	ws, err := models.CreateNewWebsocket(w, r)
	if err != nil {
		fmt.Printf("Error while trying to create new Websocket Connection: %v", err)
	}

	ws.On("refresh_token", func(event *models.Event) {
		var frontend models.FromFrontend
		eventAsByte, _ := json.Marshal(event.Content)
		fmt.Println(json.Unmarshal(eventAsByte, &frontend))
		accessToken, err := spotify.GetAccessToken(frontend.RefreshToken)
		if err != nil {
			fmt.Printf("Unable to get Access Token: %v", err)
			return
		}
		spotify.InitializeAccessToken(accessToken)
		spotify.SetXAndY(frontend.X, frontend.Y)
		spotify.SetPipeline(&ws.Out)
		go spotify.UpdateAccessTokenAfter(50, frontend.RefreshToken)
		go spotify.LookForCurrentlyPlayingSongWithTimeOut(2)
	})

	fmt.Println("Successfully upgraded incoming request to websocket connection!")
}