package summaries

import (
	"auditexport/internal/run"
	"os"
)

func WriteAuditorNotes() error {
	content := `
Auditor Notes
=============

- Evidence was collected using read-only access.
- No production data or source code was exported.
- Evidence reflects configuration state at time of execution.
- Hashes provided to verify integrity.
- ZIP archive may be stored as-is for audit records.

If additional clarification is required,
contact the system owner or re-run AuditExport.
`

	return os.WriteFile(
		run.EvidencePath("summaries", "auditor_notes.txt"),
		[]byte(content),
		0644,
	)
}
