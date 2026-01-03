package cicd

import (
	"fmt"
	"time"

	"auditexport/internal/collectors/github"
)

type Workflow struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type workflowResponse struct {
	Workflows []Workflow `json:"workflows"`
}

func WriteWorkflows(owner, repo string) error {
	client, err := github.NewClient()
	if err != nil {
		return err
	}

	var resp workflowResponse
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/workflows",
		owner,
		repo,
	)

	// if err := client.GetJSON(url, &resp); err != nil {
	// 	return err
	// }
	if err := client.GetJSON(url, &resp); err != nil {
		if github.IsNotFound(err) || github.IsForbidden(err) {
			return ErrCICDNotAvailable
		}
		return err
	}

	// ✅ RELATIVE path ONLY
	return github.WriteJSON(
		"github/workflows/workflows.json",
		resp.Workflows,
	)
}
