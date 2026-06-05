
#makes sense if it was private vpcs but instead will grab vpc params w. boto3 populate db 
# and then query. lol still need it for security group 
data "aws_vpc" "vpc" {
    filter {
    name   = "isDefault"
    values = ["true"]
  }
}

data "aws_subnets" "default_subnets" {

    filter {
        name = "vpc-id"
        values = [data.aws_vpc.vpc.id]
    }
}

data "aws_subnet" "public_subnet"{
    id = data.aws_subnets.default_subnets.ids[0]
}

