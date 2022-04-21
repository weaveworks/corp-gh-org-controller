provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "helm_release" "corp-gh-org-controller" {
  name  = "corp-gh-org-controller"
  chart = "../../charts/corp-gh-org-controller"

  set {
    name  = "billingEmail"
    value = "accounts@acme.org"
  }

  set {
    name  = "enterpriseName"
    value = "acme"
  }
}

module "engagement_007" {
  source            = "../../terraform/corp_gh_org"
  namespace         = "corp-gh-org-controller-system"
  token_secret_name = "gh-secret"
  org_id            = "test_engagement_007"
  org_name          = "Test Engagement Org"
  org_admins = [
    "morancj"
  ]

  debug = true

  teams = [{
    name = "MI5"
    maintainers = [
      "phoban01"
    ]
  }]

  repositories = [
    {
      name        = "journal"
      description = "The engagement journal repo"
    },
    {
      name        = "proof-of-concept"
      description = "Private repo for POC"
    },
    {
      name        = "demos"
      description = "Public demo materials for open-source engagement"
    },
  ]
}
