resource "aws_dynamodb_table" "this" {
  name         = "${var.project_env}-${var.table_name}"
  tags         = var.tags
  hash_key     = var.hash_key
  range_key    = var.range_key
  billing_mode = var.billing_mode
  dynamic "attribute" {
    for_each = var.attributes
    content {
      name = attribute.value.name
      type = attribute.value.type
    }
  }
  read_capacity  = var.billing_mode == "PROVISIONED" ? var.read_capacity : null
  write_capacity = var.billing_mode == "PROVISIONED" ? var.write_capacity : null

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_iam_role" "dynamodb_access" {
  name = "${aws_dynamodb_table.this.name}-access-role"
  tags = var.tags
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "dynamodb_policy" {
  name        = "${aws_dynamodb_table.this.name}-dynamodb-policy"
  tags = var.tags
  description = "Permissões mínimas para acesso à tabela DynamoDB: GetItem, PutItem, UpdateItem"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Action = [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem"
      ],
      Resource = [
        aws_dynamodb_table.this.arn,
        "${aws_dynamodb_table.this.arn}/*"
      ]
    }]
  })
}

resource "aws_iam_role_policy_attachment" "attach" {
  role       = aws_iam_role.dynamodb_access.name
  policy_arn = aws_iam_policy.dynamodb_policy.arn
}
