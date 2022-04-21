package controllers

import (
	"context"

	"github.com/shurcooL/githubv4"
	ctrl "sigs.k8s.io/controller-runtime"
)

// getEnterpiseID takes an enterprise name and returns the GitHub enterprise ID
func getEnterpiseID(ctx context.Context, client *githubv4.Client, enterpriseName string) (string, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("querying for the enterprise id")

	var query struct {
		Enterprise struct {
			ID githubv4.String
		} `graphql:"enterprise(slug: $slug)"`
	}

	variables := map[string]interface{}{
		"slug": githubv4.String(enterpriseName),
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return "", err
	}

	log.Info("retrieved id", "enterprise", enterpriseName, "id", query.Enterprise.ID)

	return string(query.Enterprise.ID), nil
}
