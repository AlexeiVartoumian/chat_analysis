
resource "aws_iam_role" "bucket_reader_main" {
  name = "the_bucket_dealer"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
          AWS = [
            for acct in var.spoke_accounts :
            "arn:aws:iam::${acct}:role/the_bucket_dealer_spoke"
          ]
        }
        Action = "sts:AssumeRole"
       }
       #,
      #   {
      #   Action = "sts:AssumeRole"
      #   Effect = "Allow"
      #   Principal = {
      #     Service = "ecs-tasks.amazonaws.com"
      #   }
      # }
    ]
  })
}
# resource "aws_iam_role" "bucket_reader_main" {

#     name = "the_bucket_dealer"
#     # assume_role_policy = templatefile("${path.module}/assume.tpl", 
#     # { spoke_accounts = var.spoke_accounts 
#     # })
#     assume_role_policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Effect    = "Allow"
#         Principal = { Service = "lambda.amazonaws.com" }
#         Action    = "sts:AssumeRole"
#       }
#     ]
#   })

# }



 
