package cicd

import (
	"auditexport/internal/run"
	"os"
)

func WriteCISummary() error {
	content := `CI/CD Evidence Summary

- Platform: GitHub Actions
- Pipelines are defined as code
- Pipelines are version controlled
- Pipelines trigger on pull requests and merges
- Secrets are stored in platform-managed secret store
- Build execution is enforced before merge
`

	return os.WriteFile(
		run.EvidencePath("summaries", "ci_summary.txt"),
		[]byte(content),
		0644,
	)
}
