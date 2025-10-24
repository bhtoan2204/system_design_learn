# Backend Infrastructure Outputs

output "dev_s3_bucket" {
  description = "S3 bucket for dev environment state"
  value       = aws_s3_bucket.dev_state.bucket
}

output "staging_s3_bucket" {
  description = "S3 bucket for staging environment state"
  value       = aws_s3_bucket.staging_state.bucket
}

output "prod_s3_bucket" {
  description = "S3 bucket for prod environment state"
  value       = aws_s3_bucket.prod_state.bucket
}

output "dev_dynamodb_table" {
  description = "DynamoDB table for dev environment state locking"
  value       = aws_dynamodb_table.dev_locks.name
}

output "staging_dynamodb_table" {
  description = "DynamoDB table for staging environment state locking"
  value       = aws_dynamodb_table.staging_locks.name
}

output "prod_dynamodb_table" {
  description = "DynamoDB table for prod environment state locking"
  value       = aws_dynamodb_table.prod_locks.name
}

output "summary" {
  description = "Summary of created backend resources"
  value = {
    s3_buckets = {
      dev     = aws_s3_bucket.dev_state.bucket
      staging = aws_s3_bucket.staging_state.bucket
      prod    = aws_s3_bucket.prod_state.bucket
    }
    dynamodb_tables = {
      dev     = aws_dynamodb_table.dev_locks.name
      staging = aws_dynamodb_table.staging_locks.name
      prod    = aws_dynamodb_table.prod_locks.name
    }
  }
}
