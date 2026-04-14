

resource "aws_iam_role" "bucket_reader_spoke" {

    name = "the_bucket_dealer_spoke"
    assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = { Service = "lambda.amazonaws.com" }
        Action    = "sts:AssumeRole"
      }
    ]
  })

}
