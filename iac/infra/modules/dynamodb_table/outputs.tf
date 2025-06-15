output "resource_name" { value = aws_dynamodb_table.this.name }
output "resource_arn" { value = aws_dynamodb_table.this.arn }
output "role_name" {
  description = "Nome da role que pode ser usada para anexar policies"
  value       = aws_iam_role.dynamodb_access.name
}

output "policy_arn" {
  description = "ARN da policy default mínima para acesso à tabela DynamoDB"
  value       = aws_iam_policy.dynamodb_policy.arn
}
