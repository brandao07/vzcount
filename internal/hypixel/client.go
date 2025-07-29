package hypixel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	headerAPIKey = "API-Key"
)

var (
	ErrInvalidAPIKey     = errors.New("invalid API key")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrUnexpectedStatus  = errors.New("unexpected status code from Hypixel API")
)

type API struct {
	URL string // Base URL for the Hypixel API
	Key string // API key for authentication
}

// VampireZCount retrieves the current player count for the VampireZ game mode
func (hype *API) VampireZCount() (int, error) {
	body, err := hype.playerCount()
	if err != nil {
		return -1, err
	}

	// Parse the JSON response to extract VampireZ player count
	var parsed struct {
		Games struct {
			Legacy struct {
				Modes struct {
					VampireZ int `json:"VAMPIREZ"`
				} `json:"modes"`
			} `json:"LEGACY"`
		} `json:"games"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return -1, err
	}

	return parsed.Games.Legacy.Modes.VampireZ, nil
}

// playerCount makes a request to the Hypixel API /counts endpoint and returns the raw response body
func (hype *API) playerCount() ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", hype.URL+"/counts", nil)

	if err != nil {
		return nil, err
	}

	// Add API key to request headers
	req.Header.Add(headerAPIKey, hype.Key)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Handle different HTTP status codes
	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 403:
			return nil, ErrInvalidAPIKey
		case 429:
			return nil, ErrRateLimitExceeded
		default:
			return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
		}
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
