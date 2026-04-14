
# module dynamodb_hub {
#     source = "./module/dynamodb_hub"

#     providers = {
#       aws = aws.hub
#     }
# }

# module "iam_hub"{
#     source = "./module/iam_hub"
#     spoke_accounts = var.spoke_accounts
# }

# module "iam_spoke" {
#     source = "./module/iam_spoke"

#        providers = {
#       aws = aws.spoke
#     }
# }

# module "iam_hub_attachments" {
#     source = "./module/iam_hub_attachments"
#     spoke_accounts = var.spoke_accounts
#     hub_account = var.hub_account
#     s3_source_bucket_arn = module.s3.s3_bucket_source_arn
#     s3_file_bucket_arn = module.s3.s3_bucket_file_arn
#     s3_output_bucket_arn = module.s3.s3_bucket_output_arn
#     file_pool_table = module.dynamodb_hub.filepool_table_name
#     account_pool_table = module.dynamodb_hub.accountpool_table_name 
#     bucket_reader_role_name= module.iam_hub.aws_iam_role_main_name
#     bucket_reader_main_arn = module.iam_hub.aws_iam_role_main_arn

#     sqs_coordinator_arn = module.sqs_hub.coordinator_sqs_queue_arn
#     providers = {
#         aws = aws.hub
#     }
#     depends_on = [module.iam_hub]
# }

# module "iam_spoke_attachments" {
#     source = "./module/iam_spoke_attachments"
#     sqs_request_access = module.sqs_spoke.workflow_requests_sqs_queue_name
#     sqs_queue_2 = module.sqs_spoke.workflow_lambda2_trigger_sqs_queue_name
#     sqs_queue_3 = module.sqs_spoke.workflow_lambda3_trigger_sqs_queue_name
#     bucket_reader_spoke = module.iam_spoke.iam_spoke_role_name
#     hub_account = var.hub_account

#     providers = {
#       aws = aws.spoke
#     }
# }


# module "s3" {
#     source = "./module/s3_hub"
#     s3_bucket_name_file = "file-store"
#     s3_bucket_name_output = "output-store"
#     s3_bucket_name_source = "source-store"

#     providers = {
#       aws = aws.hub
#     }
# }
# module "sqs_spoke"{
#     source = "./module/sqs"
#     hub_account = var.hub_account
#     aws_iam_role_main_name = module.iam_hub.aws_iam_role_main_name
#     providers = {
#       aws = aws.spoke
#     }
# }

# module "sqs_hub"{
#     source = "./module/sqs_hub"

#     providers = {
#         aws = aws.hub
#     }
   
# }

# module "lambda_hub"{
#     source = "./module/lambda_hub"
#     sqs_coordinator_arn = module.sqs_hub.coordinator_sqs_queue_arn   
#     aws_iam_role_main_arn = module.iam_hub.aws_iam_role_main_arn
#     account_pool_table =  module.dynamodb_hub.accountpool_table_name
#     file_pool_table =  module.dynamodb_hub.filepool_table_name
#      providers = {
#         aws = aws.hub
#     }
# }

# module "lambda_spoke" {
#     source = "./module/lambda_spoke"
#     iam_role_arn_spoke = module.iam_spoke.iam_spoke_role_arn
#     iam_role_main_arn = module.iam_hub.aws_iam_role_main_arn
#     s3_source_name = module.s3.s3_bucket_source_name
#     s3_filestore_name = module.s3.s3_bucket_file_name
#     s3_output_bucket_name = module.s3.s3_bucket_output_name
#     sqs_request_access_arn = module.sqs_spoke.workflow_requests_sqs_queue_arn
#     sqs_queue_2_arn = module.sqs_spoke.workflow_lambda2_trigger_sqs_queue_arn
#     sqs_queue_3_arn = module.sqs_spoke.workflow_lambda3_trigger_sqs_queue_arn
#     sqs_queue_2_id = module.sqs_spoke.workflow_lambda2_trigger_sqs_queue_id
#     sqs_queue_3_id = module.sqs_spoke.workflow_lambda3_trigger_sqs_queue_id
#     dynamodb_filetable_name = module.dynamodb_hub.filepool_table_name
#     dynamodb_accounttable_name = module.dynamodb_hub.accountpool_table_name 

#     providers = {
#       aws = aws.spoke
#     }
# }