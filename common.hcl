locals {
  module_name = "cirrostratus-oauth2"
  aws_region = get_env("AWS_DEFAULT_REGION")
  aws_stage = get_env("AWS_STAGE")
  module_bucket = "${local.module_name}-${local.aws_region}"
  common_tags = {
    module = local.module_name
  }
}