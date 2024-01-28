terraform {
  backend "s3" {}
}

locals {
  iam_role_name = "${var.module_name}-${var.function_name}-execution-role"
  iam_policy_name = "${var.module_name}-${var.function_name}-policy"
}

resource "aws_iam_role" "function" {
  name = local.iam_role_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
  tags = var.common_tags
}

resource "aws_iam_role_policy" "function" {
  name = local.iam_policy_name
  role = aws_iam_role.function.id
  policy = var.iam_policy
}

resource "aws_lambda_function" "function" {
  function_name = "${var.module_name}-${var.function_name}"
  role = aws_iam_role.function.arn
  filename = var.zip_location
  runtime = "provided.al2"
  handler = "bootstrap"
  memory_size = var.memory_size
  timeout = var.timeout
  environment {
    variables = var.environment_variables
  }
  lifecycle {
    ignore_changes = [environment]
  }
  tags = var.common_tags
}

resource "aws_cloudwatch_log_group" "function" {
  name = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 30
  tags = var.common_tags
}
