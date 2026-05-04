package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type TeamResponse struct {
	TeamID    uint   `json:"teamId"`
	TeamName  string `json:"teamName"`
	Managers  []uint `json:"managers"`
	Members   []uint `json:"members"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type listTeamsResponse struct {
	Data []TeamResponse `json:"data"`
}

type TeamClient interface {
	ListMyTeams(token string) ([]TeamResponse, error)
	IsManagerOf(token string, managerID uint, memberID uint) (bool, error)
}

type teamClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewTeamClient() TeamClient {
	baseURL := os.Getenv("TEAM_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	return &teamClientImpl{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *teamClientImpl) ListMyTeams(token string) ([]TeamResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/teams/my", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("team service returned status %d", resp.StatusCode)
	}

	var data listTeamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (c *teamClientImpl) IsManagerOf(token string, managerID uint, memberID uint) (bool, error) {
	teams, err := c.ListMyTeams(token)
	if err != nil {
		return false, err
	}

	for i := range teams {
		if containsUint(teams[i].Managers, managerID) && containsUint(teams[i].Members, memberID) {
			return true, nil
		}
	}

	return false, nil
}

func containsUint(values []uint, target uint) bool {
	for i := range values {
		if values[i] == target {
			return true
		}
	}
	return false
}
