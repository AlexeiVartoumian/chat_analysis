
resource "aws_ecs_cluster" "scroller_cluster" {
    name = "scroller-cluster"
}

resource "aws_ecs_task_definition" "scroller_task" {
  family                   = "scroller-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024"
  memory                   = "2048"
  task_role_arn            = var.task_role_arn
  execution_role_arn       = var.execution_role_arn

  volume  {
    name = "shm_volume"
  }
  container_definitions = jsonencode([
    {
      name      = "scroller-container"
      image     = "390746273208.dkr.ecr.eu-west-2.amazonaws.com/scroller:latest"
      essential = true
      command   = [
     
      ]

     
      mountPoints = [
        {
          sourceVolume = "shm_volume"
          containerPath = "/dev/shm"
          readOnly = false
        }
      ]

      # environment = [
      #   { name = "s3_source_bucket", value = var.s3_source_name},
      #   { name = "file_store", value = var.s3_output_store_name},
      #   { name = "file_id", value = "cookies-grouped-remi"},
       
      # ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.scroller_task.name
          awslogs-region        = "eu-west-2"
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}
resource "aws_cloudwatch_log_group" "scroller_task" {
  name              = "/ecs/scroller-task"
  retention_in_days = 7
}
