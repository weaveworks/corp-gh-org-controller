variable "namespace" {
  type        = string
  description = "The Kubernetes namespace in which to create the GitHubOrg resource."
}

variable "debug" {
  type        = bool
  description = "When debug mode is enabled the controller will not create organisations but will create repositories and teams."
  default     = false
}

variable "org_id" {
  type        = string
  description = "The identifier of the organization. This will be used to generate the url for the GitHub Org."
}

variable "org_name" {
  type        = string
  description = "The name of the organization. This will be the display name of the organization."
}

variable "org_admins" {
  type        = list(string)
  description = "A list of GitHub usernames that will be administrators of the organisation."
  default     = []
}

variable "token_secret_name" {
  type        = string
  description = "The name of the Kubernetes secret containing the GitHub token with permissions to create enterprise organisations."
}

variable "repositories" {
  type = list(object({
    name        = string           # the repository name
    description = optional(string) # the repository description
    visibility  = optional(string) # the repository visibility, must be one of "public", "private", "internal", defaults to "private"
    autoInit    = optional(bool)   # whether to create an initial commit with README, defaults to true
    hasIssues   = optional(bool)   # whether to enable issues for the repository, defaults to true
    hasProjects = optional(bool)   # whether to enable projects for the repository, defaults to false
  }))
  description = "A list of GitHub repositories that will be created in the organisation."
  default     = []
}

variable "teams" {
  type = list(object({
    name        = string         # the team name
    description = optional(bool) # the team description
    maintainers = list(string)   # the list of GitHub usernames that will be maintainers of the team
  }))
  description = "A list of GitHub teams that will be added to any repositories created in the organisation."
  default     = []
}
