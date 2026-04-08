variable spoke_accounts {
    description = "value"
    type = list(string)
}

variable hub_account{
    type = string
}

variable s3_cookie_bucket {
    type = string
}
variable s3_file_bucket {
    type = string
}
variable s3_output_bucket {
    type = string
}

variable s3_source_bucket{
    type = string
    description = "the s3 bucket"
}

variable s3_file_bucket{
    type = string
    description = "value"
}

variable s3_output_bucket{
    type = string
    description = "ok"
}