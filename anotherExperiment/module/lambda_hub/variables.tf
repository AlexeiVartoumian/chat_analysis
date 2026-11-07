variable "aws_iam_role_main_arn" {
    type = string
}

variable "sqs_coordinator_arn" {
    type = string
}

variable "sqs_coordinator_smart_sqs_queue_id" {
    type = string
}
variable "sqs_cordinator_smart_arn" {
    type = string
}

variable "sqs_coordinator_bambi_sqs_queue_id" {
    type = string
}
variable "sqs_cordinator_bambi_arn" {
    type = string
}
variable "sqs_coordinator_workaworka_sqs_queue_id" {
    type = string
}
variable "sqs_cordinator_workaworka_arn" {
    type = string
}

variable "sqs_cordinator_green_arn" {
    type =string
}

variable "sqs_coordinator_green_sqs_queue_id" {
    type = string 
}

variable "sqs_cordinator_ash_arn" {
    type =string
}

variable "sqs_coordinator_ash_sqs_queue_id" {
    type = string 
}

variable "sqs_cordinator_brez_arn" {
    type =string
}

variable "sqs_coordinator_brez_sqs_queue_id" {
    type = string 
}

variable "sqs_cordinator_recruit_arn" {
    type =string
}

variable "sqs_coordinator_recruit_sqs_queue_id" {
    type = string 
}

variable "sqs_cordinator_tailor_arn" {
    type =string
}

variable "sqs_coordinator_tailor_sqs_queue_id" {
    type = string 
}
variable "sqs_cordinator_rip_arn" {
    type =string
}

variable "sqs_coordinator_rip_sqs_queue_id" {
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

variable "account_pool_table_orac_arn"{
    type = string
}

variable "account_pool_table_smart_arn"{
    type = string
}

variable "account_pool_table_bambi_arn"{
    type = string
}

variable "account_pool_table_workaworka_arn"{
    type = string
}

# variable "account_pool_table_stream_arn"{
#     type = string
# }

variable "account_pool_table_lev"{
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


variable "s3_output_bucket_smart_store_name" {
    type = string
}

variable "s3_output_bucket_orac_store_name" {
    type = string
}

variable "s3_output_bucket_bambi_store_name" {
    type = string
}

variable "s3_output_bucket_workaworka_store_name" {
    type = string
}

variable "s3_output_bucket_green_cache_name" {
    type = string
}

variable "s3_output_bucket_ash_cache_name" {
    type = string
}

variable "s3_output_bucket_brez_cache_name" {
    type = string
}

variable "s3_output_bucket_recruit_cache_name" {
    type = string
}

variable "s3_output_bucket_tailor_cache_name" {
    type = string
}
variable "s3_output_bucket_rip_cache_name" {
    type = string
}

variable "eventbridge_rule_arn"{
    type = string
}
variable "sqs_coordinator_work_arn" {
    type = string
}


variable "sqs_coordinator_orac_id" {
    type = string
}

variable "sqs_cordinator_orac_arn" {
    type = string
}


variable "polling_enabled" {

    type = bool
    
}