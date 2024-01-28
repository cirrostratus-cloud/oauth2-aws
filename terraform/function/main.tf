terraform {
  backend "s3" {}
}

locals {
  iam_role_name = "${var.module_name}-${var.function_name}-execution-role"
  iam_policy_name = "${var.module_name}-${var.function_name}-policy"
  zip_name_location = "${var.zip_location}/${var.zip_name}"
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

data "archive_file" "function" {
  type = "zip"
  source_dir  = var.file_location
  output_path = local.zip_name_location
}

data "aws_s3_bucket" "function" {
  bucket = var.module_bucket
}

resource "aws_s3_object" "function" {
  bucket = data.aws_s3_bucket.function.id
  key    = var.zip_name
  source = data.archive_file.function.output_path
  etag = filemd5(data.archive_file.function.output_path)
}

resource "aws_lambda_function" "function" {
  function_name = "${var.module_name}-${var.function_name}"
  role = aws_iam_role.function.arn
  s3_bucket = data.aws_s3_bucket.function.id
  s3_key    = aws_s3_object.function.key
  source_code_hash = data.archive_file.function.output_base64sha256
  runtime = "provided.al2"
  handler = "bootstrap"
  memory_size = var.memory_size
  timeout = var.timeout
  environment {
    variables = var.environment_variables
  }
  tags = var.common_tags
}

resource "aws_cloudwatch_log_group" "function" {
  name = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 30
  tags = var.common_tags
}