

resource "aws_iam_role" "bucket_reader_spoke" {

    name = "the_bucket_dealer_spoke"
    assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = { Service = "lambda.amazonaws.com" }
        Action    = "sts:AssumeRole"
      },
    {
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }
    ]
  })

}



resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecs_reader_task_execution_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role" "fargate_spoke_invoker_role" {
  name = "ecs_task_invoker_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Principal = {
          #AWS = "arn:aws:iam::${var.hub_account}:role/hub_lambda_execution_role"
          AWS = var.iam_role_main_arn
        }
      }
    ]
  })
}
