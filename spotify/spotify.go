package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetAccessToken(refreshToken string) (string, error) {
	client := http.Client{}
	body, err := json.Marshal(map[string]string{
		"grant_type": "refresh_token",
		"refresh_token": refreshToken,
	})
	if err != nil {
		fmt.Printf("Error while trying to parse body data for request: %v", err)
		return "", err
	}

	encodedAuthHeader, err := getEncodedAuthHeader(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET"))
	if err != nil {
		fmt.Println("Error while trying to encode auth header: ", err)
	}
	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error while trying to create a new request: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic " + encodedAuthHeader)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while trying to send request: ", err)
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while trying to read response body: ", err)
	}

	return fmt.Sprint(string(body)), nil
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
