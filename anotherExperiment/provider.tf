

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