package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shurcooL/githubv4"
)

func TestCreateEnterpriseOrg(t *testing.T) {
	tests := []struct {
		name         string
		enterpriseID string
		billingEmail string
		orgID        string
		orgName      string
		adminLogins  []string
		wantErr      bool
	}{
		{
			name:         "create enterprise call should succeed",
			enterpriseID: "100:octocat",
			billingEmail: "test@example.com",
			orgID:        "1",
			orgName:      "test",
			adminLogins: []string{
				"test",
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{}`))
			}))

			client := githubv4.NewEnterpriseClient(server.URL, http.DefaultClient)

			err := createEnterpriseOrganization(ctx,
				client,
				test.enterpriseID,
				test.billingEmail,
				test.orgID,
				test.orgName,
				test.adminLogins,
			)

			if err != nil {
				t.Errorf("createEnterpriseOrganization returned an error: %v", err)
			}
		})
	}
}
