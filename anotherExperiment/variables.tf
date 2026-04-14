
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