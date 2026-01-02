// internal/collectors/github/context.go
package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"os"
)

type OrgContext struct {
	Login string `json:"login"`
}

func GetOwnerFromEvidence() (string, error) {
	data, err := os.ReadFile(
		run.EvidencePath("github", "organization.json"),
	)
	if err != nil {
		return "", err
	}

	var org OrgContext
	if err := json.Unmarshal(data, &org); err != nil {
		return "", err
	}

	return org.Login, nil
}
