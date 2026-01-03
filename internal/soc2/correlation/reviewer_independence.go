package correlation

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
	"time"
)

type ReviewerViolation struct {
	Repository string    `json:"repository"`
	PRNumber   int       `json:"pr_number"`
	Author     string    `json:"author"`
	MergedBy   string    `json:"merged_by"`
	MergedAt   time.Time `json:"merged_at"`
}

// CheckReviewerIndependence ensures PR author != merger.
func CheckReviewerIndependence(from, to time.Time) ([]ReviewerViolation, error) {

	bytes, err := os.ReadFile(
		run.EvidencePath("github", "pull_requests.json"),
	)
	if err != nil {
		return nil, err
	}

	var model struct {
		PullRequests []struct {
			Repository string     `json:"repository"`
			Number     int        `json:"number"`
			Author     string     `json:"author"`
			Merged     bool       `json:"merged"`
			MergedBy   string     `json:"merged_by"`
			MergedAt   *time.Time `json:"merged_at"`
		} `json:"pull_requests"`
	}

	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	var violations []ReviewerViolation

	for _, pr := range model.PullRequests {

		if !pr.Merged || pr.MergedAt == nil {
			continue
		}

		if pr.MergedAt.Before(from) || pr.MergedAt.After(to) {
			continue
		}

		// ❌ Self-approval violation
		if pr.Author == pr.MergedBy && pr.MergedBy != "" {
			violations = append(violations, ReviewerViolation{
				Repository: pr.Repository,
				PRNumber:   pr.Number,
				Author:     pr.Author,
				MergedBy:   pr.MergedBy,
				MergedAt:   *pr.MergedAt,
			})
		}
	}

	return violations, nil
}
