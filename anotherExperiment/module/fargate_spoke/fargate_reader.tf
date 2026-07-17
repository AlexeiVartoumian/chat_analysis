resource "aws_ecs_cluster" "reader_cluster" {
  name = "reader-cluster"
}

resource "aws_ecs_task_definition" "reader_task" {
  family                   = "reader-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  task_role_arn            = var.task_role_arn
  execution_role_arn       = var.execution_role_arn

  container_definitions = jsonencode([
    {
      name      = "reader-container"
      image     = "390746273208.dkr.ecr.eu-west-2.amazonaws.com/reader:latest"
      essential = true
      command   = [
     
      ]
      # environment = [
      #   { name = "s3_source_bucket", value = var.s3_source_name},
      #   { name = "file_store", value = var.s3_output_store_name},
      #   { name = "file_id", value = "cookies-grouped-remi"},
       
      # ]
      environment = [
        { name = "RoleArn", value = var.iam_role_main_arn},
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.reader_task_spoke.name
          awslogs-region        = "eu-west-2"
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

}


resource "aws_cloudwatch_log_group" "reader_task_spoke" {
  name              = "/ecs/backfill-task"
  retention_in_days = 7
}