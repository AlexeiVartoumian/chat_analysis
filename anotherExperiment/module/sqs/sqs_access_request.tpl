{
  "Version": "2012-10-17",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__owner_statement",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${aws_account}:root"
      },
      "Action": "SQS:*",
      "Resource": "arn:aws:sqs:eu-west-2:${aws_account}:${sqs_queuename}"
    },
    {
      "Sid": "__sender_statement",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${hub_account}:role/${orchestrator}"
      },
      "Action": "SQS:SendMessage",
      "Resource": "arn:aws:sqs:eu-west-2:${aws_account}:${sqs_queuename}"
    }
  ]
}