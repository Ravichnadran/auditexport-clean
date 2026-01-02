package cicd

import (
	"fmt"
	"time"

	"auditexport/internal/collectors/github"
)

type WorkflowRun struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Event      string    `json:"event"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	CreatedAt  time.Time `json:"created_at"`
}

type workflowRunResponse struct {
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

func WriteWorkflowRuns(owner, repo string) error {
	client, err := github.NewClient()
	if err != nil {
		return err
	}

	var resp workflowRunResponse
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs?per_page=10",
		owner,
		repo,
	)

	if err := client.GetJSON(url, &resp); err != nil {
		return err
	}

	// ✅ RELATIVE path ONLY (CRITICAL FIX)
	return github.WriteJSON(
		"github/workflows/workflow_runs.json",
		resp.WorkflowRuns,
	)
}
