package v1alpha1

import (
	"strings"
)

// GitHubError is an error type that represents an error returned from the GitHub API.
type GitHubError string

func (e GitHubError) Error() string {
	return string(e)
}

const (
	// GitHubAlreadyExistsError is returned when a resource already exists.
	GitHubAlreadyExistsError GitHubError = "GitHubResourceAlreadyExists"

	// GitHubResourceNotFound is returned when a resource already is missing.
	GitHubResourceNotFound GitHubError = "GitHubResourceNotFound"
)

// CheckGitHubError maps the GitHub error string to a GitHubError.
func CheckGitHubError(err error) error {
	s := err.Error()
	switch {
	case strings.Contains(s, "exists"):
		return GitHubAlreadyExistsError
	case strings.Contains(s, "not found"):
		return GitHubResourceNotFound
	case strings.Contains(s, "not resolve"):
		return GitHubResourceNotFound
	default:
		return err
	}
}

// IgnoreAlreadyExists returns true if the error is GitHubAlreadyExistsError.
func IgnoreAlreadyExists(err error) error {
	if err == GitHubAlreadyExistsError {
		return nil
	}
	return err
}
