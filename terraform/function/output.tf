output function_name {
  value       = aws_lambda_function.function.function_name
  description = "Function name"
}

output invoke_arn  {
  value       = aws_lambda_function.function.invoke_arn
  description = "Function invoke ARN"
}
