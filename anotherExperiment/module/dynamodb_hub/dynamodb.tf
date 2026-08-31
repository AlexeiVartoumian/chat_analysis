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


resource "aws_dynamodb_table" "filepool_deed" {
  name         = "filepoolstore_deed"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "profile_id"

  attribute {
    name = "profile_id"
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
    Name = "filepoolstore_deed"
  }
}

# AccountPool Table
resource "aws_dynamodb_table" "accountpool_deed" {
  name         = "accountpoolstore_deed"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "account_id"

  attribute {
    name = "account_id"
    type = "S"
  }

  tags = {
    Name = "accountpoolstore_deed"
  }
}



resource "aws_dynamodb_table" "filepool_ash" {
  name         = "filepoolstore_ash"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "profile_id"

  attribute {
    name = "profile_id"
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
    Name = "filepoolstore_ash"
  }
}


resource "aws_dynamodb_table" "filepool_green" {
  name         = "filepoolstore_green"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "profile_id"

  attribute {
    name = "profile_id"
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
    Name = "filepoolstore_green"
  }
}



resource "aws_dynamodb_table" "accountpoolwork" {
  name         = "accountpoolwork"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "account"

  stream_enabled = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  attribute {
    name = "account"
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
  global_secondary_index {
    name            = "status-index"
    hash_key        = "status"
    range_key       = "locked_at"
    projection_type = "ALL"
  }

  tags = {
    Name = "accountpoolwork"
  }
}
