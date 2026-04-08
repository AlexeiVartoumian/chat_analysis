

resource aws_iam_role "bucket_reader_main" {

    name = "the_bucket_dealer"
    assume_role_policy = templatefile("${path.module}/assume.tpl", 
    { spoke_accounts = var.spoke_accounts 
    })

}
resource "aws_iam_role_policy" "send_spoke_sqs" {
  role = aws_iam_role.bucket_reader_main.name
  policy = templatefile("${path.module}/send_spoke_sqs.tpl" , {
    spoke_accounts = var.spoke_accounts
  })
}


resource "aws_iam_role_policy" "bucket_permissions" {
  role = aws_iam_role.bucket_reader_main.name
  policy = templatefile("${path.module}/bucket_permissions_for_sqs.tpl" , {
    s3_cookie_bucket_arn = var.s3_cookie_bucket
    s3_file_bucket = var.s3_file_bucket
    s3_output_bucket = var.s3_output_bucket
  })
}

resource "aws_iam_role_policy" "dynamodb_permissions" {
  role = aws_iam_role.bucket_reader_main.name
  policy = templatefile("${path.module}/dynamodb.tpl" , {
    hub_account = var.hub_account
  })
}

resource "aws_iam_role_policy" "sqs_permissions" {
  role = aws_iam_role.sqs.name
  policy = templatefile("${path.module}/sqs_permissions.tpl" , {
    hub_account = var.hub_account
  })
}
