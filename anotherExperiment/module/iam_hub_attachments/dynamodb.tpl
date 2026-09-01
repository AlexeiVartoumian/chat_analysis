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
                "dynamodb:DeleteItem",
                "dynamodb:GetItem",
                "dynamodb:GetRecords",
                "dynamodb:GetShardIterator",
                "dynamodb:DescribeStream",
                "dynamodb:ListStreams"
            ],
            "Resource": [
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${accountpool_table}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table}/index/status-index",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table_deed}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${accountpool_table_deed}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${filepool_table_deed }/index/status-index",

                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${file_pool_table_ash}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${file_pool_table_ash}/index/status-index",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${file_pool_table_green}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${file_pool_table_green}/index/status-index",

                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${accountpool_table_work}",
                "arn:aws:dynamodb:eu-west-2:${hub_account}:table/${accountpool_table_work}/index/status-index",
                "${account_pool_table_stream_arn}"


                
            ]
        }
    ]
}
