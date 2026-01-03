package correlation

import (
	"auditexport/internal/run"
	"encoding/json"
	"errors"
	"os"
	"time"
)

// BuildChangeEvents correlates PRs → CI runs → merge events.
// This is the PRIMARY SOC 2 CC8.1 evidence generator.
func BuildChangeEvents(from, to time.Time) ([]ChangeEvent, error) {

	// --------------------------------------------------
	// Load Pull Requests
	// --------------------------------------------------
	prBytes, err := os.ReadFile(
		run.EvidencePath("github", "pull_requests.json"),
	)
	if err != nil {
		return nil, err
	}

	var prModel struct {
		PullRequests []struct {
			Repository string     `json:"repository"`
			Number     int        `json:"number"`
			Author     string     `json:"author"`
			Merged     bool       `json:"merged"`
			MergedAt   *time.Time `json:"merged_at"`
		} `json:"pull_requests"`
	}

	if err := json.Unmarshal(prBytes, &prModel); err != nil {
		return nil, err
	}

	// --------------------------------------------------
	// Load Workflow Runs
	// --------------------------------------------------
	runBytes, err := os.ReadFile(
		run.EvidencePath("github", "workflows", "workflow_runs.json"),
	)
	if err != nil {
		return nil, errors.New("SOC2 requires CI workflow runs")
	}

	var workflowRuns []struct {
		ID         int       `json:"id"`
		Name       string    `json:"name"`
		Conclusion string    `json:"conclusion"`
		CreatedAt  time.Time `json:"created_at"`
	}

	if err := json.Unmarshal(runBytes, &workflowRuns); err != nil {
		return nil, err
	}

	// --------------------------------------------------
	// Correlate PR → CI → Merge
	// --------------------------------------------------
	var events []ChangeEvent

	for _, pr := range prModel.PullRequests {

		// Only merged PRs count
		if !pr.Merged || pr.MergedAt == nil {
			continue
		}

		mergedAt := *pr.MergedAt

		// Apply audit window
		if mergedAt.Before(from) || mergedAt.After(to) {
			continue
		}

		var ciEvidence []CIRunEvidence

		for _, run := range workflowRuns {
			// CI must complete BEFORE merge
			if run.CreatedAt.After(mergedAt) {
				continue
			}

			ciEvidence = append(ciEvidence, CIRunEvidence{
				RunID:      run.ID,
				Name:       run.Name,
				Conclusion: run.Conclusion,
				CreatedAt:  run.CreatedAt,
			})
		}

		// SOC 2 requires AT LEAST ONE SUCCESSFUL CI
		hasSuccess := false
		for _, ci := range ciEvidence {
			if ci.Conclusion == "success" {
				hasSuccess = true
				break
			}
		}

		if !hasSuccess {
			continue // ❌ reject non-compliant merge
		}

		events = append(events, ChangeEvent{
			Repository: pr.Repository,
			PRNumber:   pr.Number,
			Author:     pr.Author,
			MergedAt:   mergedAt,
			CIRuns:     ciEvidence,
		})
	}

	return events, nil
}
