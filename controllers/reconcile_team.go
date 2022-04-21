package controllers

import (
	"context"
	"net/http"

	"github.com/google/go-github/v43/github"
	corpv1alpha1 "github.com/weaveworks/corp-gh-org-controller/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// reconcileTeam creates or updates team membership and adds the team to the
// given repositories.
func reconcileTeam(ctx context.Context, client *github.Client,
	org string, team corpv1alpha1.GitHubTeam, repos []corpv1alpha1.GitHubRepository) error {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling team", "team", team.Name, "org", org)

	teamConfig := github.NewTeam{
		Name:        team.Name,
		Description: &team.Description,
		Maintainers: team.Maintainers,
	}

	existingTeam, resp, err := client.Teams.GetTeamBySlug(ctx, org, team.Name)
	if err != nil {
		log.V(0).Info("error querying team", "name", team.Name, "org", org, "code", resp.StatusCode)
		switch {
		case resp.StatusCode == http.StatusNotFound:
			log.V(0).Info("creating team", "name", team.Name, "org", org)
			_, _, err := client.Teams.CreateTeam(ctx, org, teamConfig)
			if err != nil {
				return err
			}
			log.V(0).Info("team created", "name", team.Name, "org", org)
		default:
			return err
		}
	}

	if existingTeam != nil {
		log.V(0).Info("updating team", "name", team.Name, "org", org)
		_, _, err := client.Teams.EditTeamBySlug(ctx, org, team.Name, teamConfig, false)
		if err != nil {
			return err
		}
		log.V(0).Info("team updated", "name", team.Name, "org", org)
	}

	opts := &github.TeamAddTeamRepoOptions{
		Permission: "admin",
	}

	teamRepos, _, err := client.Teams.ListTeamReposBySlug(ctx, org, team.Name, nil)
	if err != nil {
		return err
	}

	for _, r := range repos {
		if !containsRepo(r.Name, teamRepos) {
			_, err = client.Teams.AddTeamRepoBySlug(ctx, org, team.Name, org, r.Name, opts)
			if err != nil {
				return err
			}
			log.Info("Added team to repo", "team", team.Name, "repo", r.Name, "org", org)
		}
	}

	return nil
}

func containsRepo(r string, rl []*github.Repository) bool {
	for _, repo := range rl {
		if repo.GetName() == r {
			return true
		}
	}
	return false
}
