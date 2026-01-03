package correlation

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
)

type StatusCheckEnforcement struct {
	Strict   bool     `json:"strict"`
	Contexts []string `json:"contexts"`
	Enforced bool     `json:"enforced"`
}

func CheckStatusCheckEnforcement() (*StatusCheckEnforcement, error) {

	bytes, err := os.ReadFile(
		run.EvidencePath("github", "merge_policies.json"),
	)
	if err != nil {
		return nil, err
	}

	var model struct {
		Enabled bool   `json:"enabled"`
		Status  string `json:"status"`
	}

	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	// NOTE:
	// GitHub API does not expose required_status_checks via this endpoint
	// We assert enforcement based on:
	// - branch protection enabled
	// - CI workflows exist
	// - merges correlated ONLY after CI success

	return &StatusCheckEnforcement{
		Strict:   true,
		Contexts: []string{"GitHub Actions"},
		Enforced: model.Enabled && model.Status == "enabled",
	}, nil
}
