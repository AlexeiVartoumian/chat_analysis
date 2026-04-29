

data "archive_file" "orchestrator_path" {
    type = "zip"
    source_file = "${path.root}/module/sources/orchestrator/orchestrator.py"
    output_path = "${path.root}/module/sources/orchestrator/orchestrator.zip"
}

resource "aws_lambda_function" "orchestrator" {
    filename = data.archive_file.orchestrator_path.output_path
    source_code_hash = data.archive_file.orchestrator_path.output_base64sha256

    function_name = "orchestratortest"
    role = var.aws_iam_role_main_arn
    handler = "orchestrator.lambda_handler"
    runtime = "python3.13" 
    environment {
        variables = {
            account_pool_table= var.account_pool_table
            file_pool_table = var.file_pool_table
        }
    }

    #depends_on = [aws_cloudwatch_log_group.orchestrator]
}

resource "aws_lambda_permission" "allow_sqs_request" {
  statement_id  = "AllowExecutionFromSqs"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.orchestrator.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = var.sqs_coordinator_arn
}

resource "aws_lambda_event_source_mapping" "processor_trigger" {
  event_source_arn = var.sqs_coordinator_arn
  function_name    = aws_lambda_function.orchestrator.arn
  batch_size       = 10
  enabled          = true
  depends_on = [aws_lambda_function.orchestrator]
}

# resource "aws_cloudwatch_log_group" "orchestrator" {
#     name = "/aws/lambda/orchestratortest"
#     retention_in_days = 7
# }