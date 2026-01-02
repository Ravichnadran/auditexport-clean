package github

import (
	"fmt"
	"net/http"
)

// ValidateAuth ensures GITHUB_TOKEN exists AND is valid
func ValidateAuth() error {
	client, err := NewClient()
	if err != nil {
		return err // handles "GITHUB_TOKEN not set"
	}

	req, err := client.newRequest("GET", baseURL+"/user")
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"GITHUB_TOKEN invalid or unauthorized (HTTP %d)",
			resp.StatusCode,
		)
	}

	return nil
}
