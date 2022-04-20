package controllers

import (
	"context"

	"github.com/shurcooL/githubv4"
	corpv1alpha1 "github.com/weaveworks/corp-gh-org-controller/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var createOrgMutation struct {
	CreateEnterpriseOrganization struct {
		Organization struct {
			Login githubv4.String
		}
	} `graphql:"createOrganization(input: $input)"`
}

// Not implemented as we need to understand the query which will return a list of existing
// enterprise admins before we allow mutation of admins.
//
// For now add or remove admins using the Web UI.
//

// var inviteEnterpriseAdminMutation struct {
//     InviteEnterpriseAdmin struct {
//         Organization struct {
//             Login githubv4.String
//         }
//     } `graphql:"inviteEnterpriseAdmin(input: $input)"`
// }
//
// var removeEnterpriseAdminMutation struct {
//     RemoveEnterpriseAdmin struct {
//         Organization struct {
//             Login githubv4.String
//         }
//     } `graphql:"removeEnterpriseAdmin(input: $input)"`
// }

func createEnterpriseOrganization(
	ctx context.Context,
	client *githubv4.Client,
	enterpriseID, billingEmail, orgID, orgName string,
	adminLogins []string,
) error {
	log := ctrl.LoggerFrom(ctx)
	var admins []githubv4.String
	for _, adminLogin := range adminLogins {
		admins = append(admins, githubv4.String(adminLogin))
	}

	// return if the org already exists or we encounter an error
	if orgExists, err := checkIfOrgExists(ctx, client, orgID); orgExists || err != nil {
		return err
	}

	input := githubv4.CreateEnterpriseOrganizationInput{
		EnterpriseID: githubv4.String(enterpriseID),
		Login:        githubv4.String(orgID),
		ProfileName:  githubv4.String(orgName),
		BillingEmail: githubv4.String(billingEmail),
		AdminLogins:  admins,
	}

	log.Info("creating organisation", "name", orgName, "id", orgID)

	if err := client.Mutate(ctx, createOrgMutation, input, nil); err != nil {
		log.Info("error code", "error", err.Error())
		return err
	}

	log.Info("organisation created", "name", orgName, "id", orgID)

	return nil
}

func checkIfOrgExists(ctx context.Context, client *githubv4.Client, orgID string) (bool, error) {
	var getOrg struct {
		Organization struct {
			ID githubv4.String
		} `graphql:"Organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(orgID),
	}

	err := client.Query(ctx, &getOrg, variables)
	if err != nil {
		if corpv1alpha1.CheckGitHubError(err) == corpv1alpha1.GitHubResourceNotFound {
			return false, nil
		}
		return true, err
	}

	return true, nil
}
