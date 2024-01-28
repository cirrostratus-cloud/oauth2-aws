variable module_name {
  type        = string
  description = "Module name"
}

variable function_name {
  type        = string
  description = "Function name"
}

variable routes {
  type        = list(object({
    path    = string
    method = string
  }))
  description = "Routes"
}

variable invoke_arn {
  type        = string
  description = "Invoke ARN"
}

variable common_tags {
  type        = map(string)
  description = "Common tags"
}


