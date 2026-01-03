package correlation

import "time"

// ChangeEvent represents ONE production change
// backed by enforced PR + CI evidence.
type ChangeEvent struct {
	Repository string    `json:"repository"`
	PRNumber   int       `json:"pr_number"`
	Author     string    `json:"author"`
	MergedAt   time.Time `json:"merged_at"`
	MergeSHA   string    `json:"merge_sha"`

	CIRuns []CIRunEvidence `json:"ci_runs"`
}

// CIRunEvidence represents a CI run linked to a PR.
type CIRunEvidence struct {
	RunID      int       `json:"run_id"`
	Name       string    `json:"name"`
	Conclusion string    `json:"conclusion"`
	CreatedAt  time.Time `json:"created_at"`
}
