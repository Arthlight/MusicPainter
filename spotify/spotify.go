package spotify

import (
	"Spotify-Visualizer/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var maxX, maxY int
var recentX, recentY int
var currentAccessToken string
var currentTrackID string
var trackResponse models.TrackResponse
var audioFeatures models.AudioFeatures
var out *chan []byte


func SetPipeline(cptr *chan []byte) {
	out = cptr
}

func SetXAndY(canvasX, canvasY int) {
	maxX = canvasX
	maxY = canvasY
}

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
		trackID, ok := GetCurrentTrackID()
		if ok {
			// TODO: send a notification to the frontend that user is currently playing a song
			currentTrackID = trackID
			SetCurrentAudioFeaturesOfTrack()
			computeNextCoordinatesFromSongInfo()

		}
		time.Sleep(time.Second * time.Duration(timeout))
	}
}

func GetCurrentTrackID() (string, bool){
	client := http.Client{}
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		fmt.Printf("Error while trying to create a new request in getCurrentTrackID: %v", err)
		return "", false
	}
	request.Header.Set("Authorization", "Bearer " + currentAccessToken)
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
	err = json.Unmarshal(body, &trackResponse)
	if err != nil {
		fmt.Println("Error while trying to unmarshal response body in getCurrentTrackID: ", err)
		return "", false
	}
	fmt.Println(string(body))
	if !trackResponse.IsPlaying {
		return "", false
	}

	return trackResponse.ID, true
}

func SetCurrentAudioFeaturesOfTrack() {
	client := http.Client{}
	trackURL := fmt.Sprintf("https://api.spotify.com/v1/audio-features/%s", currentTrackID)
	request, err := http.NewRequest("GET", trackURL, nil)
	if err != nil {
		fmt.Printf("Error while trying to create a new request in SetCurrentAudioFeaturesOfTrack: %v", err)
		return
	}
	request.Header.Set("Authorization", "Bearer " + currentAccessToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while trying to send request in SetCurrentAudioFeaturesOfTrack: ", err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body in SetCurrentAudioFeaturesOfTrack: ", err)
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &audioFeatures)
}

// TODO: Hier hab ich wahrsch auch einen while loop wodrin ich mir dann eine calculatete anzahl ein punkten die in eine bestimmte richtung
// TODO: gehen ans frontend zurücksende, bevor ich dann die richtung wechsel und sich der prozess wiederholt. Am Anfang
// TODO: des while loops hol ich mir dann immer wieder neu werte von der trackID die sich geändert haben könnte und passe
// TODO: dadurch dann ggf. den ouput den ich zurücksende ans frontend an (neue kreisgroeße, farbe, form etc)
// Y: 829 – X: 1680
func computeNextCoordinatesFromSongInfo() {
	colorPalette := getColorForCurrentTrack()
	ellipseWidth, ellipseHeight := getEllipseWidthHeight()
	stepRange := getStepRange()

	rand.Seed(time.Now().UnixNano())
	numberOfSteps := rand.Intn(stepRange[1] - stepRange[0]) + stepRange[0]
	currentDirection := rand.Intn(7)

	for numberOfSteps >= 0 && isPositionOnCanvas(positionAfterStep(numberOfSteps)) {
		randomColorIndex := rand.Intn(len(colorPalette))
		currentColor := colorPalette[randomColorIndex]


		numberOfSteps--
	}

}

func getColorForCurrentTrack() [5]models.RGB {
	switch EnergyAndDanceability := audioFeatures.Danceability + audioFeatures.Energy; {
	case EnergyAndDanceability > 1.4:
		currentColors := models.FunkyColors
		return [5]models.RGB{
			currentColors.FunkyOrange,
			currentColors.FunkyDarkBlue,
			currentColors.FunkyLightBlue,
			currentColors.FunkyLightGreen,
			currentColors.FunkyDarkGreen,
		}
	case EnergyAndDanceability > 1.1:
		currentColors := models.WarmColors
		return [5]models.RGB{
			currentColors.WarmRed,
			currentColors.WarmOrange,
			currentColors.WarmPink,
			currentColors.WarmMagenta,
			currentColors.WarmYellow,
		}
	case EnergyAndDanceability > 1.0:
		currentColors := models.ConfidentColors
		return [5]models.RGB{
			currentColors.ConfidentPink,
			currentColors.ConfidentMagenta,
			currentColors.ConfidentRed,
			currentColors.ConfidentOrange,
			currentColors.ConfidentYellow,
		}
	case EnergyAndDanceability > 0.8:
		currentColors := models.RelaxedColors
		return [5]models.RGB{
			currentColors.RelaxedPink,
			currentColors.RelaxedYellow,
			currentColors.RelaxedBrown,
			currentColors.RelaxedOrange,
			currentColors.RelaxedBlue,
		}
	default:
		currentColors := models.SadColors
		return [5]models.RGB{
			currentColors.SadDarkOrange,
			currentColors.SadLightOrange,
			currentColors.SadGreen,
			currentColors.SadRed,
			currentColors.SadBlue,
		}
	}
}

func getEllipseWidthHeight() (float64, float64) {
	value := audioFeatures.Liveness * 10

	return value, value
}

func getStepRange() [2]int {
	switch {
	case audioFeatures.Tempo > 150:
		return [2]int{1, 4}
	case audioFeatures.Tempo > 140:
		return [2]int{2, 7}
	case audioFeatures.Tempo > 120:
		return [2]int{3, 9}
	case audioFeatures.Tempo > 100:
		return [2]int{4, 13}
	case audioFeatures.Tempo > 80:
		return [2]int{8, 20}
	case audioFeatures.Tempo > 60:
		return [2]int{10, 25}
	default:
		return [2]int{11, 27}
	}
}

func positionAfterStep(numberOfSteps int) (int, int) {

}

func isPositionOnCanvas(x, y int) bool {

}

