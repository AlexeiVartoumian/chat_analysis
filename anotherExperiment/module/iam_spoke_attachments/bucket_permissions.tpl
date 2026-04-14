{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Statement1",
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetObject",
                "s3:PutObject"
            ],
            "Resource": [
                
                "${s3_source_bucket_arn}",
                "${s3_source_bucket_arn}/*",
                "${s3_file_bucket_arn}",
                "${s3_file_bucket_arn}/*",
                "${s3_output_bucket_arn}",
                "${s3_output_bucket_arn}/*"
              
            ]
        }
    ]
}