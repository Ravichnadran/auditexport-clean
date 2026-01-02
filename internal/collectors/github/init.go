package github

import (
	"auditexport/internal/run"
	"os"
)

func InitGithubDir() error {
	run.InitRunContext()
	return os.MkdirAll(
		run.EvidencePath("github"),
		0755,
	)
}
