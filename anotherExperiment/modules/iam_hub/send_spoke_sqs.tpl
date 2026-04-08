{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "sqs:ReceiveMessage",
                "sqs:DeleteMessage",
                "sqs:GetQueueAttributes",
                "sqs:SendMessage"
            ],
            "Resource": [
                ${jsonencode([
                    for accnt in spoke_accounts:
                    "arn:aws:sqs:eu-west-2:${accnt}:workflow-requests"
                ])}
                
               
            ]
        }
    ]
}