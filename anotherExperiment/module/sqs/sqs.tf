

resource "aws_sqs_queue" "workflow_requests" {
  name                      = "workflow-requests-test"
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 1500
  tags = {
    Environment = "production"
  }
}

resource "aws_sqs_queue" "workflow_lambda2_trigger" {
  name                      = "workflow-lambda2-trigger-test"
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 1500
  tags = {
    Environment = "production"
  }
}

resource "aws_sqs_queue" "workflow_lambda3_trigger" {
  name                      = "workflow-lambda3-trigger-test"
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 1500
  tags = {
    Environment = "production"
  }
}

resource "aws_sqs_queue_policy" "request_access"{
    queue_url = aws_sqs_queue.workflow_requests.id
    policy = templatefile("${path.module}/sqs_access_request.tpl" ,{
        aws_account  = data.aws_caller_identity.current.account_id
        hub_account  = var.hub_account
        sqs_queuename  = aws_sqs_queue.workflow_requests.name
        orchestrator = var.aws_iam_role_main_name
    })
}

resource "aws_sqs_queue_policy" "access_lambda_2"{
    queue_url = aws_sqs_queue.workflow_lambda2_trigger.id
    policy = templatefile("${path.module}/sqs_access_request.tpl" ,{
        aws_account  = data.aws_caller_identity.current.account_id
        hub_account  = var.hub_account
        sqs_queuename  = aws_sqs_queue.workflow_lambda2_trigger.name
        orchestrator = var.aws_iam_role_main_name
    })
}

resource "aws_sqs_queue_policy" "access_lambda_3"{
    queue_url = aws_sqs_queue.workflow_lambda3_trigger.id
    policy = templatefile("${path.module}/sqs_access_request.tpl" ,{
        aws_account  = data.aws_caller_identity.current.account_id
        hub_account  = var.hub_account
        sqs_queuename  = aws_sqs_queue.workflow_lambda3_trigger.name
        orchestrator = var.aws_iam_role_main_name
    })
}