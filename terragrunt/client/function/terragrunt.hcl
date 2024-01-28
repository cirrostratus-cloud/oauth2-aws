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
        "Resource": "arn:aws:logs:*:*:*" # TODO: Restrict to resource creating log group outside of this module
      }
    ]
  })
  environment_variables = {
    LOG_LEVEL = "INFO"
  }
  zip_location = "${get_parent_terragrunt_dir()}/dist/client/client.zip"
  common_tags = local.common_tags
}
