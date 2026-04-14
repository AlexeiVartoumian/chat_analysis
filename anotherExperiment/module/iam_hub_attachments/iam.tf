



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



# resource "aws_iam_role_policy" "send_spoke_sqs" {
#   role = var.bucket_reader_role_name
#   policy = templatefile("${path.module}/send_spoke_sqs.tpl" , {
#     spoke_accounts = var.spoke_accounts
#   })
# }

resource "aws_iam_role_policy" "send_spoke_sqs" {
  name = "send_spoke_sqs"
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
          "sqs:SendMessage"
        ]
        Resource = [
          for accnt in var.spoke_accounts :
          "arn:aws:sqs:eu-west-2:${accnt}:workflow-requests-test"
        ]
      }
    ]
  })
}
resource "aws_iam_role_policy" "bucket_permissions" {
  role = var.bucket_reader_role_name
  policy = templatefile("${path.module}/bucket_permissions.tpl" , {
    s3_source_bucket_arn = var.s3_source_bucket_arn
    s3_file_bucket_arn = var.s3_file_bucket_arn
    s3_output_bucket_arn = var.s3_output_bucket_arn
  })
}

resource "aws_iam_role_policy" "dynamodb_permissions" {
  role = var.bucket_reader_role_name
  policy = templatefile("${path.module}/dynamodb.tpl" , {
    hub_account = var.hub_account
    filepool_table = var.file_pool_table
    accountpool_table = var.account_pool_table
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
          "sqs:SendMessage"
        ]
        Resource = var.sqs_coordinator_arn  # pass the actual SQS ARN directly
      }
    ]
  })
}


resource "aws_iam_role_policy_attachment" "lambda_exec" {
  role = var.bucket_reader_role_name
  policy_arn = data.aws_iam_policy.lambda_basic_execution.arn
}

 
