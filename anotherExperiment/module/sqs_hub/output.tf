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


output "deadletter_requests_sqs_queue_id"{
    description = "sqs_queue_url "
    value = aws_sqs_queue.deadletter_requests.id
}

output "deadletter_requests_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.deadletter_requests.name
}


output "deadletter_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.deadletter_requests.arn   
}

output "coordinator_deed_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_deed.arn   
}


output "coordinator_work_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_work.arn   
}

output "coordinator_work_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_work.id  
}

output "coordinator_work_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_work.name
}

output "coordinator_lev_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_lev.arn   
}

output "coordinator_lev_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_lev.id  
}

output "coordinator_lev_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_lev.name
}


output "coordinator_orac_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_orac.arn   
}

output "coordinator_orac_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_orac.id  
}

output "coordinator_orac_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_orac.name
}

output "coordinator_smart_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_smart.arn   
}

output "coordinator_smart_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_smart.id  
}

output "coordinator_smart_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_smart.name
}

output "coordinator_bambi_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_bambi.arn   
}

output "coordinator_bambi_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_bambi.id  
}

output "coordinator_bambi_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_bambi.name
}

output "coordinator_workaworka_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_workaworka.arn   
}

output "coordinator_workaworka_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_workaworka.id  
}

output "coordinator_workaworka_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_workaworka.name
}

output "coordinator_green_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_green.arn   
}

output "coordinator_green_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_green.id  
}

output "coordinator_green_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_green.name
}

output "coordinator_ash_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_ash.arn   
}

output "coordinator_ash_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_ash.id  
}

output "coordinator_ash_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_ash.name
}


output "coordinator_brez_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_brez.arn   
}

output "coordinator_brez_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_brez.id  
}

output "coordinator_brez_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_brez.name
}

output "coordinator_recruit_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_recruit.arn   
}

output "coordinator_recruit_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_recruit.id  
}

output "coordinator_recruit_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_recruit.name
}

output "coordinator_tailor_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_tailor.arn   
}

output "coordinator_tailor_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_tailor.id  
}

output "coordinator_tailor_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_tailor.name
}

output "coordinator_rip_sqs_queue_arn" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_rip.arn   
}

output "coordinator_rip_sqs_queue_id" {
    description = "sqs arn"
    value = aws_sqs_queue.sqs_hub_requests_rip.id  
}

output "coordinator_rip_sqs_queue_name" {
    description = "sqs_queue_name"
    value = aws_sqs_queue.sqs_hub_requests_rip.name
}
