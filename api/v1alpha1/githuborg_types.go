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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GitHubOrgSpec defines the desired state of GitHubOrg
type GitHubOrgSpec struct {
	// The logins for the administrators of the new organization.
	// +required
	AdminLogins []string `json:"adminLogins"`

	// The id of the new organization.
	// +required
	OrgID string `json:"orgId"`

	// The name of the new organization.
	// +required
	OrgName string `json:"orgName"`

	// Then name of the secret containing the GitHub token.
	// +required
	TokenSecretName string `json:"tokenSecretName"`

	// +optional
	Debug bool `json:"debug,omitempty"`

	// +optional
	Repositories []GitHubRepository `json:"repositories,omitempty"`

	// +optional
	Teams []GitHubTeam `json:"teams,omitempty"`
}

// GitHubRepository defines the GitHub repository to be created.
type GitHubRepository struct {
	// The name of the new repository.
	// +required
	Name string `json:"name"`

	// The description of the new repository.
	// +optional
	Description string `json:"description,omitempty"`

	// The visibility of the new repository.
	// +kubebuilder:default:=private
	// +kubebuilder:enum:=public|private|internal
	// +optional
	Visibility string `json:"visibility,omitempty"`

	// Either true to enable issues for this repository or false to disable them.
	// +kubebuilder:default:=true
	// +optional
	HasIssues bool `json:"hasIssues,omitempty"`

	// Either true to enable projects for this repository or false to disable them.
	// +kubebuilder:default:=false
	// +optional
	HasProjects bool `json:"hasProjects,omitempty"`

	// Create an initial commit with empty README.
	// +kubebuilder:default:=true
	// +optional
	AutoInit bool `json:"autoInit,omitempty"`
}

// GitHubTeam is a group of GitHub users.
type GitHubTeam struct {
	// The name of the new team.
	// +required
	Name string `json:"name"`

	// The description of the new team.
	// +optional
	Description string `json:"description,omitempty"`

	// The logins of the members of the new team.
	// +required
	Maintainers []string `json:"maintainers"`

	// TODO: implement support for adding members to teams.
	// The logins of the members of the new team.
	// +optional
	// Members []string `json:"members,omitempty"`
}

// GitHubOrgStatus defines the observed state of GitHubOrg
type GitHubOrgStatus struct {
	// ObservedGeneration is the last reconciled generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GitHubOrg is the Schema for the githuborgs API
type GitHubOrg struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitHubOrgSpec   `json:"spec,omitempty"`
	Status GitHubOrgStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GitHubOrgList contains a list of GitHubOrg
type GitHubOrgList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitHubOrg `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitHubOrg{}, &GitHubOrgList{})
}
