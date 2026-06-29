data "aws_vpc" "vpc" {
#  filter {
#   name = "tag:myvpc"
#   values = ["project-vpc"]
#  }
id = "vpc-06f5aaa7886920c69"
}
data "aws_subnets" "default_subnets" {
 filter {
  name = "vpc-id"
  values = [data.aws_vpc.vpc.id]
 }
}
data "aws_subnet" "public_subnet" {
 id = data.aws_subnets.default_subnets.ids[0] 
}

data "aws_security_group" "sg" {
  name = "test"
}




