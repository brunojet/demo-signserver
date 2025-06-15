resource "aws_cloudwatch_event_rule" "this" {
  name          = "${var.project_env}-${var.name}-eventbridge"
  tags = var.tags
  event_pattern = var.event_pattern
  event_bus_name = var.event_bus_name
}

resource "aws_iam_role" "eventbridge_invoke" {
  name = "${aws_cloudwatch_event_rule.this.name}-invoke-role"
  tags = var.tags
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = { Service = "events.amazonaws.com" },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_cloudwatch_event_target" "this" {
  rule         = aws_cloudwatch_event_rule.this.name
  arn          = var.target_arn
  role_arn     = aws_iam_role.eventbridge_invoke.arn
  event_bus_name = var.event_bus_name
}
