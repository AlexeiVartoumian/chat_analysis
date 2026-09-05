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
    sqs_cordinator_deed_arn = module.sqs_hub.coordinator_deed_sqs_queue_arn
    sqs_coordinator_work_arn = module.sqs_hub.coordinator_work_sqs_queue_arn
    aws_iam_role_main_arn = module.iam_hub.aws_iam_role_main_arn
    account_pool_table =  module.dynamodb_hub.accountpool_table_name
    file_pool_table =  module.dynamodb_hub.filepool_table_name
    file_pool_table_deed = module.dynamodb_hub.filepooldeed_table_name
    account_pool_table_deed = module.dynamodb_hub.accountpooldeed_table_name
    coordinator_work_sqs_queue_id = module.sqs_hub.coordinator_work_sqs_queue_id  
    s3_output_bucket_work_store_name = module.s3.s3_bucket_output_work_store_name
    account_pool_table_work = module.dynamodb_hub.accountpoolwork_table_name
    account_pool_table_work_arn = module.dynamodb_hub.accountpoolwork_table_arn
    account_pool_table_stream_arn = module.dynamodb_hub.accountpoolwork_stream_arn
    eventbridge_rule_arn = module.eventbridge_hub.file_created_rule_arn

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
    file_pool_table_ash = module.dynamodb_hub.filepoolash_table_name
    file_pool_table_green = module.dynamodb_hub.filepoolgreen_table_name  
    file_pool_table_deed = module.dynamodb_hub.filepooldeed_table_name

    account_pool_table = module.dynamodb_hub.accountpool_table_name 
    bucket_reader_role_name= module.iam_hub.aws_iam_role_main_name
    bucket_reader_main_arn = module.iam_hub.aws_iam_role_main_arn
    sqs_coordinator_arn = module.sqs_hub.coordinator_sqs_queue_arn
    sqs_deadletter_arn = module.sqs_hub.deadletter_sqs_queue_arn
    sqs_coordinator_arn_deed = module.sqs_hub.coordinator_deed_sqs_queue_arn
    sqs_coordinator_arn_work = module.sqs_hub.coordinator_work_sqs_queue_arn
    
    account_pool_table_deed = module.dynamodb_hub.accountpooldeed_table_name
    account_pool_table_work = module.dynamodb_hub.accountpoolwork_table_name
    account_pool_table_stream_arn = module.dynamodb_hub.accountpoolwork_stream_arn

    
    s3_output_bucket_ash_cache_arn = module.s3.s3_bucket_output_ash_cache_arn
    s3_output_bucket_ash_store_arn = module.s3.s3_bucket_output_ash_store_arn

    s3_output_bucket_deed_cache_arn = module.s3.s3_bucket_output_deed_cache_arn
    s3_output_bucket_deed_store_arn = module.s3.s3_bucket_output_deed_store_arn

    s3_output_bucket_green_cache_arn = module.s3.s3_bucket_output_green_cache_arn
    s3_output_bucket_green_store_arn = module.s3.s3_bucket_output_green_store_arn

    s3_output_bucket_work_cache_arn = module.s3.s3_bucket_output_work_cache_arn
    s3_output_bucket_work_store_arn = module.s3.s3_bucket_output_work_store_arn
    providers = {
        aws = aws.hub
    }
    depends_on = [module.iam_hub]
}

module "fargate_hub"{
    source = "./module/fargate_hub"
    iam_role_main_arn = module.iam_hub.aws_iam_role_main_arn
    spoke_accounts = var.spoke_accounts
    # s3_source_name = module.s3.s3_bucket_source_name
    # s3_filestore_name = module.s3.s3_bucket_file_name

     providers = {
        aws = aws.hub
    }
}

module "ec2_hub"{
    source = "./module/ec2_hub"
    scroller_profile = module.iam_hub.aws_iam_role_scroller_profile
    ami_id = var.ami_id
    keys = var.keys
}

module "eventbridge_hub"{
    source = "./module/eventbridge_hub"
    #bucket_work_id = module.s3.s3_bucket_output_work_store_id
    lambda_work_arn = module.lambda_hub.lambda_work_arn

}