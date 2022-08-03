# GitHub Enterprise Org Controller

A Kubernetes controller for managing Weaveworks GitHub Enterprise Organisations.

Maintainers:
- `piaras@weave.works` (@phoban01)

Features:
- Create Enterprise Organisations under the Weaveworks Enterprise
- Create repositories in the new organisation
- Create teams in the organization with access to repositories

## Installation

A Helm chart has been provided to deploy the controller. You must configure the `billingEmail` and `enterpriseName` via the Chart when deploying the controller. See `e2e/terraform/main.tf` for a sample Terraform configuration that uses the Helm chart to deploy the controller and a Terraform module to create an `GitHubOrg` Kubernetes resource.

The helm chart has been released to oci://ghcr.io/weaveworks/charts/gh-org-controller

## Usage

**Please note that this controller will not delete or remove any resources from GitHub.**

A PAT token with sufficient privileges to create organisations, repositories and teams must exist in the same namespace as the `GitHubOrg` resource. It is recommended to rotate this token frequently.

### Organisations

To create an organisations provide the `spec.orgId`, `spec.orgName` and GitHub usernames of organisation administrators (`spec.adminLogins`).

The controller will not update organisation details once the organisation has been created. So if you wish to modify `adminLogins` or delete the organisation you must do this via the Web UI.

```yaml
apiVersion: corp.weave.works/v1alpha1
kind: GitHubOrg
metadata:
  name: my-github-org
  namespace: organisations
spec:
  orgId: my-org
  orgName: My Organisation
  adminLogins:
  - orgAdmin01
  - orgAdmin02
  tokenSecretName: gh-secret
  debug: true
  repositories:
  - name: my-repo
    visibility: private
    autoInit: true
  teams:
  - name: red-team
    description: |-
      A truly awesome team of GitHubers.
    maintainers:
    - alice
    - bob
  - name: blue-team
    maintainers:
    - jerry
    - mary
```

### Repositories

Repositories can be added to the organisation using the `spec.repositories`. Repositories can be added to the organisation at any time however they must be removed manually using the Web UI. The `repositories` field takes the following parameters:

```yaml
name: string # the name of the repository
description: string # a description of the repository
visibility: string # whether the repository is public, private or internal (default: private)
autoInit: bool #  create aninitial commit in the repository (default: true)
hasIssues: bool # enable issues for the repository (default: true)
hasProjects: bool # enable projects for repository (default: false)
```

### Teams

Teams can be created using the `spec.teams` field. Teams are added to all organisation repositories with `admin` privileges. Teams can be added at any time however the must be removed manually using the Web UI. The `teams` field takes the following parameters:

```yaml
name: string # the name of the team
description: string # a description of the team
maintainers: []string # a list of GitHub usernames who will be added to the team as maintainers
```
