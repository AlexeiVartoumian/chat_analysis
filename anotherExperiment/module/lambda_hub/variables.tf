variable "aws_iam_role_main_arn" {
    type = string
}

variable "sqs_coordinator_arn" {
    type = string
}

variable "account_pool_table"{
    type = string
}

variable "file_pool_table"{
    type = string
}


variable "account_pool_table_deed"{
    type = string
}

variable "account_pool_table_work"{
    type = string
}

variable "account_pool_table_work_arn"{
    type = string
}

variable "account_pool_table_stream_arn"{
    type = string
}



variable "file_pool_table_deed"{
    type = string
}

variable "sqs_cordinator_deed_arn" {
    type = string
}

variable "coordinator_work_sqs_queue_id" {
    type = string
}

variable "s3_output_bucket_work_store_name" {
    type = string
}

variable "eventbridge_rule_arn"{
    type = string
}
