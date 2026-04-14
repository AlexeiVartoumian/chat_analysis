# FilePool Table
resource "aws_dynamodb_table" "filepool" {
  name         = "filepoolstore"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "file_id"

  attribute {
    name = "file_id"
    type = "S"
  }

  attribute {
    name = "status"
    type = "S"
  }

  attribute {
    name = "locked_at"
    type = "N"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  global_secondary_index {
    name            = "status-index"
    hash_key        = "status"
    range_key       = "locked_at"
    projection_type = "ALL"
  }

  tags = {
    Name = "filepoolstore"
  }
}

# AccountPool Table
resource "aws_dynamodb_table" "accountpool" {
  name         = "accountpoolstore"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "account_id"

  attribute {
    name = "account_id"
    type = "S"
  }

  tags = {
    Name = "accountpoolstore"
  }
}
