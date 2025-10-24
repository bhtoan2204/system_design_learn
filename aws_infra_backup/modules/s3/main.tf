# S3 Module
# This module creates S3 buckets with appropriate configurations

terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# Application S3 Bucket
resource "aws_s3_bucket" "app_bucket" {
  count = var.create_app_bucket ? 1 : 0
  
  bucket = "${var.environment}-${var.app_name}-app-${random_string.bucket_suffix[0].result}"
  
  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.app_name}-app"
    Environment = var.environment
    Purpose     = "Application Data"
  })
}

# Static Website S3 Bucket
resource "aws_s3_bucket" "static_website_bucket" {
  count = var.create_static_website_bucket ? 1 : 0
  
  bucket = "${var.environment}-${var.app_name}-static-${random_string.bucket_suffix[0].result}"
  
  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.app_name}-static"
    Environment = var.environment
    Purpose     = "Static Website"
  })
}

# Backup S3 Bucket
resource "aws_s3_bucket" "backup_bucket" {
  count = var.create_backup_bucket ? 1 : 0
  
  bucket = "${var.environment}-${var.app_name}-backup-${random_string.bucket_suffix[0].result}"
  
  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.app_name}-backup"
    Environment = var.environment
    Purpose     = "Backup Storage"
  })
}

# Random string for bucket suffix
resource "random_string" "bucket_suffix" {
  count = var.create_app_bucket || var.create_static_website_bucket || var.create_backup_bucket ? 1 : 0
  
  length  = 8
  special = false
  upper   = false
}

# S3 Bucket Versioning
resource "aws_s3_bucket_versioning" "app_bucket_versioning" {
  count = var.create_app_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.app_bucket[0].id
  versioning_configuration {
    status = var.enable_versioning ? "Enabled" : "Suspended"
  }
}

resource "aws_s3_bucket_versioning" "backup_bucket_versioning" {
  count = var.create_backup_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.backup_bucket[0].id
  versioning_configuration {
    status = "Enabled"
  }
}

# S3 Bucket Server Side Encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "app_bucket_encryption" {
  count = var.create_app_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.app_bucket[0].id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "backup_bucket_encryption" {
  count = var.create_backup_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.backup_bucket[0].id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

# S3 Bucket Public Access Block
resource "aws_s3_bucket_public_access_block" "app_bucket_pab" {
  count = var.create_app_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.app_bucket[0].id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_public_access_block" "backup_bucket_pab" {
  count = var.create_backup_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.backup_bucket[0].id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 Bucket Lifecycle Configuration
resource "aws_s3_bucket_lifecycle_configuration" "app_bucket_lifecycle" {
  count = var.create_app_bucket && var.enable_lifecycle ? 1 : 0
  
  bucket = aws_s3_bucket.app_bucket[0].id
  
  rule {
    id     = "app_bucket_lifecycle"
    status = "Enabled"
    
    expiration {
      days = var.lifecycle_expiration_days
    }
    
    noncurrent_version_expiration {
      noncurrent_days = var.lifecycle_noncurrent_expiration_days
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "backup_bucket_lifecycle" {
  count = var.create_backup_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.backup_bucket[0].id
  
  rule {
    id     = "backup_bucket_lifecycle"
    status = "Enabled"
    
    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }
    
    transition {
      days          = 90
      storage_class = "GLACIER"
    }
    
    transition {
      days          = 365
      storage_class = "DEEP_ARCHIVE"
    }
  }
}

# S3 Bucket CORS Configuration for Static Website
resource "aws_s3_bucket_cors_configuration" "static_website_cors" {
  count = var.create_static_website_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.static_website_bucket[0].id
  
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

# S3 Bucket Website Configuration
resource "aws_s3_bucket_website_configuration" "static_website_config" {
  count = var.create_static_website_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.static_website_bucket[0].id
  
  index_document {
    suffix = "index.html"
  }
  
  error_document {
    key = "error.html"
  }
}

# S3 Bucket Policy for Static Website
resource "aws_s3_bucket_policy" "static_website_policy" {
  count = var.create_static_website_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.static_website_bucket[0].id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.static_website_bucket[0].arn}/*"
      }
    ]
  })
  
  depends_on = [aws_s3_bucket_public_access_block.static_website_pab]
}

# Public Access Block for Static Website (less restrictive)
resource "aws_s3_bucket_public_access_block" "static_website_pab" {
  count = var.create_static_website_bucket ? 1 : 0
  
  bucket = aws_s3_bucket.static_website_bucket[0].id
  
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
