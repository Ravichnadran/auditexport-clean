package run

import (
	"encoding/json"
	"os"
	"time"
)

type RunMetadata struct {
	Tool        string    `json:"tool"`
	Version     string    `json:"version"`
	Standard    string    `json:"standard"`
	RunID       string    `json:"run_id"`
	GeneratedAt time.Time `json:"generated_at"`
}

// These should come from build-time flags later
var (
	ToolName = "AuditExport"
	Version  = "v0.2.0"
)

func WriteRunMetadata(standard string) error {
	meta := RunMetadata{
		Tool:        ToolName,
		Version:     Version,
		Standard:    standard,
		RunID:       "run-001",
		GeneratedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		EvidencePath("run", "run_metadata.json"),
		data,
		0644,
	)
}
