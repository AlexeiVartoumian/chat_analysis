



# resource "aws_iam_role_policy" "bucket_reader_trust" {
#   name = "bucket_reader_trust"
#   role = var.bucket_reader_role_name  # passed in from iam_hub module output

#   policy = templatefile("${path.module}/assume.tpl", {
#     spoke_accounts = var.spoke_accounts
#      iam_role_spoke  = var.iam_role_spoke
#   })
# }

# resource "aws_iam_role_policy" "bucket_reader_trust" {
#   name = "bucket_reader_trust"
#   role = var.bucket_reader_role_name

#   policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Effect = "Allow"
#         Principal = {
#           Service = "lambda.amazonaws.com"
#           AWS = [
#             for acct in var.spoke_accounts :
#             "arn:aws:iam::${acct}:role/${var.iam_role_spoke}"
#           ]
#         }
#         Action = "sts:AssumeRole"
#       }
#     ]
#   })
# }


resource "aws_iam_role_policy" "send_spoke_sqs" {
  role = var.bucket_reader_role_name
  policy = jsonencode({
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
                for account_id in var.spoke_accounts:
                "arn:aws:sqs:eu-west-2:${account_id}:workflow-requests-test"
            ]
        }
    ]

  })
}

resource "aws_iam_role_policy" "esc_invoke" {
  name = "assume-ecs-invoker-role"
  role = var.bucket_reader_role_name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Resource = [
          for account_id in var.spoke_accounts :
          "arn:aws:iam::${account_id}:role/ecs_task_invoker_role"
        ]
      }
    ]
  })
}
variable "worker_role_tag" {
  description = "Value of the Role tag used to scope self-termination to this worker pool"
  type        = string
  default     = "time-to-go"
}

data "aws_iam_policy_document" "self_terminate" {
  statement {
    effect    = "Allow"
    actions   = ["ec2:TerminateInstances"]
    resources = ["*"]

    condition {
      test     = "StringEquals"
      variable = "ec2:ResourceTag/Role"
      values   = [var.worker_role_tag]
    }
  }
}
resource "aws_iam_role_policy" "time-to-go" {
  name = "time-2-go"
  role = var.bucket_reader_role_name
  policy = data.aws_iam_policy_document.self_terminate.json
}
# resource "aws_iam_role_policy" "send_spoke_sqs" {
#   role = var.bucket_reader_role_name
#   policy = templatefile("${path.module}/send_spoke_sqs.tpl" , {
#     spoke_accounts = var.spoke_accounts
#   })
# }




resource "aws_iam_role_policy" "bucket_permissions" {
  role = var.bucket_reader_role_name
  policy = templatefile("${path.module}/bucket_permissions.tpl" , {
    s3_source_bucket_arn = var.s3_source_bucket_arn
    s3_file_bucket_arn = var.s3_file_bucket_arn
    s3_output_bucket_arn = var.s3_output_bucket_arn
    s3_backfill_bucket_arn = var.s3_backfill_bucket_arn
    s3_output_bucket_ash_cache_arn = var.s3_output_bucket_ash_cache_arn
    s3_output_bucket_ash_store_arn = var.s3_output_bucket_ash_store_arn
    s3_output_bucket_deed_cache_arn = var.s3_output_bucket_deed_cache_arn
    s3_output_bucket_deed_store_arn = var.s3_output_bucket_deed_store_arn
    s3_output_bucket_green_cache_arn = var.s3_output_bucket_green_cache_arn
    s3_output_bucket_green_store_arn = var.s3_output_bucket_green_store_arn
    s3_output_bucket_work_cache_arn = var.s3_output_bucket_work_cache_arn
    s3_output_bucket_work_store_arn = var.s3_output_bucket_work_store_arn
  })
}

resource "aws_iam_role_policy" "dynamodb_permissions" {
  role = var.bucket_reader_role_name
  policy = templatefile("${path.module}/dynamodb.tpl" , {
    hub_account = var.hub_account
    filepool_table = var.file_pool_table
    accountpool_table = var.account_pool_table
    filepool_table_deed = var.file_pool_table_deed
    accountpool_table_deed = var.account_pool_table_deed
    file_pool_table_ash = var.file_pool_table_ash
    file_pool_table_green = var.file_pool_table_green
    accountpool_table_work = var.account_pool_table_work
    account_pool_table_stream_arn = var.account_pool_table_stream_arn
  })
}

#TODO add sqs queue for hub
# resource "aws_iam_role_policy" "sqs_permissions" {
#   role = aws_iam_role.bucket_reader_main.name
#   policy = templatefile("${path.module}/sqs_permissions.tpl" , {
#    aws_account = var.hub_account
#   })
# }



# resource "aws_iam_role_policy" "sqs_coordinator_permissions" {
#   role = var.bucket_reader_role_name
#   policy = templatefile("${path.module}/sqs_hub_permissions.tpl" , {
#     aws_account  = data.aws_caller_identity.current.account_id
#     sqs_workflow_coordinator = var.bucket_reader_main_arn
#   })
# }

resource "aws_iam_role_policy" "sqs_coordinator_permissions" {
  role = var.bucket_reader_role_name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "sqs:SendMessage",
          "sqs:ChangeMessageVisibility"
        ]
        Resource = [ var.sqs_coordinator_arn , var.sqs_deadletter_arn , var.sqs_coordinator_arn_deed , var.sqs_coordinator_arn_work]  # pass the actual SQS ARN directly
      }
    ]
  })
}


resource "aws_iam_role_policy_attachment" "lambda_exec" {
  role = var.bucket_reader_role_name
  policy_arn = data.aws_iam_policy.lambda_basic_execution.arn
}


resource "aws_iam_policy" "hub_assume_spoke_invoker" {
  name = var.bucket_reader_role_name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AssumeSpokeinvokerRole"
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Resource = [
          for acct in var.spoke_accounts :
          "arn:aws:iam::${acct}:role/ecs_task_invoker_role"
        ]
      }
    ]
  })
}
 
