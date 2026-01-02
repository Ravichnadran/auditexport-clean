package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type BranchProtection struct {
	Repository        string    `json:"repository"`
	Branch            string    `json:"branch"`
	Protected         bool      `json:"protected"`
	RequiredReviews   bool      `json:"required_reviews"`
	AdminEnforced     bool      `json:"admin_enforced"`
	ForcePushDisabled bool      `json:"force_push_disabled"`
	DeletionDisabled  bool      `json:"deletion_disabled"`
	CollectedAt       time.Time `json:"collected_at"`
}

type ProtectedBranchesOutput struct {
	CollectedAt time.Time          `json:"collected_at"`
	Total       int                `json:"total"`
	Branches    []BranchProtection `json:"branches"`
}

func WriteProtectedBranches() error {
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

	var results []BranchProtection

	for _, repo := range repoData.Repositories {
		page := 1

		for {
			branchesURL := fmt.Sprintf(
				"%s/repos/%s/branches?per_page=100&page=%d",
				baseURL,
				repo.FullName,
				page,
			)

			req, err := client.newRequest("GET", branchesURL)
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
					"branches API error for %s",
					repo.FullName,
				)
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}

			var branches []struct {
				Name      string `json:"name"`
				Protected bool   `json:"protected"`
			}

			if err := json.Unmarshal(body, &branches); err != nil {
				return err
			}

			if len(branches) == 0 {
				break
			}

			for _, b := range branches {
				entry := BranchProtection{
					Repository:  repo.FullName,
					Branch:      b.Name,
					Protected:   b.Protected,
					CollectedAt: time.Now().UTC(),
				}

				if b.Protected {
					protectURL := fmt.Sprintf(
						"%s/repos/%s/branches/%s/protection",
						baseURL,
						repo.FullName,
						b.Name,
					)

					pReq, err := client.newRequest("GET", protectURL)
					if err != nil {
						return err
					}

					pResp, err := client.httpClient.Do(pReq)
					if err != nil {
						return err
					}

					if pResp.StatusCode == 200 {
						pBody, err := io.ReadAll(pResp.Body)
						pResp.Body.Close()
						if err != nil {
							return err
						}

						var protection struct {
							RequiredPullRequestReviews interface{} `json:"required_pull_request_reviews"`
							EnforceAdmins              struct {
								Enabled bool `json:"enabled"`
							} `json:"enforce_admins"`
							AllowForcePushes struct {
								Enabled bool `json:"enabled"`
							} `json:"allow_force_pushes"`
							AllowDeletions struct {
								Enabled bool `json:"enabled"`
							} `json:"allow_deletions"`
						}

						if err := json.Unmarshal(pBody, &protection); err != nil {
							return err
						}

						entry.RequiredReviews = protection.RequiredPullRequestReviews != nil
						entry.AdminEnforced = protection.EnforceAdmins.Enabled
						entry.ForcePushDisabled = !protection.AllowForcePushes.Enabled
						entry.DeletionDisabled = !protection.AllowDeletions.Enabled
					} else {
						pResp.Body.Close()
					}
				}

				results = append(results, entry)
			}

			page++
		}
	}

	output := ProtectedBranchesOutput{
		CollectedAt: time.Now().UTC(),
		Total:       len(results),
		Branches:    results,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "protected_branches.json"),
		data,
		0644,
	)
}
