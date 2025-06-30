provider "aws" {
  region = var.region
}

module "persistence_bucket" {
  source      = "../modules/s3_bucket"
  bucket_name = "storage"
  versioning  = false
  tags        = local.tags
  project_env = local.project_env
}

module "device_profile_table" {
  source     = "../modules/dynamodb_table"
  table_name = "profile"
  hash_key   = "pk"
  range_key  = "sk"
  attributes = [
    { name = "pk", type = "S" },
    { name = "sk", type = "S" }
  ]
  tags        = local.tags
  project_env = local.project_env
}

module "signature_request_table" {
  source     = "../modules/dynamodb_table"
  table_name = "request"
  hash_key   = "pk"
  attributes = [
    { name = "pk", type = "S" }
  ]
  tags        = local.tags
  project_env = local.project_env
}

resource "aws_resourcegroups_group" "env_group" {
  name        = "${local.project_env}-resources"
  tags        = local.tags
  description = "Resource Group para o ambiente ${local.project_env}"
  resource_query {
    query = jsonencode({
      ResourceTypeFilters = ["AWS::AllSupported"],
      TagFilters = [
        { Key = "Name", Values = [local.tags.Name] },
        { Key = "Environment", Values = [local.tags.Environment] }
      ]
    })
  }
}
