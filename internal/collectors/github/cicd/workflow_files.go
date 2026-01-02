package cicd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"auditexport/internal/collectors/github"
	"auditexport/internal/run"
)

func WriteWorkflowFiles(owner, repo string) error {
	client, err := github.NewClient()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/.github/workflows",
		owner, repo,
	)

	var files []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	if err := client.GetJSON(url, &files); err != nil {
		// workflows folder may not exist — valid audit state
		return nil
	}

	dir := run.EvidencePath("github", "workflows")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for _, f := range files {
		var content struct {
			Content  string `json:"content"`
			Encoding string `json:"encoding"`
		}

		if err := client.GetJSON(f.URL, &content); err != nil {
			continue
		}

		if content.Encoding != "base64" {
			continue
		}

		data, _ := base64.StdEncoding.DecodeString(content.Content)
		path := filepath.Join(dir, f.Name)

		_ = os.WriteFile(path, data, 0644)
	}

	return nil
}
