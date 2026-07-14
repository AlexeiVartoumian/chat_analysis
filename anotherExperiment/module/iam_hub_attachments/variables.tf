variable bucket_reader_role_name {
    description = "value"
    type =string
}

variable bucket_reader_main_arn {
     description = "value"
    type =string
}


variable spoke_accounts {
    description = "value"
    type = list(string)
}

variable hub_account{
    type = string
}

variable s3_source_bucket_arn {
    type = string
}

variable s3_file_bucket_arn{
    type = string
    description = "the s3 bucket"
}

variable s3_output_bucket_arn{
    type = string
    description = "ok"
}

variable  s3_backfill_bucket_arn{
    type = string
    description = "v"
}
               

variable file_pool_table {
    type = string 
} 
variable account_pool_table{
    type = string 
}

variable file_pool_table_deed {
    type = string 
} 
variable account_pool_table_deed{
    type = string 
}

variable sqs_coordinator_arn {
    type = string
}

variable sqs_coordinator_arn_deed {
    type = string
}

variable sqs_deadletter_arn {
    type = string
}


variable s3_output_bucket_ash_cache_arn {
    type =string
}

variable s3_output_bucket_ash_store_arn {
    type = string
}
variable s3_output_bucket_deed_cache_arn {
    type =string
}

variable s3_output_bucket_deed_store_arn {
    type = string
}

variable s3_output_bucket_green_cache_arn {
    type =string
}

variable s3_output_bucket_green_store_arn {
    type = string
}