package run

import "os"

func InitEvidenceRunDir() error {
	InitRunContext()
	return os.MkdirAll(
		EvidencePath("run"),
		0755,
	)
}

func InitSummariesDir() error {
	InitRunContext()
	return os.MkdirAll(
		EvidencePath("summaries"),
		0755,
	)
}
