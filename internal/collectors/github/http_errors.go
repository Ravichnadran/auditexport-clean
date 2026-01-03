package github

import (
	"errors"
	"net/http"
)

// GitHubAPIError wraps HTTP status failures
type GitHubAPIError struct {
	StatusCode int
	Message    string
}

func (e *GitHubAPIError) Error() string {
	return e.Message
}

// IsNotFound checks for 404 errors
func IsNotFound(err error) bool {
	var apiErr *GitHubAPIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsForbidden checks for 403 errors
func IsForbidden(err error) bool {
	var apiErr *GitHubAPIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}
