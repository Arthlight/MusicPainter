package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
	"strings"
)

func GetAccessToken(refreshToken string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"grant_type": "refresh_token",
		"refresh_token": refreshToken,
	})
	if err != nil {
		fmt.Printf("Error while trying to parse body data for request: %v", err)
		return "", err
	}

	AuthHeaderString := os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")

	var sb strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &sb)
	_, err = encoder.Write([]byte(AuthHeaderString))
	if err != nil {
		fmt.Println("Error encoding AuthHeaderString into base64: ", err)
		return "", nil
	}
	err = encoder.Close()
	if err != nil {
		fmt.Println("Error closing base64 encoder: ", err)
		return "", nil
	}
	encodedAuthHeaderString := sb.String()

	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error while trying to create a new request: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic " + encodedAuthHeaderString)
}
