data "archive_file" "reader_path" {
    type = "zip"
    source_file = "${path.root}/module/sources/package/reader.py"
    output_path = "${path.root}/module/sources/package/reader.zip"
}


data "archive_file" "processor_path" {
    type = "zip"
    source_file = "${path.root}/module/sources/file/processfile.py"
    output_path = "${path.root}/module/sources/file/processfile.zip"
}


data "archive_file" "go_metadata_path" {
    type = "zip"
    source_file = "${path.root}/module/sources/booter/bootstrap"
    output_path = "${path.root}/module/sources/booter/bootstrap.zip"
}

data "archive_file" "requests_layer" {
  type        = "zip"
  source_dir  = "${path.root}/module/sources/layer"
  output_path = "${path.root}/module/sources/layer/python.zip"
}

resource "aws_lambda_function" "reader" {
    filename = data.archive_file.reader_path.output_path
    source_code_hash = data.archive_file.reader_path.output_base64sha256
    function_name = "readerv2"
    role = var.iam_role_arn_spoke 
    handler = "reader.lambda_handler"
    runtime = "python3.13" 
    timeout     = 900
    layers = [aws_lambda_layer_version.requests_layer.arn]

    environment {
        variables = {
            RoleArn = var.iam_role_main_arn
            s3_source_bucket = var.s3_source_name
            file_store = var.s3_filestore_name
            sqs_2_id = var.sqs_queue_2_id
        }
    }

    depends_on = [ aws_cloudwatch_log_group.reader ]
}

resource "aws_lambda_function" "processor" {
    filename = data.archive_file.processor_path.output_path
     source_code_hash = data.archive_file.processor_path.output_base64sha256
    function_name = "processorv2"
    role = var.iam_role_arn_spoke
    handler = "processfile.lambda_handler"
    runtime = "python3.13" 
    timeout     = 900
    layers = [aws_lambda_layer_version.requests_layer.arn]
    environment {
        variables = {
            RoleArn = var.iam_role_main_arn
            output_bucket= var.s3_output_bucket_name
            sqs_queue_3_id = var.sqs_queue_3_id
        }
    }
    depends_on = [aws_cloudwatch_log_group.processor]
}

resource "aws_lambda_function" "go_metadata" {
    filename = data.archive_file.go_metadata_path.output_path
    //source_code_hash = data.archive_file.zip.output_base64sha256
    source_code_hash = data.archive_file.go_metadata_path.output_base64sha256
    function_name = "go_metadatav2"
    role =  var.iam_role_arn_spoke
    handler = "bootstrap"
    runtime = "provided.al2023"
    timeout     = 900
    environment {
        variables = {
            RoleArn = var.iam_role_main_arn
            output_bucket= var.s3_output_bucket_name
            s3_source_bucket = var.s3_source_name
            file_pool = var.dynamodb_filetable_name
            account_pool = var.dynamodb_accounttable_name
        }
    }
    depends_on = [ aws_cloudwatch_log_group.go_metadata ]
}

resource "aws_lambda_permission" "allow_sqs_request" {
  statement_id  = "AllowExecutionFromSqs"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.reader.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = var.sqs_request_access_arn
}


resource "aws_lambda_permission" "allow_sqs2" {
  statement_id  = "AllowExecutionFromSqs"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.processor.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = var.sqs_queue_2_arn
}


resource "aws_lambda_permission" "allow_sqs3" {
  statement_id  = "AllowExecutionFromSqs"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.go_metadata.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = var.sqs_queue_3_arn
}

resource "aws_lambda_event_source_mapping" "reader_trigger" {
  event_source_arn = var.sqs_request_access_arn
  function_name    = aws_lambda_function.reader.arn
  batch_size       = 10
  enabled          = true
  depends_on = [aws_lambda_function.reader]
}

resource "aws_lambda_event_source_mapping" "processor_trigger" {
  event_source_arn = var.sqs_queue_2_arn
  function_name    = aws_lambda_function.processor.arn
  batch_size       = 10
  enabled          = true
depends_on = [aws_lambda_function.processor]
}

resource "aws_lambda_event_source_mapping" "go_metadata_trigger" {
  event_source_arn = var.sqs_queue_3_arn
  function_name    = aws_lambda_function.go_metadata.arn
  batch_size       = 10
  enabled          = true
  depends_on = [aws_lambda_function.go_metadata]
}


resource "aws_cloudwatch_log_group" "reader" {
    name = "/aws/lambda/readerv2"
    retention_in_days = 7

}


resource "aws_cloudwatch_log_group" "processor" {
    name = "/aws/lambda/processorv2"
    retention_in_days = 7
}
resource "aws_cloudwatch_log_group" "go_metadata" {
    name = "/aws/lambda/go_metadatav2"
    retention_in_days = 7
}