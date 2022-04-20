/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/google/go-github/v43/github"
	"github.com/shurcooL/githubv4"
	corpv1alpha1 "github.com/weaveworks/corp-gh-org-controller/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// GitHubOrgReconciler reconciles a GitHubOrg object
type GitHubOrgReconciler struct {
	client.Client
	Scheme            *runtime.Scheme
	billingEmail      string
	enterpriseName    string
	githubAPIEndpoint string
}

type GitHubOrgReconcilerOptions struct {
	BillingEmail      string
	EnterpriseName    string
	GitHubAPIEndpoint string
}

//+kubebuilder:rbac:groups=corp.weave.works,resources=githuborgs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=corp.weave.works,resources=githuborgs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=corp.weave.works,resources=githuborgs/finalizers,verbs=update
//+kubebuilder:rbac:groups=corp.weave.works,resources=githuborgs/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=configmaps;secrets;serviceaccounts,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// SetupWithManager sets up the controller with the Manager.
func (r *GitHubOrgReconciler) SetupWithManager(mgr ctrl.Manager, opts *GitHubOrgReconcilerOptions) error {
	r.billingEmail = opts.BillingEmail
	r.enterpriseName = opts.EnterpriseName

	rateLimiter := workqueue.NewItemExponentialFailureRateLimiter(
		5*time.Second,
		5*time.Minute,
	)

	return ctrl.NewControllerManagedBy(mgr).
		For(&corpv1alpha1.GitHubOrg{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 1,
			RateLimiter:             rateLimiter,
		}).
		Complete(r)
}

// Reconcile reads that state of the cluster for a GitHubOrg object and makes changes based on the state read
func (r *GitHubOrgReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	reconcileStart := time.Now()

	var org corpv1alpha1.GitHubOrg
	if err := r.Get(ctx, req.NamespacedName, &org); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	tokenKey := types.NamespacedName{
		Namespace: org.Namespace,
		Name:      org.Spec.TokenSecretName,
	}

	oauthClient, err := r.getOauth2Client(ctx, tokenKey)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Token secret not found, skipping reconciliation")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	graphqlClient := githubv4.NewClient(oauthClient)
	restClient := github.NewClient(oauthClient)

	err = r.reconcile(ctx, graphqlClient, restClient, &org)
	if err != nil {
		return ctrl.Result{}, err
	}

	msg := fmt.Sprintf("Reconciliation finished in %s", time.Since(reconcileStart).String())

	log.Info(msg)

	return ctrl.Result{}, nil
}

// reconcile will reconcile the GitHubOrg.
func (r *GitHubOrgReconciler) reconcile(ctx context.Context, graphqlClient *githubv4.Client, restClient *github.Client, obj *corpv1alpha1.GitHubOrg) error {
	log := ctrl.LoggerFrom(ctx)

	var org string
	org = obj.Spec.OrgID

	// skip creating org if we are in debug mode
	if !obj.Spec.Debug {
		enterpriseID, err := getEnterpiseID(ctx, graphqlClient, r.enterpriseName)
		if err != nil {
			return err
		}

		if err := createEnterpriseOrganization(ctx,
			graphqlClient,
			enterpriseID,
			r.billingEmail,
			org,
			obj.Spec.OrgName,
			obj.Spec.AdminLogins,
		); corpv1alpha1.IgnoreAlreadyExists(err) != nil {
			return err
		}
	}

	for _, repo := range obj.Spec.Repositories {
		// TODO: check if the repo exists
		// if !repoExists {
		err := reconcileRepository(ctx, restClient, org, repo)
		if err != nil {
			log.Info("error creating repo", "error", err)
		}
		// if corpv1alpha1.IgnoreAlreadyExists(err) != nil {
		//     return err
		// }
		// }
		// TODO: update the repo settings
	}

	// ensure teams exist and have access to the repositories
	for _, team := range obj.Spec.Teams {
		err := reconcileTeam(ctx, restClient, org, team, obj.Spec.Repositories)
		if err != nil {
			log.Info("error creating team", "error", err)
		}
	}

	return nil
}

func (r *GitHubOrgReconciler) getOauth2Client(ctx context.Context, tokenKey types.NamespacedName) (*http.Client, error) {
	creds := &corev1.Secret{}
	if err := r.Get(ctx, tokenKey, creds); err != nil {
		return nil, err
	}

	if _, ok := creds.Data["GITHUB_TOKEN"]; ok != true {
		return nil, fmt.Errorf("GITHUB_TOKEN not found in secret %s", tokenKey)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(creds.Data["GITHUB_TOKEN"])},
	)

	return oauth2.NewClient(ctx, src), nil
}
