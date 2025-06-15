locals {
  allowed_environments = ["dev", "hml", "prod"]
  environment          = terraform.workspace
  project_env          = "${var.project}-${local.environment}"
  is_feature_branch    = can(regex("^feature-", local.environment))
  is_hotfix_branch     = can(regex("^hotfix-", local.environment))
  is_bugfix_branch     = can(regex("^bugfix-", local.environment))
  is_valid_env         = contains(local.allowed_environments, local.environment) || local.is_feature_branch || local.is_hotfix_branch || local.is_bugfix_branch

  tags = {
    Name        = var.project
    Environment = local.environment
    ManagedBy   = "terraform"
    Module      = "orchestrator"
  }
}

resource "null_resource" "validate_workspace" {
  count = local.is_valid_env ? 0 : 1
  provisioner "local-exec" {
    command = "echo Workspace inválido: ${local.environment}. Use dev, qa, prod, feature-*, hotfix-* ou bugfix-*. && exit 1"
  }
}
