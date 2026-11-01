
variable "spoke_accounts"{
    type = list(string)
    default = []
}


variable "hub_account" {
    type = string
}

variable "admin_role_name" {
  type    = string
}

variable "keys"{
    type = string
}

variable "ami_id" {
    type =string 
}

variable "polling_enabled" {
  type        = bool
  description = "Turn SQS polling on or off for all spokes"
  nullable    = false
}