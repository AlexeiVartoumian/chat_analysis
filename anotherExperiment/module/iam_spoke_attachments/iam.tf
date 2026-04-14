




resource "aws_iam_role_policy" "sqs" {

    role = var.bucket_reader_spoke
    policy = templatefile("${path.module}/spoke_sqs_permissions.tpl" , {
         aws_account  = data.aws_caller_identity.current.account_id
         sqs_queue_request = var.sqs_request_access
         sqs_queue_2 = var.sqs_queue_2
         sqs_queue_3 = var.sqs_queue_3
  })
}
resource "aws_iam_role_policy" "assume_hub_role" {

    role = var.bucket_reader_spoke
    policy = templatefile("${path.module}/spoke_iam_assume.tpl" , {
        hub_account = var.hub_account
  })
}

resource "aws_iam_role_policy_attachment" "lambda_exec" {
  role = var.bucket_reader_spoke
  policy_arn = data.aws_iam_policy.lambda_basic_execution.arn
}

