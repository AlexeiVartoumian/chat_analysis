

variable "task_role_arn" {
    type = string
    description = "responsible for giveing the permissions needed during runtime"
}

variable "execution_role_arn" {
    type = string 
    description = "responsible for getting the container up and running"
}

