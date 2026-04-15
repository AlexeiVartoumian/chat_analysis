resource "aws_ecs_cluster" "reader_cluster" {
  name = "reader-cluster"
}

resource "aws_iam_role" "ecs_task_role" {
  name = "ecs_reader_task_role"

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

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_task_definition" "reader_task" {
  family                   = "reader-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  task_role_arn            = var.iam_role_main_arn
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "reader-container"
      image     = "390746273208.dkr.ecr.eu-west-2.amazonaws.com/reader:latest"
      essential = true
      command   = [
     
      ]
      environment = [
        { name = "s3_source_bucket", value = var.s3_source_name},
        { name = "file_store", value = var.s3_filestore_name},
        { name = "file_id", value = "cookies-grouped-remi"},
        { name = "search_term", value = "cloud engineer"}
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.reader_task.name
          awslogs-region        = "eu-west-2"
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}


resource "aws_cloudwatch_log_group" "reader_task" {
  name              = "/ecs/reader-task"
  retention_in_days = 7
}
