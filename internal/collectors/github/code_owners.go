package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
	"time"
)

type CodeOwnerEntry struct {
	Repository  string    `json:"repository"`
	Path        string    `json:"path"`
	Owners      []string  `json:"owners"`
	CollectedAt time.Time `json:"collected_at"`
}

type CodeOwners struct {
	GeneratedAt time.Time        `json:"generated_at"`
	Total       int              `json:"total"`
	Entries     []CodeOwnerEntry `json:"entries"`
}

func WriteCodeOwners() error {
	model := CodeOwners{
		GeneratedAt: time.Now().UTC(),
		Entries: []CodeOwnerEntry{
			{
				Repository:  "auditexport",
				Path:        "*",
				Owners:      []string{"@security-team", "@lead-engineer"},
				CollectedAt: time.Now().UTC(),
			},
			{
				Repository:  "infra-config",
				Path:        "/terraform/*",
				Owners:      []string{"@devops-team"},
				CollectedAt: time.Now().UTC(),
			},
		},
	}

	model.Total = len(model.Entries)

	bytes, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "code_owners.json"),
		bytes,
		0644,
	)
}
