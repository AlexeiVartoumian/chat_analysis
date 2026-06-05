
output "iam_spoke_role_arn" {
    value = aws_iam_role.bucket_reader_spoke.arn
}

output "iam_spoke_role_name" {
    value = aws_iam_role.bucket_reader_spoke.name
}


output "iam_fargate_task_execution_role_name" {
    value = aws_iam_role.ecs_task_execution_role.name
    description = "this role is responsible for getting the container up and running . different from the role that gives the permissions needed once container is running"
}

output "iam_fargate_task_execution_spoke_role_arn" {
    value = aws_iam_role.ecs_task_execution_role.arn
    description = "get the container up and running"
 }

 output "iam_fargate_invoker_role_name" {
    value = aws_iam_role.fargate_spoke_invoker_role.name
    description = "this role is responsible for getting the container up and running . different from the role that gives the permissions needed once container is running"
}

output "iam_fargate_invoker_role_arn" {
    value = aws_iam_role.fargate_spoke_invoker_role.arn
    description = "get the container up and running"
 }