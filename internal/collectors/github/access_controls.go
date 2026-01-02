package github

import (
	"auditexport/internal/run"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type RepoAccess struct {
	Repository  string    `json:"repository"`
	Username    string    `json:"username"`
	Permission  string    `json:"permission"`
	CollectedAt time.Time `json:"collected_at"`
}

type AccessControlsOutput struct {
	CollectedAt time.Time    `json:"collected_at"`
	Total       int          `json:"total"`
	Entries     []RepoAccess `json:"entries"`
}

func WriteAccessControls() error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	repoBytes, err := os.ReadFile(
		run.EvidencePath("github", "repositories.json"),
	)
	if err != nil {
		return err
	}

	var repoData RepositoriesOutput
	if err := json.Unmarshal(repoBytes, &repoData); err != nil {
		return err
	}

	var results []RepoAccess

	for _, repo := range repoData.Repositories {
		page := 1

		for {
			url := fmt.Sprintf(
				"%s/repos/%s/collaborators?per_page=100&page=%d",
				baseURL,
				repo.FullName,
				page,
			)

			req, err := client.newRequest("GET", url)
			if err != nil {
				return err
			}

			resp, err := client.httpClient.Do(req)
			if err != nil {
				return err
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				return fmt.Errorf(
					"collaborators API error for %s: %s",
					repo.FullName,
					resp.Status,
				)
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}

			var collaborators []struct {
				Login string `json:"login"`
			}

			if err := json.Unmarshal(body, &collaborators); err != nil {
				return err
			}

			if len(collaborators) == 0 {
				break
			}

			for _, c := range collaborators {
				permURL := fmt.Sprintf(
					"%s/repos/%s/collaborators/%s/permission",
					baseURL,
					repo.FullName,
					c.Login,
				)

				permReq, err := client.newRequest("GET", permURL)
				if err != nil {
					return err
				}

				permResp, err := client.httpClient.Do(permReq)
				if err != nil {
					return err
				}

				if permResp.StatusCode != 200 {
					permResp.Body.Close()
					return fmt.Errorf(
						"permission API error for %s/%s",
						repo.FullName,
						c.Login,
					)
				}

				permBody, err := io.ReadAll(permResp.Body)
				permResp.Body.Close()
				if err != nil {
					return err
				}

				var permData struct {
					Permission string `json:"permission"`
				}

				if err := json.Unmarshal(permBody, &permData); err != nil {
					return err
				}

				results = append(results, RepoAccess{
					Repository:  repo.FullName,
					Username:    c.Login,
					Permission:  permData.Permission,
					CollectedAt: time.Now().UTC(),
				})
			}

			page++
		}
	}

	output := AccessControlsOutput{
		CollectedAt: time.Now().UTC(),
		Total:       len(results),
		Entries:     results,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		run.EvidencePath("github", "access_controls.json"),
		data,
		0644,
	)
}
