
resource "aws_s3_bucket" "transcribe_store_logs" {
      bucket = "transcribefile-access-logs-${data.aws_caller_identity.current.account_id}"

}

resource "aws_s3_bucket_lifecycle_configuration" "transcribe_store_logs" {
  bucket = aws_s3_bucket.transcribe_store_logs.id

  rule {
    id     = "expire-logs"
    status = "Enabled"

    filter {}

    expiration {
      days = 30
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "transcribe_store_logs" {
  bucket = aws_s3_bucket.transcribe_store_logs.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "transcribe_store_logs" {
  depends_on = [aws_s3_bucket_ownership_controls.transcribe_store_logs]
  bucket     = aws_s3_bucket.transcribe_store_logs.id
  acl        = "log-delivery-write"
}


resource "aws_s3_bucket_logging" "transcribe_store_logs" {
  bucket        = "alexeitranscribefile"
  target_bucket = aws_s3_bucket.source_store_logs.id
  target_prefix = "logs/"
}


