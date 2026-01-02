package summaries

import (
	"auditexport/internal/run"
	"os"
	"time"
)

func WriteGitHubSummary() error {
	content := `
GitHub Technical Evidence Summary
================================

This report summarizes governance and access controls
observed within the GitHub environment.

Evidence Collected:
-------------------
- Organization metadata
- Repository inventory
- Branch structures
- Commit history
- Pull request workflow
- Contributor visibility
- Access permissions
- Protected branch enforcement
- Code ownership rules
- Required review policies
- Merge restrictions
- Audit log configuration
- Retention policies

Audit Relevance:
----------------
ISO/IEC 27001:
- A.8  Asset Management
- A.9  Access Control
- A.12 Operations Security
- A.14 Secure Development

SOC 2:
- CC6 Logical Access
- CC7 Change Management
- CC8 System Operations

Collection Method:
------------------
Read-only GitHub APIs.
No mutations or write operations performed.

Generated At:
-------------
` + time.Now().UTC().Format(time.RFC3339)

	return os.WriteFile(
		run.EvidencePath("summaries", "github_summary.txt"),
		[]byte(content),
		0644,
	)
}
