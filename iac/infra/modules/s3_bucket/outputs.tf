output "resource_name" { value = aws_s3_bucket.this.bucket }
output "resource_arn" { value = aws_s3_bucket.this.arn }
output "role_name" { value = aws_iam_role.s3_access.name }
output "policy_arn" { value = aws_iam_policy.s3_policy.arn }
