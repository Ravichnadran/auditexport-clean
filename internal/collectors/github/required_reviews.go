package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
	"time"
)

type RequiredReviews struct {
	Repository    string    `json:"repository"`
	Branch        string    `json:"branch"`
	Enabled       bool      `json:"enabled"`
	RequiredCount int       `json:"required_approving_review_count,omitempty"`
	CodeOwners    bool      `json:"require_code_owner_reviews,omitempty"`
	DismissStale  bool      `json:"dismiss_stale_reviews,omitempty"`
	CollectedAt   time.Time `json:"collected_at"`
	Status        string    `json:"status"`
}

func WriteRequiredReviews(owner, repo, branch string) error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	url := baseURL +
		"/repos/" + owner + "/" + repo +
		"/branches/" + branch +
		"/protection/required_pull_request_reviews"

	req, err := client.newRequest("GET", url)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := RequiredReviews{
		Repository:  owner + "/" + repo,
		Branch:      branch,
		CollectedAt: time.Now().UTC(),
	}

	// ❗ Expected cases: 404 / 403
	if resp.StatusCode != 200 {
		result.Enabled = false
		result.Status = "not_configured_or_no_access"

		data, _ := json.MarshalIndent(result, "", "  ")
		return os.WriteFile(
			run.EvidencePath("github", "required_reviews.json"),
			data,
			0644,
		)
	}

	var apiResp struct {
		RequiredApprovingReviewCount int  `json:"required_approving_review_count"`
		RequireCodeOwnerReviews      bool `json:"require_code_owner_reviews"`
		DismissStaleReviews          bool `json:"dismiss_stale_reviews"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}

	result.Enabled = true
	result.RequiredCount = apiResp.RequiredApprovingReviewCount
	result.CodeOwners = apiResp.RequireCodeOwnerReviews
	result.DismissStale = apiResp.DismissStaleReviews
	result.Status = "configured"

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "required_reviews.json"),
		data,
		0644,
	)
}
