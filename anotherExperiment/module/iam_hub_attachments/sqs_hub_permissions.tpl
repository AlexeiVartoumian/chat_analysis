{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "sqs:ReceiveMessage",
                "sqs:DeleteMessage",
                "sqs:GetQueueAttributes",
                "sqs:SendMessage",
                "sqs:ChangeMessageVisibility"
            ],
            "Resource": [
                "arn:aws:sqs:eu-west-2:${aws_account}:${sqs_deadletter_arn}",
                "arn:aws:sqs:eu-west-2:${aws_account}:${sqs_workflow_coordinator}"
                "arn:aws:sqs:eu-west-2:${aws_account}:${sqs_workflow_coordinator_deed}"
                "arn:aws:sqs:eu-west-2:${aws_account}:${sqs_workflow_coordinator_work}"
            ]
        }
    ]
}