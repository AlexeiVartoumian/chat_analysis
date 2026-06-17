
resource "aws_s3_bucket" "output_store_logs" {
  bucket = "${local.s3_bucket_name_output}-access-logs-${data.aws_caller_identity.current.account_id}"
}

resource "aws_s3_bucket_lifecycle_configuration" "output_store_logs_lifecycle" {
  bucket = aws_s3_bucket.output_store_logs.id

  rule {
    id     = "expire-logs"
    status = "Enabled"

    filter {}

    expiration {
      days = 30
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "output_store_logs" {
  bucket = aws_s3_bucket.output_store_logs.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "output_store_logs" {
  depends_on = [aws_s3_bucket_ownership_controls.output_store_logs]
  bucket     = aws_s3_bucket.output_store_logs.id
  acl        = "log-delivery-write"
}



resource "aws_s3_bucket_logging" "output_store" {
  bucket        = aws_s3_bucket.output_store.id
  target_bucket = aws_s3_bucket.output_store_logs.id
  target_prefix = "logs/"
}



resource "aws_s3_bucket" "backfill_store_logs" {
 
  bucket = "${local.s3_bucket_name_backfill}-access-logs-${data.aws_caller_identity.current.account_id}"

}

resource "aws_s3_bucket_lifecycle_configuration" "backfill_store_logs_lifecycle" {
  bucket = aws_s3_bucket.backfill_store_logs.id

  rule {
    id     = "expire-logs"
    status = "Enabled"

    filter {}

    expiration {
      days = 30
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "backfill_store_logs" {
  bucket = aws_s3_bucket.backfill_store_logs.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "backfill_store_logs" {
  depends_on = [aws_s3_bucket_ownership_controls.backfill_store_logs]
  bucket     = aws_s3_bucket.backfill_store_logs.id
  acl        = "log-delivery-write"
}


resource "aws_s3_bucket_logging" "backfill_store" {
  bucket        = aws_s3_bucket.backfill_store.id
  target_bucket = aws_s3_bucket.backfill_store_logs.id
  target_prefix = "logs/"
}


resource "aws_s3_bucket" "file_store_logs" {
 
      bucket = "${local.s3_bucket_name_file}-access-logs-${data.aws_caller_identity.current.account_id}"
}

resource "aws_s3_bucket_lifecycle_configuration" "file_store_logs_lifecycle" {
  bucket = aws_s3_bucket.file_store_logs.id

  rule {
    id     = "expire-logs"
    status = "Enabled"

    filter {}

    expiration {
      days = 30
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "file_store_logs" {
  bucket = aws_s3_bucket.file_store_logs.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "file_store_logs" {
  depends_on = [aws_s3_bucket_ownership_controls.file_store_logs]
  bucket     = aws_s3_bucket.file_store_logs.id
  acl        = "log-delivery-write"
}


resource "aws_s3_bucket_logging" "file_store" {
  bucket        = aws_s3_bucket.file_store.id
  target_bucket = aws_s3_bucket.file_store_logs.id
  target_prefix = "logs/"
}



resource "aws_s3_bucket" "source_store_logs" {
      bucket = "${local.s3_bucket_name_source}-access-logs-${data.aws_caller_identity.current.account_id}"

}

resource "aws_s3_bucket_lifecycle_configuration" "source_store_logs_lifecycle" {
  bucket = aws_s3_bucket.source_store_logs.id

  rule {
    id     = "expire-logs"
    status = "Enabled"

    filter {}

    expiration {
      days = 30
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "source_store_logs" {
  bucket = aws_s3_bucket.source_store_logs.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "source_store_logs" {
  depends_on = [aws_s3_bucket_ownership_controls.source_store_logs]
  bucket     = aws_s3_bucket.source_store_logs.id
  acl        = "log-delivery-write"
}


resource "aws_s3_bucket_logging" "source_store" {
  bucket        = aws_s3_bucket.source_store.id
  target_bucket = aws_s3_bucket.source_store_logs.id
  target_prefix = "logs/"
}


