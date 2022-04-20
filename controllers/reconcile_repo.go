package controllers

import (
	"context"
	"net/http"

	"github.com/google/go-github/v43/github"
	corpv1alpha1 "github.com/weaveworks/corp-gh-org-controller/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func reconcileRepository(ctx context.Context, client *github.Client, org string, repo corpv1alpha1.GitHubRepository) error {
	log := ctrl.LoggerFrom(ctx)

	repoConfig := &github.Repository{
		Name:        &repo.Name,
		Description: &repo.Description,
		Visibility:  &repo.Visibility,
		HasIssues:   &repo.HasIssues,
		HasProjects: &repo.HasProjects,
		AutoInit:    &repo.AutoInit,
	}

	existingRepo, resp, err := client.Repositories.Get(ctx, org, repo.Name)
	if err != nil {
		log.V(0).Info("error querying repo", "name", repo.Name, "org", org, "code", resp.StatusCode)
		switch {
		case resp.StatusCode == http.StatusNotFound:
			log.V(0).Info("creating repository", "name", repo.Name, "org", org)
			_, _, err := client.Repositories.Create(ctx, org, repoConfig)
			if err != nil {
				return err
			}
			log.V(0).Info("repository created", "name", repo.Name, "org", org)
			return nil
		default:
			return err
		}
	}

	log.Info("updating repository", "name", repo.Name, "org", org)

	repoModified := false
	var nullString *string

	// currently we only allow modifying the visibility and description
	if *existingRepo.Visibility != repo.Visibility {
		repoConfig.Visibility = &repo.Visibility
		repoModified = true
	} else {
		repoConfig.Visibility = nullString
	}

	if existingRepo.Description != nil &&
		*existingRepo.Description != repo.Description {
		repoConfig.Description = &repo.Description
		repoModified = true
	} else {
		repoConfig.Description = nullString
	}

	if repoModified {
		_, _, err = client.Repositories.Edit(ctx, org, *repoConfig.Name, repoConfig)
		if err != nil {
			return err
		}
	}

	log.Info("repository updated", "name", repo.Name, "org", org)

	return nil
}
