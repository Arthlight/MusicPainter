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
var recentX, recentY = 30, 30
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
			currentTrackID = trackID
			notifyFrontend()
			SetCurrentAudioFeaturesOfTrack()
			sendNextCoordinatesFromSongInfoToFrontend()
		}
		time.Sleep(time.Second * time.Duration(timeout))
	}
}

func notifyFrontend() {
	isPlaying := map[string]bool{"isPlaying": true}
	frontendData := &models.Event{
		Name:    "isPlaying",
		Content: isPlaying,
	}
	rawFrontendData, err := frontendData.ToBinary()
	if err != nil {
		fmt.Println("Error while trying to marshal frontendData in notifyFrontend: ", err)
	}
	*out <- rawFrontendData
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

func sendNextCoordinatesFromSongInfoToFrontend() {
	colorPalette := getColorForCurrentTrack()
	ellipseWidth, ellipseHeight := getEllipseWidthHeight()
	stepRange := getStepRange()
	stepSize := getRandomNumInRange(int(audioFeatures.Valence), int(2 * audioFeatures.Valence))
	numberOfSteps := getRandomNumInRange(stepRange[0], stepRange[1])
	currentDirection := getRandomNumInRange(0, 7)

	for numberOfSteps >= 0 && isPositionOnCanvas(positionAfterStep(stepSize, currentDirection)) {
		randomColorIndex := rand.Intn(len(colorPalette))
		currentColor := colorPalette[randomColorIndex]
		err := sendToFrontend(currentColor, stepSize, int(ellipseHeight), int(ellipseWidth), recentX, recentY, trackResponse.Name)
		if err != nil {
			fmt.Println("Error while trying to send data to frontend: ", err)
		}
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
		return [2]int{2, 4}
	case audioFeatures.Tempo > 140:
		return [2]int{3, 7}
	case audioFeatures.Tempo > 120:
		return [2]int{4, 10}
	case audioFeatures.Tempo > 100:
		return [2]int{5, 13}
	case audioFeatures.Tempo > 80:
		return [2]int{8, 20}
	case audioFeatures.Tempo > 60:
		return [2]int{10, 25}
	default:
		return [2]int{11, 27}
	}
}

func positionAfterStep(stepSize, currentDirection int) []int {
	switch currentDirection {
	case 0:
		// go up
		recentY = recentY - stepSize
		return []int{recentX, recentY}
	case 1:
		// go up-right
		recentX = recentX + stepSize
		recentY = recentY - stepSize
		return []int{recentX, recentY}
	case 2:
		// go up-left
		recentX = recentX - stepSize
		recentY = recentY - stepSize
		return []int{recentX, recentY}
	case 3:
		// go left
		recentX = recentX - stepSize
		return []int{recentX, recentY}
	case 4:
		// go right
		recentX = recentX + stepSize
		return []int{recentX, recentY}
	case 5:
		// go down-right
		recentX = recentX + stepSize
		recentY = recentY + stepSize
		return []int{recentX, recentY}
	case 6:
		// go down-left
		recentX = recentX - stepSize
		recentY = recentY + stepSize
		return []int{recentX, recentY}
	case 7:
		// go down
		recentY = recentY + stepSize
		return []int{recentX, recentY}
	default:
		panic("Invalid case in switch statement")
	}
}

// Y: 829 | X: 1680
func isPositionOnCanvas(coordinates []int) bool {
	currentX, currentY := coordinates[0], coordinates[1]
	return currentX < maxX && currentX >= 0 && currentY < maxY && currentY >= 0
}

func getRandomNumInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min) + min
}

func sendToFrontend(currentColor models.RGB, stepSize int, ellipseHeight int, ellipseWidth int, x int, y int, songName string) error {
	dataPackage := models.ForFrontend{
		StepSize:      stepSize,
		ColorPalette:  currentColor,
		EllipseHeight: ellipseHeight,
		EllipseWidth:  ellipseWidth,
		X:             x,
		Y:             y,
		SongName:      songName,
	}
	frontendData := &models.Event{
		Name:    "data",
		Content: dataPackage,
	}
	rawFrontendData, err := frontendData.ToBinary()
	if err != nil {
		return err
	}
	*out <- rawFrontendData

	return nil
}

