variable "s3_bucket_name_source" {
    type = string
    description = "name of bucket "

    default = "source-store-2026-3s-intricate"
}


variable "s3_bucket_name_file" {
    type = string
    description = "name of bucket "

    default = "file-store-2026-3s-intricate"
}

variable "s3_bucket_name_output" {
    type = string
    description = "name of bucket "

    default = "output-store"
}

variable "s3_bucket_name_backfill" {
    type = string
    description = "name of bucket"

    default = "backfill-store"
}

variable "s3_bucket_name_output_ash_cache" {
    type = string
    description = "name of bucket"
    default = "output-store-ash-cache"
}

variable "s3_bucket_name_output_ash_store" {
    type = string
    description = "name of bucket"
    default = "output-store-ash-store"
}

variable "s3_bucket_name_output_deed_cache" {
    type = string
    description = "name of bucket"
    default = "output-store-deed-cache"
}

variable "s3_bucket_name_output_deed_store" {
    type = string
    description = "name of bucket"
    default = "output-store-deed-store"
}

