output "lambda_reader" {
    description = "lambda reader"
    value = aws_lambda_function.reader.arn
}

output "lambda_processor" {
    description = "lambda processor"
    value = aws_lambda_function.processor.arn   
}

output "lambda_gometadata" {
    description = "lambda go metadata"
    value = aws_lambda_function.go_metadata.arn 
}