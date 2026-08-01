output "filepool_table_name" {
  value       = aws_dynamodb_table.filepool.name
}

output "accountpool_table_name" {
  value       = aws_dynamodb_table.accountpool.name
  
}



output "filepooldeed_table_name" {
  value       = aws_dynamodb_table.filepool_deed.name
}

output "accountpooldeed_table_name" {
  value       = aws_dynamodb_table.accountpool_deed.name
  
}


output "filepoolash_table_name" {
  value       = aws_dynamodb_table.filepool_ash.name
}

output "filepoolgreen_table_name" {
  value       = aws_dynamodb_table.filepool_green.name
}

