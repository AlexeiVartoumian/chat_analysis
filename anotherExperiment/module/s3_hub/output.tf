

output "s3_bucket_source_arn" {
    description = "s3_bucket_id_source"
    value = aws_s3_bucket.source_store.arn
}

output "s3_bucket_file_arn" {
    description = "s3 bucket file"
    value = aws_s3_bucket.file_store.arn
}

output "s3_bucket_output_arn" {
    description = "s3 bucket file"
    value = aws_s3_bucket.output_store.arn
}

output "s3_bucket_backfill_arn" {
    description = "s3 bucket file"
    value = aws_s3_bucket.backfill_store.arn
}
  

output "s3_bucket_source_id" {
    description = "s3_bucket_id_source"
    value = aws_s3_bucket.source_store.id
}

output "s3_bucket_file_id" {
    description = "s3 bucket file"
    value = aws_s3_bucket.file_store.id
}

output "s3_bucket_output_id" {
    description = "s3 bucket file"
    value = aws_s3_bucket.output_store.id
}
output "s3_bucket_backfill_id" {
    description = "s3 bucket file"
    value = aws_s3_bucket.backfill_store.id
}


output "s3_bucket_source_name" {
    description = "s3_bucket_id_source"
    value = aws_s3_bucket.source_store.bucket
}

output "s3_bucket_file_name" {
    description = "s3 bucket file"
    value = aws_s3_bucket.file_store.bucket
}

output "s3_bucket_output_name" {
    description = "s3 bucket file"
    value = aws_s3_bucket.output_store.bucket
}

output "s3_bucket_backfill_name" {
    description = "s3 bucket file"
    value = aws_s3_bucket.backfill_store.bucket
}

