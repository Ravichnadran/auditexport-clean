package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Commit struct {
	Repository  string    `json:"repository"`
	SHA         string    `json:"sha"`
	Message     string    `json:"message"`
	Author      string    `json:"author"`
	Committer   string    `json:"committer"`
	CommittedAt time.Time `json:"committed_at"`
	CollectedAt time.Time `json:"collected_at"`
}

type CommitsOutput struct {
	CollectedAt time.Time `json:"collected_at"`
	Total       int       `json:"total"`
	Commits     []Commit  `json:"commits"`
}

func WriteCommits() error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	repoBytes, err := os.ReadFile(
		run.EvidencePath("github", "repositories.json"),
	)
	if err != nil {
		return err
	}

	var repoData RepositoriesOutput
	if err := json.Unmarshal(repoBytes, &repoData); err != nil {
		return err
	}

	var allCommits []Commit

	for _, repo := range repoData.Repositories {
		page := 1

		for {
			url := fmt.Sprintf(
				"%s/repos/%s/commits?per_page=100&page=%d",
				baseURL,
				repo.FullName,
				page,
			)

			req, err := client.newRequest("GET", url)
			if err != nil {
				return err
			}

			resp, err := client.httpClient.Do(req)
			if err != nil {
				return err
			}

			// ✅ HANDLE EMPTY REPOSITORY (409 Conflict)
			if resp.StatusCode == 409 {
				resp.Body.Close()
				break // valid audit state: no commits
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				return fmt.Errorf(
					"commits API error for %s: %s",
					repo.FullName,
					resp.Status,
				)
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}

			var apiCommits []struct {
				SHA    string `json:"sha"`
				Commit struct {
					Message string `json:"message"`
					Author  struct {
						Name string    `json:"name"`
						Date time.Time `json:"date"`
					} `json:"author"`
					Committer struct {
						Name string    `json:"name"`
						Date time.Time `json:"date"`
					} `json:"committer"`
				} `json:"commit"`
			}

			if err := json.Unmarshal(body, &apiCommits); err != nil {
				return err
			}

			if len(apiCommits) == 0 {
				break
			}

			for _, c := range apiCommits {
				allCommits = append(allCommits, Commit{
					Repository:  repo.FullName,
					SHA:         c.SHA,
					Message:     c.Commit.Message,
					Author:      c.Commit.Author.Name,
					Committer:   c.Commit.Committer.Name,
					CommittedAt: c.Commit.Committer.Date,
					CollectedAt: time.Now().UTC(),
				})
			}

			page++
		}
	}

	output := CommitsOutput{
		CollectedAt: time.Now().UTC(),
		Total:       len(allCommits),
		Commits:     allCommits,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "commits.json"),
		data,
		0644,
	)
}
