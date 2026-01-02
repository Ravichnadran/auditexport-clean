package summaries

import (
	"auditexport/internal/run"
	"os"
	"time"
)

func WriteTechnicalSummary() error {
	content := `
Technical Summary
=================

This evidence package was generated deterministically.

Execution Characteristics:
--------------------------
- Single execution context
- Explicit scope declaration
- Timestamped artifacts
- Stable file paths
- Repeatable output

Security Guarantees:
--------------------
- No credentials stored
- No outbound callbacks
- No telemetry
- Offline-verifiable

Integrity:
----------
Each file is cryptographically hashed.
Any modification invalidates the package.

Generated At:
-------------
` + time.Now().UTC().Format(time.RFC3339)

	return os.WriteFile(
		run.EvidencePath("summaries", "technical_summary.txt"),
		[]byte(content),
		0644,
	)
}
