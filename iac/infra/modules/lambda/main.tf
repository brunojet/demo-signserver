resource "aws_lambda_function" "this" {
  function_name = "${var.project_env}-${var.name}"
  tags = var.tags
  role          = aws_iam_role.lambda_exec.arn
  handler       = var.handler
  runtime       = var.runtime
  filename      = var.filename != null && var.filename != "" ? var.filename : data.archive_file.lambda_zip[0].output_path
  source_code_hash = var.filename != null && var.filename != "" ? filebase64sha256(var.filename) : data.archive_file.lambda_zip[0].output_base64sha256
  environment {
    variables = var.environment_variables
  }
  timeout = var.timeout
  memory_size = var.memory_size
}

resource "aws_iam_role" "lambda_exec" {
  name = "${var.name}-exec-role"
  tags = var.tags
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_cloudwatch_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "local_file" "lambda_hello" {
  count    = var.filename == null || var.filename == "" ? 1 : 0
  content  = <<EOF
exports.handler = async (event) => {
  // Extrai bucket e chave do evento S3 (via EventBridge)
  const record = event && event.detail && event.detail.requestParameters ? event.detail.requestParameters : (event && event.detail ? event.detail : event);
  const bucket = record && record.bucketName ? record.bucketName : (record && record.bucket && record.bucket.name ? record.bucket.name : "(bucket não encontrado)");
  const key = record && record.key ? record.key : (record && record.object && record.object.key ? record.object.key : "(chave não encontrada)");
  console.log('Bucket: ' + bucket + ', Key: ' + key);
  return {
    statusCode: 200,
    body: JSON.stringify({ message: "Hello, world!", bucket, key })
  };
};
EOF
  filename = "${path.module}/hello.js"
}

data "archive_file" "lambda_zip" {
  count       = var.filename == null || var.filename == "" ? 1 : 0
  type        = "zip"
  source_file = local_file.lambda_hello[0].filename
  output_path = "${path.module}/hello.zip"
}
