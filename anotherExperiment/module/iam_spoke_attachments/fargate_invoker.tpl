{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "RunECSTasks",
            "Effect": "Allow",
            "Action": "ecs:RunTask",
            "Resource":  "arn:aws:ecs:eu-west-2:${spoke_account_id}:task-definition/scroller-task"

        },
        {
            "Sid": "PassTaskRoles",
            "Effect": "Allow",
            "Action": [
               "iam:PassRole"
            ],
            "Resource": [
                "arn:aws:iam::${spoke_account_id}:role/ecs_reader_task_execution_role",
                "arn:aws:iam::${spoke_account_id}:role/${bucket_reader_spoke}"
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