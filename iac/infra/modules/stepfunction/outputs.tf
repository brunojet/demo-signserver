output "resource_arn" {
  value = aws_sfn_state_machine.this.arn
}

output "role_name" {
  value = aws_iam_role.stepfunction_exec.name
}
