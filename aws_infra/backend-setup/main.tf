# Backend Infrastructure Setup
# This Terraform configuration creates S3 buckets and DynamoDB tables for remote state management

terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project     = "terraform-backend"
      ManagedBy   = "Terraform"
      Purpose     = "Remote State Management"
    }
  }
}

# S3 bucket for dev environment state
resource "aws_s3_bucket" "dev_state" {
  bucket = "terraform-state-dev-demo"
  
  tags = {
    Name        = "terraform-state-dev-demo"
    Environment = "dev"
    Purpose     = "Terraform State Storage"
  }
}

resource "aws_s3_bucket_versioning" "dev_state_versioning" {
  bucket = aws_s3_bucket.dev_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "dev_state_encryption" {
  bucket = aws_s3_bucket.dev_state.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_public_access_block" "dev_state_pab" {
  bucket = aws_s3_bucket.dev_state.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 bucket for staging environment state
resource "aws_s3_bucket" "staging_state" {
  bucket = "terraform-state-staging-demo"
  
  tags = {
    Name        = "terraform-state-staging-demo"
    Environment = "staging"
    Purpose     = "Terraform State Storage"
  }
}

resource "aws_s3_bucket_versioning" "staging_state_versioning" {
  bucket = aws_s3_bucket.staging_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "staging_state_encryption" {
  bucket = aws_s3_bucket.staging_state.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_public_access_block" "staging_state_pab" {
  bucket = aws_s3_bucket.staging_state.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 bucket for prod environment state
resource "aws_s3_bucket" "prod_state" {
  bucket = "terraform-state-prod-demo"
  
  tags = {
    Name        = "terraform-state-prod-demo"
    Environment = "prod"
    Purpose     = "Terraform State Storage"
  }
}

resource "aws_s3_bucket_versioning" "prod_state_versioning" {
  bucket = aws_s3_bucket.prod_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "prod_state_encryption" {
  bucket = aws_s3_bucket.prod_state.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_public_access_block" "prod_state_pab" {
  bucket = aws_s3_bucket.prod_state.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# DynamoDB table for dev environment state locking
resource "aws_dynamodb_table" "dev_locks" {
  name           = "terraform-state-locks-dev"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "LockID"
  
  attribute {
    name = "LockID"
    type = "S"
  }
  
  tags = {
    Name        = "terraform-state-locks-dev"
    Environment = "dev"
    Purpose     = "Terraform State Locking"
  }
}

# DynamoDB table for staging environment state locking
resource "aws_dynamodb_table" "staging_locks" {
  name           = "terraform-state-locks-staging"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "LockID"
  
  attribute {
    name = "LockID"
    type = "S"
  }
  
  tags = {
    Name        = "terraform-state-locks-staging"
    Environment = "staging"
    Purpose     = "Terraform State Locking"
  }
}

# DynamoDB table for prod environment state locking
resource "aws_dynamodb_table" "prod_locks" {
  name           = "terraform-state-locks-prod"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "LockID"
  
  attribute {
    name = "LockID"
    type = "S"
  }
  
  tags = {
    Name        = "terraform-state-locks-prod"
    Environment = "prod"
    Purpose     = "Terraform State Locking"
  }
}
