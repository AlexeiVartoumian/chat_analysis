

resource aws_iam_role "bucket_reader_spoke" {

    name = "the_bucket_dealer_spoke"
    assume_role_policy = templatefile("${path.module}/assume.tpl", { none = "none" })

}