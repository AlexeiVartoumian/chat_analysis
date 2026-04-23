module "sqs_hub"{
    source = "./module/sqs_hub"

    providers = {
        aws = aws.hub
    }
   
}
module dynamodb_hub {
    source = "./module/dynamodb_hub"
    providers = {
      aws = aws.hub
    }
}

module "iam_hub"{
    source = "./module/iam_hub"
    spoke_accounts = var.spoke_accounts
}

module "s3" {
    source = "./module/s3_hub"
    s3_bucket_name_file = "file-store"
    s3_bucket_name_output = "output-store"
    s3_bucket_name_source = "source-store"

    providers = {
      aws = aws.hub
    }
}

module "lambda_hub"{
    source = "./module/lambda_hub"
    sqs_coordinator_arn = module.sqs_hub.coordinator_sqs_queue_arn   
    aws_iam_role_main_arn = module.iam_hub.aws_iam_role_main_arn
    account_pool_table =  module.dynamodb_hub.accountpool_table_name
    file_pool_table =  module.dynamodb_hub.filepool_table_name
     providers = {
        aws = aws.hub
    }
}
module "iam_hub_attachments" {
    source = "./module/iam_hub_attachments"
    spoke_accounts = var.spoke_accounts
    hub_account = var.hub_account
    s3_source_bucket_arn = module.s3.s3_bucket_source_arn
    s3_file_bucket_arn = module.s3.s3_bucket_file_arn
    s3_output_bucket_arn = module.s3.s3_bucket_output_arn
    s3_backfill_bucket_arn = module.s3.s3_bucket_backfill_arn
    file_pool_table = module.dynamodb_hub.filepool_table_name
    account_pool_table = module.dynamodb_hub.accountpool_table_name 
    bucket_reader_role_name= module.iam_hub.aws_iam_role_main_name
    bucket_reader_main_arn = module.iam_hub.aws_iam_role_main_arn
    sqs_coordinator_arn = module.sqs_hub.coordinator_sqs_queue_arn
    providers = {
        aws = aws.hub
    }
    depends_on = [module.iam_hub]
}

module "fargate_hub"{
    source = "./module/fargate_hub"
    iam_role_main_arn = module.iam_hub.aws_iam_role_main_arn
    # s3_source_name = module.s3.s3_bucket_source_name
    # s3_filestore_name = module.s3.s3_bucket_file_name

     providers = {
        aws = aws.hub
    }
}