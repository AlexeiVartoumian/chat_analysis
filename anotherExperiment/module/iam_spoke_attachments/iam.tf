




resource "aws_iam_role_policy" "sqs" {

    role = var.bucket_reader_spoke
    policy = templatefile("${path.module}/spoke_sqs_permissions.tpl" , {
         aws_account  = data.aws_caller_identity.current.account_id
         sqs_queue_request = var.sqs_request_access
         sqs_queue_2 = var.sqs_queue_2
         sqs_queue_3 = var.sqs_queue_3
  })
}
resource "aws_iam_role_policy" "assume_hub_role" {

    role = var.bucket_reader_spoke
    policy = templatefile("${path.module}/spoke_iam_assume.tpl" , {
        hub_account = var.hub_account
  })
}

resource "aws_iam_role_policy_attachment" "lambda_exec" {
  role = var.bucket_reader_spoke
  policy_arn = data.aws_iam_policy.lambda_basic_execution.arn
}


resource "aws_iam_role_policy" "fargate_spoke" {
  role = var.fargate_spoke

  policy = templatefile("${path.module}/spoke_fargate_execution.tpl" , {
        hub_account = var.hub_account
  })
}

resource "aws_iam_role_policy" "assume_fargatehub_role" {

    role = var.fargate_spoke
    policy = templatefile("${path.module}/spoke_iam_assume.tpl" , {
        hub_account = var.hub_account
  })
}


resource "aws_iam_policy" "spoke_invoker_policy" { 
    name = "ecs_task_invoker_policy"
    policy = templatefile("${path.module}/fargate_invoker.tpl" ,{
      spoke_account_id = data.aws_caller_identity.current.account_id
      bucket_reader_spoke = var.bucket_reader_spoke
    })
}

resource "aws_iam_role_policy_attachment" "spoke_invoker_attach" {
  
  role = var.fargate_invoker_name
  policy_arn = aws_iam_policy.spoke_invoker_policy.arn
}