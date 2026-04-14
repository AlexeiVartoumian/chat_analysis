
output "aws_iam_role_main_arn" {
    value = aws_iam_role.bucket_reader_main.arn
}

output "aws_iam_role_main_name" {
    value = aws_iam_role.bucket_reader_main.name
}

