

resource "aws_sqs_queue" "workflow_requests" {
  name                      = "workflow-requests"
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
  name                      = "workflow-lambda2-trigger"
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
  name                      = "workflow-lambda3-trigger"
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  visibility_timeout_seconds = 1500
  tags = {
    Environment = "production"
  }
}