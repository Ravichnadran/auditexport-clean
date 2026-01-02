package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type PullRequest struct {
	Repository  string     `json:"repository"`
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	State       string     `json:"state"`
	Author      string     `json:"author"`
	Merged      bool       `json:"merged"`
	MergedBy    string     `json:"merged_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	MergedAt    *time.Time `json:"merged_at,omitempty"`
	CollectedAt time.Time  `json:"collected_at"`
}

type PullRequestsOutput struct {
	CollectedAt  time.Time     `json:"collected_at"`
	Total        int           `json:"total"`
	PullRequests []PullRequest `json:"pull_requests"`
}

func WritePullRequests() error {
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

	var allPRs []PullRequest

	for _, repo := range repoData.Repositories {
		page := 1

		for {
			url := fmt.Sprintf(
				"%s/repos/%s/pulls?state=all&per_page=100&page=%d",
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
					"pull requests API error for %s: %s",
					repo.FullName,
					resp.Status,
				)
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}

			var apiPRs []struct {
				Number int    `json:"number"`
				Title  string `json:"title"`
				State  string `json:"state"`
				User   struct {
					Login string `json:"login"`
				} `json:"user"`
				CreatedAt time.Time  `json:"created_at"`
				MergedAt  *time.Time `json:"merged_at"`
				URL       string     `json:"url"`
			}

			if err := json.Unmarshal(body, &apiPRs); err != nil {
				return err
			}

			if len(apiPRs) == 0 {
				break
			}

			for _, pr := range apiPRs {
				// Fetch PR details (for merged_by)
				reqDetail, err := client.newRequest("GET", pr.URL)
				if err != nil {
					return err
				}

				respDetail, err := client.httpClient.Do(reqDetail)
				if err != nil {
					return err
				}

				detailBody, err := io.ReadAll(respDetail.Body)
				respDetail.Body.Close()
				if err != nil {
					return err
				}

				var detail struct {
					Merged   bool `json:"merged"`
					MergedBy *struct {
						Login string `json:"login"`
					} `json:"merged_by"`
				}

				_ = json.Unmarshal(detailBody, &detail)

				var mergedBy string
				if detail.MergedBy != nil {
					mergedBy = detail.MergedBy.Login
				}

				allPRs = append(allPRs, PullRequest{
					Repository:  repo.FullName,
					Number:      pr.Number,
					Title:       pr.Title,
					State:       pr.State,
					Author:      pr.User.Login,
					Merged:      detail.Merged,
					MergedBy:    mergedBy,
					CreatedAt:   pr.CreatedAt,
					MergedAt:    pr.MergedAt,
					CollectedAt: time.Now().UTC(),
				})
			}

			page++
		}
	}

	output := PullRequestsOutput{
		CollectedAt:  time.Now().UTC(),
		Total:        len(allPRs),
		PullRequests: allPRs,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "pull_requests.json"),
		data,
		0644,
	)
}
