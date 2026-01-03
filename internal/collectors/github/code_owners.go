package github

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrCodeOwnersNotConfigured = errors.New("codeowners not configured")

type CodeOwnerEntry struct {
	Repository  string    `json:"repository"`
	Path        string    `json:"path"`
	Owners      []string  `json:"owners"`
	CollectedAt time.Time `json:"collected_at"`
}

type CodeOwners struct {
	GeneratedAt time.Time        `json:"generated_at"`
	Total       int              `json:"total"`
	Entries     []CodeOwnerEntry `json:"entries"`
}

func WriteCodeOwners() error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	owner, err := GetOwnerFromEvidence()
	if err != nil {
		return err
	}

	repos, err := LoadRepositoriesFromEvidence()
	if err != nil {
		return err
	}

	entries := []CodeOwnerEntry{}

	for _, repo := range repos {
		content, err := fetchCodeOwnersFile(client, owner, repo.Name)
		if err != nil {
			if errors.Is(err, ErrCodeOwnersNotConfigured) {
				continue // valid audit state
			}
			return err
		}

		parsed := parseCodeOwners(repo.FullName, content)
		entries = append(entries, parsed...)
	}

	model := CodeOwners{
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
		Total:       len(entries),
	}

	return WriteJSON(
		"github/code_owners.json",
		model,
	)
}

func fetchCodeOwnersFile(client *Client, owner, repo string) (string, error) {
	paths := []string{
		".github/CODEOWNERS",
		"CODEOWNERS",
		"docs/CODEOWNERS",
	}

	for _, path := range paths {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/contents/%s",
			owner,
			repo,
			path,
		)

		var resp struct {
			Content  string `json:"content"`
			Encoding string `json:"encoding"`
		}

		err := client.GetJSON(url, &resp)
		if err == nil && resp.Encoding == "base64" {
			raw, _ := base64.StdEncoding.DecodeString(resp.Content)
			return string(raw), nil
		}

		if IsNotFound(err) {
			continue
		}
	}

	return "", ErrCodeOwnersNotConfigured
}

func parseCodeOwners(repo, content string) []CodeOwnerEntry {
	lines := strings.Split(content, "\n")
	entries := []CodeOwnerEntry{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		entries = append(entries, CodeOwnerEntry{
			Repository:  repo,
			Path:        fields[0],
			Owners:      fields[1:],
			CollectedAt: time.Now().UTC(),
		})
	}

	return entries
}
