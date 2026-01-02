package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
	"time"
)

type Repository struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	Private       bool      `json:"private"`
	HTMLURL       string    `json:"html_url"`
	DefaultBranch string    `json:"default_branch"`
	Archived      bool      `json:"archived"`
	CollectedAt   time.Time `json:"collected_at"`
}

type RepositoriesOutput struct {
	GeneratedAt  time.Time    `json:"generated_at"`
	Total        int          `json:"total"`
	Repositories []Repository `json:"repositories"`
}

func WriteRepositories() error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	req, err := client.newRequest(
		"GET",
		baseURL+"/user/repos?per_page=100",
	)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiRepos []struct {
		ID            int64  `json:"id"`
		Name          string `json:"name"`
		FullName      string `json:"full_name"`
		Private       bool   `json:"private"`
		HTMLURL       string `json:"html_url"`
		DefaultBranch string `json:"default_branch"`
		Archived      bool   `json:"archived"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiRepos); err != nil {
		return err
	}

	output := RepositoriesOutput{
		GeneratedAt: time.Now().UTC(),
	}

	for _, r := range apiRepos {
		output.Repositories = append(output.Repositories, Repository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Private,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.DefaultBranch,
			Archived:      r.Archived,
			CollectedAt:   time.Now().UTC(),
		})
	}

	output.Total = len(output.Repositories)

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "repositories.json"),
		data,
		0644,
	)
}
