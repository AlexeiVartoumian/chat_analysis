resource "aws_s3_bucket" "source_store" {
    bucket = "${local.s3_bucket_name_source}-${data.aws_caller_identity.current.account_id}"
}


resource "aws_s3_bucket" "file_store" {

    bucket = "${local.s3_bucket_name_file}-${data.aws_caller_identity.current.account_id}"
}

resource "aws_s3_bucket" "output_store" {

    bucket = "${local.s3_bucket_name_output}-${data.aws_caller_identity.current.account_id}"

}

resource "aws_s3_bucket" "backfill_store" {

    bucket = "${local.s3_bucket_name_backfill}-${data.aws_caller_identity.current.account_id}"

}


resource "aws_s3_bucket" "output_store_ash_cache" {

    bucket = "${local.s3_bucket_name_output_ash_cache}-${data.aws_caller_identity.current.account_id}"

}

resource "aws_s3_bucket" "output_store_ash_store" {

    bucket = "${local.s3_bucket_name_output_ash_store }-${data.aws_caller_identity.current.account_id}"

}


resource "aws_s3_bucket" "output_store_deed_cache" {

    bucket = "${local.s3_bucket_name_output_deed_cache}-${data.aws_caller_identity.current.account_id}"

}

resource "aws_s3_bucket" "output_store_deed_store" {

    bucket = "${local.s3_bucket_name_output_deed_store }-${data.aws_caller_identity.current.account_id}"

}


resource "aws_s3_bucket" "output_store_green_cache" {

    bucket = "${local.s3_bucket_name_output_green_cache}-${data.aws_caller_identity.current.account_id}"

}

resource "aws_s3_bucket" "output_store_green_store" {

    bucket = "${local.s3_bucket_name_output_green_store }-${data.aws_caller_identity.current.account_id}"

}