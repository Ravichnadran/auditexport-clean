package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const baseURL = "https://api.github.com"

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient() (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("GITHUB_TOKEN not set")
	}

	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "auditexport")

	return req, nil
}

func (c *Client) GetJSON(url string, target interface{}) error {
	req, err := c.newRequest("GET", url)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &GitHubAPIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("github api error: %s", resp.Status),
		}
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func WriteJSON(path string, data interface{}) error {
	// `path` is always relative to evidence/
	fullPath := run.EvidencePath(path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fullPath, bytes, 0644)
}
