package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
)

type RepositoryRef struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type repositoriesEvidence struct {
	Repositories []RepositoryRef `json:"repositories"`
}

// LoadRepositoriesFromEvidence reads github/repositories.json
// and returns repository name + full_name pairs.
// This guarantees NO API calls and preserves audit immutability.
func LoadRepositoriesFromEvidence() ([]RepositoryRef, error) {
	bytes, err := os.ReadFile(
		run.EvidencePath("github", "repositories.json"),
	)
	if err != nil {
		return nil, err
	}

	var model repositoriesEvidence
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	return model.Repositories, nil
}
