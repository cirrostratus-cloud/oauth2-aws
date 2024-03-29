terraform {
  source = "${get_parent_terragrunt_dir()}/terraform/function"
}

locals {
  common_vars  = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  module_name = local.common_vars.locals.module_name
  function_name = "client"
  common_tags = local.common_vars.locals.common_tags
}

include {
  path = find_in_parent_folders()
}

dependency log {
    config_path = "${get_parent_terragrunt_dir()}/terragrunt/client/http/log"
    mock_outputs = {
        log_arn = "log_arn"
    }
}

inputs = {
  function_name = local.function_name
  module_name = local.module_name
  iam_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        "Resource": "${dependency.log.outputs.log_arn}:*"
      }
    ]
  })
  environment_variables = {
    LOG_LEVEL = "INFO"
    AWS_STAGE = local.common_vars.locals.aws_stage
  }
  module_bucket = local.common_vars.locals.module_bucket
  file_location = "${get_parent_terragrunt_dir()}/bin/client/http"
  zip_location = "${get_parent_terragrunt_dir()}/dist/client/http"
  zip_name = "${local.function_name}.zip"
  common_tags = local.common_tags
}
