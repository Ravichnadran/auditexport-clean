package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Branch struct {
	Repository  string    `json:"repository"`
	Name        string    `json:"name"`
	Protected   bool      `json:"protected"`
	CommitSHA   string    `json:"commit_sha"`
	CollectedAt time.Time `json:"collected_at"`
}

type BranchesOutput struct {
	CollectedAt time.Time `json:"collected_at"`
	Total       int       `json:"total"`
	Branches    []Branch  `json:"branches"`
}

func WriteBranches() error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	// Load repositories already collected
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

	var allBranches []Branch

	for _, repo := range repoData.Repositories {
		url := fmt.Sprintf(
			"%s/repos/%s/branches?per_page=100",
			baseURL,
			repo.FullName,
		)

		req, err := client.newRequest("GET", url)
		if err != nil {
			return err
		}

		resp, err := client.httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			return fmt.Errorf(
				"branches API error for %s: %s",
				repo.FullName,
				resp.Status,
			)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}

		var apiBranches []struct {
			Name      string `json:"name"`
			Protected bool   `json:"protected"`
			Commit    struct {
				SHA string `json:"sha"`
			} `json:"commit"`
		}

		if err := json.Unmarshal(body, &apiBranches); err != nil {
			return err
		}

		for _, b := range apiBranches {
			allBranches = append(allBranches, Branch{
				Repository:  repo.FullName,
				Name:        b.Name,
				Protected:   b.Protected,
				CommitSHA:   b.Commit.SHA,
				CollectedAt: time.Now().UTC(),
			})
		}
	}

	output := BranchesOutput{
		CollectedAt: time.Now().UTC(),
		Total:       len(allBranches),
		Branches:    allBranches,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "branches.json"),
		data,
		0644,
	)
}
