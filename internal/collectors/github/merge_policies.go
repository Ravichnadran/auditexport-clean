package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type MergePolicyRule struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

type MergePolicies struct {
	Repository  string            `json:"repository"`
	Branch      string            `json:"branch"`
	Enabled     bool              `json:"enabled"`
	Status      string            `json:"status"`
	Rules       []MergePolicyRule `json:"rules"`
	CollectedAt time.Time         `json:"collected_at"`
}

func WriteMergePolicies(owner, repo, branch string) error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/repos/%s/%s/branches/%s/protection",
		baseURL,
		owner,
		repo,
		branch,
	)

	req, err := client.newRequest("GET", url)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	policies := MergePolicies{
		Repository:  repo,
		Branch:      branch,
		CollectedAt: time.Now().UTC(),
	}

	// ✅ IMPORTANT: non-fatal handling
	switch resp.StatusCode {

	case http.StatusOK:
		// Branch protection exists
		policies.Enabled = true
		policies.Status = "enabled"
		policies.Rules = []MergePolicyRule{
			{
				Type:   "branch_protection",
				Detail: "Branch protection rules enabled",
			},
		}

	case http.StatusNotFound:
		// Branch protection NOT enabled (VALID AUDIT STATE)
		policies.Enabled = false
		policies.Status = "not_enabled"
		policies.Rules = []MergePolicyRule{}

	case http.StatusForbidden:
		// Insufficient permissions (VALID AUDIT STATE)
		policies.Enabled = false
		policies.Status = "insufficient_permissions"
		policies.Rules = []MergePolicyRule{}

	default:
		return fmt.Errorf("unexpected GitHub response: %d", resp.StatusCode)
	}

	data, err := json.MarshalIndent(policies, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "merge_policies.json"),
		data,
		0644,
	)
}
