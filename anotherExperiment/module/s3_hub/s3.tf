resource "aws_s3_bucket" "source_store" {
    bucket = "${local.s3_bucket_name_source}-${data.aws_caller_identity.current.account_id}"
}


resource "aws_s3_bucket" "file_store" {

    bucket = "${local.s3_bucket_name_file}-${data.aws_caller_identity.current.account_id}"
}

resource "aws_s3_bucket" "output_store" {

    bucket = "${local.s3_bucket_name_output}-${data.aws_caller_identity.current.account_id}"

}
