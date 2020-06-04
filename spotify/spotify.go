package spotify

import (
	"Spotify-Visualizer/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var currentAccessToken string
var currentTrackID string

func GetAccessToken(refreshToken string) (string, error) {
	client := http.Client{}
	form := url.Values{
		"grant_type": {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	encodedAuthHeader, err := getEncodedAuthHeader(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET"))
	if err != nil {
		fmt.Println("Error while trying to encode auth header: ", err)
	}
	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Printf("Error while trying to create a new request in GetAccessToken: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic " + encodedAuthHeader)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while trying to send request in GetAccessToken: ", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while trying to read response body in GetAccessToken: ", err)
	}
	var tr models.TokenResponse
	err = json.Unmarshal(body, &tr)
	if err != nil {
		fmt.Println("Error while trying to unmarshal response body in GetAccessToken: ", err)
	}

	return tr.AccessToken, nil
}

func getEncodedAuthHeader(credentials string) (string, error) {
	var sb strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &sb)
	_, err := encoder.Write([]byte(credentials))
	if err != nil {
		return "", err
	}
	err = encoder.Close()
	if err != nil {
		return "", err
	}
	encodedAuthHeader := sb.String()

	return encodedAuthHeader, nil
}

func InitializeAccessToken(accessToken string) {
	currentAccessToken = accessToken
}

func UpdateAccessTokenAfter(timeout int, refreshToken string) {
	for {
		time.Sleep(time.Second * time.Duration(timeout))
		accessToken, err := GetAccessToken(refreshToken)
		if err != nil {
			fmt.Println("Error receiving Access Token in UpdateAccessTokenAfter: ", err)
		}
		currentAccessToken = accessToken
	}
}

// TODO: Here I probably need to notify the frontend every 3 seconds whether a song is currently playing or not,
// TODO: So dont forget to send a notification into the "out" channel of the socket here later on
func LookForCurrentlyPlayingSongWithTimeOut(timeout int) {
	for {
		trackID, ok := getCurrentTrackID()
		if ok {
			// TODO: send a notification to the frontend that user is currently playing a song
			currentTrackID = trackID
		}
		time.Sleep(time.Second * time.Duration(timeout))
	}
}


func getCurrentTrackID() (string, bool){
	client := http.Client{}
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		fmt.Printf("Error while trying to create a new request in getCurrentTrackID: %v", err)
		return "", false
	}
	request.Header.Set("Authorization", currentAccessToken)
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while trying to send request in getCurrentTrackID: ", err)
		return "", false
	}
	defer response.Body.Close()

	// The check "contentLength == 0" may be insufficient and I might need to provide more complex measures
	// in order to assure that the response.Body is actually empty
	if response.StatusCode == 204 || response.ContentLength == 0 {
		return "", false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while trying to read response body in getCurrentTrackID: ", err)
		return "", false
	}
	var flags struct {
		isPlaying bool `json:"is_playing"`
		id string `json:"id"`
	}
	err = json.Unmarshal(body, &flags)
	if err != nil {
		fmt.Println("Error while trying to unmarshal response body in getCurrentTrackID: ", err)
		return "", false
	}
	if !flags.isPlaying {
		return "", false
	}

	return flags.id, true




	// TODO: Check whether the response.StatusCode is 200 or 204 and whether "currently_playing" is set to true. In case it is
	// TODO: set to True and the statuscode is 200, return the track ID of the song (make sure to actually return the track
	// TODO: id of the song and not the album)
}

// TODO: Hier hab ich wahrsch auch einen while loop wodrin ich mir dann immer werte von der current track id
// TODO: in einer separaten function nochmal hole und eine calculatete anzahl ein punkten die in eine bestimmte richtung
// TODO: gehen ans frontend zurücksende, bevor ich dann die richtung wechsel und sich der prozess wiederholt. Am Anfang
// TODO: des while loops hol ich mir dann immer wieder neu werte von der trackID die sich geändert haben könnte und passe
// TODO: dadurch dann ggf. den ouput den ich zurücksende ans frontend an (neue kreisgroeße, farbe, form etc)
func ComputeNextCoordinatesFromSongInfo(x, y int) {
	for {
		if currentTrackID != "" {
			audioFeatures, err := getAudioFeaturesOfTrack()
		}
	}
}

func getAudioFeaturesOfTrack() {

}