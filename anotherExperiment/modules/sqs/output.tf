
output "workflow_requests_sqs_queue_id"{
    description = "sqs_queue_url "
    value = aws_sqs_queue.workflow_requests.id
}

output "workflow_requests sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.workflow_requests.name
}


output "sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.workflow_requests.zero_proof_queue.arn   
}



output "workflow_lambda2_trigger_sqs_queue_id"{
    description = "sqs_queue_url "
    value = aws_sqs_queue.workflow_lambda2_trigger.id
}

output "workflow_lambda2_trigger sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.workflow_lambda2_trigger.name
}


output "workflow_lambda2_trigger_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.workflow_lambda2_trigger.arn   
}

output "workflow_lambda3_trigger_sqs_queue_id"{
    description = "sqs_queue_url "
    value = aws_sqs_queue.workflow_lambda3_trigger.id
}

output "workflow_lambda2_trigger sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.workflow_lambda3_trigger.name
}


output "workflow_lambda2_trigger_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.workflow_lambda3_trigger.arn   
}