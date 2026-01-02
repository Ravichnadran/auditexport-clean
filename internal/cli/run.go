package cli

import (
	"fmt"

	"auditexport/internal/collectors/github"
	"auditexport/internal/run"
)

// Run executes the 'run' command.
func Run() error {
	if err := run.InitEvidenceRunDir(); err != nil {
		return err
	}

	_ = run.WriteExecutionLog("run started")

	// GitHub evidence
	if err := github.InitGithubDir(); err != nil {
		return fmt.Errorf("failed to init github dir: %w", err)
	}

	if err := github.WriteOrganization(); err != nil {
		return fmt.Errorf("failed to write github organization: %w", err)
	}

	_ = run.WriteExecutionLog("github organization collected")
	_ = run.WriteExecutionLog("run completed")

	fmt.Println("auditexport run completed")
	return nil
}
