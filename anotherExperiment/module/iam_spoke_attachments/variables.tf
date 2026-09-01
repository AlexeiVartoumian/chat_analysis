
variable "sqs_request_access" {
    type = string
}

variable "sqs_queue_2" {
    type = string
}

variable "sqs_queue_3" {
    type = string
}

variable "sqs_workd_queue" {
    type = string
}

variable "bucket_reader_spoke" {
    type = string
}

variable "fargate_spoke" {
    type = string
}

variable "hub_account" {
    type = string
}



variable "fargate_invoker_name" {
    type = string 
    description= "this guy is what is assumed to actually run the task"
} 
variable "fargate_invoker_arn" {
    type = string
}