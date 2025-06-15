provider "aws" {
  region = var.region
}

data "terraform_remote_state" "persistence" {
  backend = "s3"
  config = {
    bucket  = "brunojet-tfstate"
    key     = "env:/${terraform.workspace}/persistence/terraform.tfstate"
    region  = "${var.region}"
    encrypt = true
  }
}

# ---
# [Ponto crítico] O remote_state é usado para buscar os ARNs dos recursos de persistência (S3/DynamoDB) já criados em outro stack.
# Certifique-se de que o state de persistência está atualizado e acessível antes de rodar este orquestrador.
# ---

# Lambda (placeholder)
module "signer_lambda" {
  source                = "../modules/lambda"
  name                  = "signer-lambda"
  tags                  = local.tags
  project_env           = local.project_env
  handler               = "hello.handler"
  runtime               = "nodejs18.x"
  filename              = null
  environment_variables = {}
  timeout               = 10
  memory_size           = 128
}

# Anexar policies default do persistence às roles do Lambda
resource "aws_iam_role_policy_attachment" "lambda_s3_persistence" {
  role       = module.signer_lambda.role_name
  policy_arn = data.terraform_remote_state.persistence.outputs["persistence_bucket_policy_arn"]
}

resource "aws_iam_role_policy_attachment" "lambda_device_profile_persistence" {
  role       = module.signer_lambda.role_name
  policy_arn = data.terraform_remote_state.persistence.outputs["device_profile_policy_arn"]
}

resource "aws_iam_role_policy_attachment" "lambda_signature_request_persistence" {
  role       = module.signer_lambda.role_name
  policy_arn = data.terraform_remote_state.persistence.outputs["signature_request_policy_arn"]
}

# [Ponto crítico] O nome da função Lambda é usado em policies e permissões. Evite caracteres especiais.
# Se precisar anexar policies extras, use o output role_name deste módulo.

# Step Function
module "signer_flow" {
  source      = "../modules/stepfunction"
  name        = "signer-flow"
  tags        = local.tags
  project_env = local.project_env
  definition  = <<EOF
{
  "Comment": "Signer flow skeleton",
  "StartAt": "InvokeLambda",
  "States": {
    "InvokeLambda": {
      "Type": "Task",
      "Resource": "arn:aws:states:::lambda:invoke",
      "Parameters": {
        "FunctionName": "${module.signer_lambda.resource_arn}",
        "Payload.$": "$"
      },
      "End": true
    }
  }
}
EOF
}

# [Ponto crítico] A role de execução da Step Function é exposta via output stepfunction_exec_role_name.
# Use este output para anexar policies de invoke Lambda ou outras permissões necessárias.

module "s3_object_created_eventbridge" {
  source        = "../modules/eventbridge"
  name          = "s3-object-created"
  tags          = local.tags
  project_env   = local.project_env
  target_arn    = module.signer_flow.resource_arn
  event_pattern = <<EOF
{
  "source": ["aws.s3"],
  "detail-type": ["Object Created"],
  "resources": ["arn:aws:s3:::${data.terraform_remote_state.persistence.outputs.persistence_bucket_name}"]
}
EOF
}

# Caso precise de permissões extras além das policies default, crie policies adicionais aqui.
# Exemplo: Policy para Step Function invocar Lambda
resource "aws_iam_policy" "stepfunction_lambda_invoke" {
  name        = "${local.project_env}-stepfunction-lambda-invoke"
  tags        = local.tags
  description = "Permite à Step Function invocar Lambda"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect   = "Allow",
      Action   = ["lambda:InvokeFunction"],
      Resource = "${module.signer_lambda.resource_arn}"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "stepfunction_invoke_lambda" {
  role       = module.signer_flow.role_name
  policy_arn = aws_iam_policy.stepfunction_lambda_invoke.arn
}

resource "aws_iam_policy" "eventbridge_invoke_stepfunction" {
  name        = "${local.project_env}-eventbridge-invoke-stepfunction"
  description = "Permite ao EventBridge invocar a Step Function"
  tags        = local.tags
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect   = "Allow",
      Action   = ["states:StartExecution"],
      Resource = module.signer_flow.resource_arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "eventbridge_invoke_stepfunction" {
  role       = module.s3_object_created_eventbridge.role_name
  policy_arn = aws_iam_policy.eventbridge_invoke_stepfunction.arn
}
