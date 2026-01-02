package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Contributor struct {
	Repository    string    `json:"repository"`
	Username      string    `json:"username"`
	Contributions int       `json:"contributions"`
	CollectedAt   time.Time `json:"collected_at"`
}

type ContributorsOutput struct {
	CollectedAt  time.Time     `json:"collected_at"`
	Total        int           `json:"total"`
	Contributors []Contributor `json:"contributors"`
}

func WriteContributors() error {
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

	var allContributors []Contributor

	for _, repo := range repoData.Repositories {
		page := 1

		for {
			url := fmt.Sprintf(
				"%s/repos/%s/contributors?per_page=100&page=%d",
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

			if resp.StatusCode != 200 {
				resp.Body.Close()
				return fmt.Errorf(
					"contributors API error for %s: %s",
					repo.FullName,
					resp.Status,
				)
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}

			var apiContributors []struct {
				Login         string `json:"login"`
				Contributions int    `json:"contributions"`
			}

			if err := json.Unmarshal(body, &apiContributors); err != nil {
				return err
			}

			if len(apiContributors) == 0 {
				break
			}

			for _, c := range apiContributors {
				allContributors = append(allContributors, Contributor{
					Repository:    repo.FullName,
					Username:      c.Login,
					Contributions: c.Contributions,
					CollectedAt:   time.Now().UTC(),
				})
			}

			page++
		}
	}

	output := ContributorsOutput{
		CollectedAt:  time.Now().UTC(),
		Total:        len(allContributors),
		Contributors: allContributors,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "contributors.json"),
		data,
		0644,
	)
}
