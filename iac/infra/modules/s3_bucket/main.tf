resource "aws_s3_bucket" "this" {
  bucket = "${var.project_env}-${var.bucket_name}"
  tags   = var.tags

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket_versioning" "this" {
  count  = var.versioning ? 1 : 0
  bucket = aws_s3_bucket.this.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_iam_role" "s3_access" {
  name = "${aws_s3_bucket.this.bucket}-access-role"
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

resource "aws_iam_policy" "s3_policy" {
  name        = "${aws_s3_bucket.this.bucket}-access-policy"
  tags = var.tags
  description = "Policy for S3 bucket access"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ]
      Resource = ["${aws_s3_bucket.this.arn}", "${aws_s3_bucket.this.arn}/*"]
    }]
  })
}

resource "aws_iam_role_policy_attachment" "attach" {
  role       = aws_iam_role.s3_access.name
  policy_arn = aws_iam_policy.s3_policy.arn
}
