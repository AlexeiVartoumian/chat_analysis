output "workflow_requests_sqs_queue_id"{
    description = "sqs_queue_url "
    value = aws_sqs_queue.sqs_hub_requests.id
}

output "workflow_requests_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests.name
}


output "coordinator_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests.arn   
}
