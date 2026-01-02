package run

import "path/filepath"

// EvidencePath builds a path inside the run's evidence directory.
func EvidencePath(parts ...string) string {
	InitRunContext()

	all := append([]string{BaseDir(), "evidence"}, parts...)
	return filepath.Join(all...)
}
