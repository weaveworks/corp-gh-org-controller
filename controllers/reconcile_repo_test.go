package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v43/github"
	corpv1alpha1 "github.com/weaveworks/corp-gh-org-controller/api/v1alpha1"
)

func TestReconcileRepository(t *testing.T) {
	tests := []struct {
		name     string
		orgID    string
		repoName string
		wantErr  bool
	}{
		{
			name:     "create repo call should succeed",
			orgID:    "1",
			repoName: "test-repo",
			wantErr:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			}))

			client, err := github.NewEnterpriseClient(server.URL, "", http.DefaultClient)
			if err != nil {
				t.Errorf("create client returned an error: %v", err)
			}

			err = reconcileRepository(ctx,
				client,
				test.orgID,
				corpv1alpha1.GitHubRepository{Name: test.repoName},
			)

			if err != nil {
				t.Errorf("createRepo returned an error: %v", err)
			}
		})
	}
}
