package mapping

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
	"time"
)

func WriteControlMapping() error {

	mapping := map[string][]string{
		"CC6.3": {
			"github/access_controls.json",
		},
		"CC7.2": {
			"github/commits.json",
		},
		"CC8.1": {
			"github/pull_requests.json",
			"github/workflows/workflow_runs.json",
			"summaries/soc2_change_management_assertions.txt",
		},
	}

	data, _ := json.MarshalIndent(struct {
		GeneratedAt time.Time           `json:"generated_at"`
		Controls    map[string][]string `json:"controls"`
	}{
		GeneratedAt: time.Now().UTC(),
		Controls:    mapping,
	}, "", "  ")

	return os.WriteFile(
		run.EvidencePath("summaries", "soc2_control_mapping.json"),
		data,
		0644,
	)
}
