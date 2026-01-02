package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Organization struct {
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CollectedAt time.Time `json:"collected_at"`
}

func WriteOrganization() error {
	client, err := NewClient()
	if err != nil {
		return err
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
		return fmt.Errorf("github api error: %s", resp.Status)
	}

	var org Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return err
	}

	org.CollectedAt = time.Now().UTC()

	data, err := json.MarshalIndent(org, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "organization.json"),
		data,
		0644,
	)
}
