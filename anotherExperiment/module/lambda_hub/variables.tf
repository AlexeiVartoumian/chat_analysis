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

variable "file_pool_table_deed"{
    type = string
}

variable "sqs_cordinator_deed_arn" {
    type = string
}