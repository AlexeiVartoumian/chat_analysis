
output "iam_spoke_role_arn" {
    value = aws_iam_role.bucket_reader_spoke.arn
}

output "iam_spoke_role_name" {
    value = aws_iam_role.bucket_reader_spoke.name
}