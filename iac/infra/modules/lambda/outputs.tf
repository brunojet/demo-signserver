output "resource_arn" { value = aws_lambda_function.this.arn }

output "resource_name" { value = aws_lambda_function.this.function_name }

output "role_name" {
  value = aws_iam_role.lambda_exec.name
}
