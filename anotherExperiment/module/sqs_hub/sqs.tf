
resource "aws_sqs_queue" "sqs_hub_requests" {
  name                      = "workflow-cordinator-test"
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 1500
  tags = {
    Environment = "production"
  }
}

resource "aws_sqs_queue_policy" "sqs_hub_requests"{
    queue_url = aws_sqs_queue.sqs_hub_requests.id
    policy = templatefile("${path.module}/sqs_access.tpl" ,{
        aws_account  = data.aws_caller_identity.current.account_id
        sqs_queuename  = aws_sqs_queue.sqs_hub_requests.name
    })
}