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
                "dynamodb:Query",
                "dynamodb:DeleteItem"
            ],
            "Resource": [
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${accountpool_table}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table}/index/status-index",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table_deed}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${accountpool_table_deed}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table_deed }/index/status-index"
            ]
        }
    ]
}
