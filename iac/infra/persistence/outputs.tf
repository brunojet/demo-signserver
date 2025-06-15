output "persistence_bucket_name" {
  description = "Name of the S3 bucket created for persistence."
  value       = module.persistence_bucket.resource_name
}

output "persistence_bucket_arn" {
  description = "ARN do bucket S3 de persistência."
  value       = module.persistence_bucket.resource_arn
}

output "persistence_bucket_role_name" {
  description = "ARN da IAM Role para acesso ao bucket de persistência."
  value       = module.persistence_bucket.role_name
}

output "persistence_bucket_policy_arn" {
  description = "ARN da IAM Policy para acesso ao bucket de persistência."
  value       = module.persistence_bucket.policy_arn
}

output "device_profile_table_name" {
  description = "Nome da tabela DynamoDB device-profile."
  value       = module.device_profile_table.resource_name
}

output "device_profile_table_arn" {
  description = "ARN da tabela DynamoDB device-profile."
  value       = module.device_profile_table.resource_arn
}

output "device_profile_role_name" {
  description = "ARN da IAM Role para acesso à tabela device-profile."
  value       = module.device_profile_table.role_name
}

output "device_profile_policy_arn" {
  description = "ARN da IAM Policy para acesso à tabela device-profile."
  value       = module.device_profile_table.policy_arn
}

output "signature_request_table_name" {
  description = "Nome da tabela DynamoDB signature-request."
  value       = module.signature_request_table.resource_name
}

output "signature_request_table_arn" {
  description = "ARN da tabela DynamoDB signature-request."
  value       = module.signature_request_table.resource_arn
}

output "signature_request_role_name" {
  description = "ARN da IAM Role para acesso à tabela signature-request."
  value       = module.signature_request_table.role_name
}

output "signature_request_policy_arn" {
  description = "ARN da IAM Policy para acesso à tabela signature-request."
  value       = module.signature_request_table.policy_arn
}
