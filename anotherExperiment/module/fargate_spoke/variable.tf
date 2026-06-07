

variable "task_role_arn" {
    type = string
    description = "responsible for giveing the permissions needed during runtime"
}

variable "execution_role_arn" {
    type = string 
    description = "responsible for getting the container up and running"
}



variable "iam_role_main_arn" {
    type = string

    description = "main role arn"
}