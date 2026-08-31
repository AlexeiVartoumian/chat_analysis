

resource "aws_s3_bucket_notification" "workstore" {
    bucket = var.bucket_work_id
    eventbridge = true
}



module "eventbridge" {
    source  = "terraform-aws-modules/eventbridge/aws"
    version = "~> 3.0"

    create_bus = false #use default bus 

    rules = {
        csv-created = {

            description = "s3 object created"

            event_pattern = jsonencode({
                source = ["aws.s3"]
                detail-type = ["Object Created"]
                detail = {
                    bucket = {name = [var.bucket_work_id]}
                    object = {key = [{suffix = ".csv"}]}
                }
            })
        }
    }

    targets = {
        csv-created = [
            {
                name = "routing-work-lambda"
                arn = var.lambda_work_arn

            }
        ]
    }

}