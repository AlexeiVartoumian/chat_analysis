variable "iam_role_arn_spoke" {
    description = "role arn used by lambdas"
    type = string
}

variable "iam_role_main_arn" {
    type = string 
}

variable "sqs_request_access_arn" {
    type =string
}

variable "sqs_queue_2_arn" {
    type =string
}
variable "sqs_queue_3_arn" {
    type =string
}

variable "s3_source_name"{
    type = string
}

variable "s3_filestore_name"{
    type = string
}

variable "sqs_queue_2_id" {
    type = string
}

variable "s3_output_bucket_name" {
    type = string
}

variable "sqs_queue_3_id" { 
    type = string
}

variable "dynamodb_filetable_name"{
    type = string
}

variable dynamodb_accounttable_name{
    type = string
}