terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.10"
    }
  }

  experiments = [module_variable_optional_attrs]
}

resource "kubernetes_manifest" "this" {
  manifest = {
    "apiVersion" = "corp.weave.works/v1alpha1"
    "kind"       = "GitHubOrg"
    "metadata" = {
      "name"      = var.org_id
      "namespace" = var.namespace
    }
    "spec" = {
      "orgId"           = var.org_id
      "orgName"         = var.org_name
      "adminLogins"     = var.org_admins
      "tokenSecretName" = var.token_secret_name
      "debug"           = var.debug
      "repositories" : var.repositories
      "teams" : var.teams
    }
  }
}
