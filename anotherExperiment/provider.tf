

provider "aws" {
  region = "eu-west-2"
}


//TODO add assume role and jinja this ?
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  alias = "hub"
  region = "eu-west-2"

  assume_role {
    role_arn = "arn:aws:iam::${var.hub_account}:role/${var.admin_role_name}"
  }
}

# provider "aws"{
#   alias  = "spoke"
#   region = "eu-west-2"

#   assume_role {
#     role_arn = "arn:aws:iam::${var.spoke_accounts[0]}:role/${var.admin_role_name}"
#   }
# }
