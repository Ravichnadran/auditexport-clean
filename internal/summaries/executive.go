package summaries

import (
	"auditexport/internal/run"
	"os"
	"time"
)

func WriteExecutiveSummary() error {
	content := `
AuditExport — Executive Summary
================================

This evidence package was generated automatically using AuditExport.

Purpose:
--------
To provide immutable, read-only, audit-ready technical evidence
for ISO/IEC 27001 and SOC 2 compliance audits.

Scope Included:
---------------
- GitHub organization configuration
- Repository governance controls
- Branch protection & merge enforcement
- Access control & contributor visibility
- Change history (commits & pull requests)

Key Properties:
---------------
- No production data collected
- No secrets stored
- Read-only API access
- Locally generated evidence
- Hash-verified integrity

Generated At:
-------------
` + time.Now().UTC().Format(time.RFC3339) + `

Prepared for:
-------------
External auditors, compliance consultants, and internal GRC teams.
`

	return os.WriteFile(
		run.EvidencePath("summaries", "executive_summary.txt"),
		[]byte(content),
		0644,
	)
}
