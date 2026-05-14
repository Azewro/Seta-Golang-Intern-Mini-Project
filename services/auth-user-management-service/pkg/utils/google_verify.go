package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TokenInfoResponse struct {
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Error         string `json:"error"`
	ErrorDesc     string `json:"error_description"`
}

func VerifyGoogleIDToken(idToken string) (*TokenInfoResponse, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token, status code: %d", resp.StatusCode)
	}

	var tokenInfo TokenInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if tokenInfo.Error != "" {
		return nil, fmt.Errorf("google error: %s", tokenInfo.ErrorDesc)
	}

	return &tokenInfo, nil
}

