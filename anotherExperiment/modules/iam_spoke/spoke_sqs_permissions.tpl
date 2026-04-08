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
                "arn:aws:sqs:eu-west-2:${aws_account}:workflow-requests",
                "arn:aws:sqs:eu-west-2:${aws_account}:workflow-lambda2-trigger",
                "arn:aws:sqs:eu-west-2:${aws_account}:workflow-lambda3-trigger"
            ]
        }
    ]
}