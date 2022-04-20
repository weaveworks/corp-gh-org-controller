package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shurcooL/githubv4"
)

func TestGetEnterpriseId(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "should return an id",
			input:    "octocat",
			expected: "1000:octocat",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data": {"enterprise": {"id": "1000:octocat"}}}`))
			}))

			client := githubv4.NewEnterpriseClient(server.URL, http.DefaultClient)

			got, err := getEnterpiseID(ctx, client, test.input)

			if err != nil {
				t.Errorf("getEnterpiseID returned an error: %v", err)
			}

			if got != test.expected {
				t.Errorf("getEnterpiseID(%q) = %q, want %q", test.input, got, test.expected)
			}
		})
	}
}
