{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com",
                "AWS": ${jsonencode([
                    for acct in spoke_accounts :
                        "arn:aws:iam::${acct}:role/${iam_role_spoke}"              
                ])}
            },
            "Action": "sts:AssumeRole"
        }
    ]
}