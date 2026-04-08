{
    "Version": "2012-10-17",
    "Statement": [     
        {
            "Effect": "Allow",
            "Action": [
                "dynamodb:UpdateItem",
                "dynamodb:DescribeTable",
                "dynamodb:Scan",
                "dynamodb:PutItem",
                "dynamodb:Query"
            ],
            "Resource": [
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/filepool",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/accountpool",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/filepool/index/status-index""
            ]
        }
    ]
}