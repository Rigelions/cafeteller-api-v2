package instagram

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type Profile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func GetMe(accessToken string) (Profile, error) {
	profileURI := "https://graph.instagram.com/me?fields=id,username&access_token=" + accessToken
	resp, err := http.Get(profileURI)
	if err != nil {
		log.Println("Error making request to Instagram:", err)
		return Profile{}, errors.New("failed to make request to Instagram")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return Profile{}, errors.New("failed to read response from Instagram")
	}

	log.Println("Response from Instagram:", string(body))

	var profile Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Println("Error unmarshalling response:", err)
		return Profile{}, errors.New("failed to parse response from Instagram")
	}

	return profile, nil
}
