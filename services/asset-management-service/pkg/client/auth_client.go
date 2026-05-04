package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code from auth service")
	ErrUnauthorized         = errors.New("unauthorized by auth service")
)

type UserResponse struct {
	ID       uint   `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type bulkUsersResponse struct {
	Data []UserResponse `json:"data"`
}

type AuthClient interface {
	VerifyToken(token string) (*UserResponse, error)
	GetUsers(token string, ids []uint) ([]UserResponse, error)
}

type authClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient() AuthClient {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &authClientImpl{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *authClientImpl) VerifyToken(token string) (*UserResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/users/me", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var data UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

type bulkRequest struct {
	UserIDs []uint `json:"userIds"`
}

func (c *authClientImpl) GetUsers(token string, ids []uint) ([]UserResponse, error) {
	if len(ids) == 0 {
		return []UserResponse{}, nil
	}

	reqBody, _ := json.Marshal(bulkRequest{UserIDs: ids})
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/users/bulk", c.baseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var data bulkUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Data, nil
}
