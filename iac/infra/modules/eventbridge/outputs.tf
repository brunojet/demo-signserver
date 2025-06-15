output "resource_name" {
  description = "Nome da regra EventBridge criada."
  value       = aws_cloudwatch_event_rule.this.name
}

output "resource_arn" {
  description = "ARN da regra EventBridge criada."
  value       = aws_cloudwatch_event_rule.this.arn
}

output "role_name" {
  description = "Nome da role que pode ser usada para anexar policies para invocar a regra EventBridge."
  value       = aws_iam_role.eventbridge_invoke.name
}
