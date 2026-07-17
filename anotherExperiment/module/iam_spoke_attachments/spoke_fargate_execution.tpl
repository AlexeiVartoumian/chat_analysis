{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AuthToHubECR",
            "Effect": "Allow",
            "Action": "ecr:GetAuthorizationToken",
            "Resource": "*"
        },
        {
            "Sid": "PullFromHubECR",
            "Effect": "Allow",
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability"
            ],
            "Resource": ["arn:aws:ecr:eu-west-2:${hub_account}:repository/scroller",
                         "arn:aws:ecr:eu-west-2:${hub_account}:repository/reader"  
             ]
        },
        {
            "Sid": "WriteLogs",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "arn:aws:logs:*:*:log-group:/ecs/*"
    }
    ]
}