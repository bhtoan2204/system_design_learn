# S3 Module Outputs

output "app_bucket_id" {
  description = "ID of the application S3 bucket"
  value       = var.create_app_bucket ? aws_s3_bucket.app_bucket[0].id : null
}

output "app_bucket_arn" {
  description = "ARN of the application S3 bucket"
  value       = var.create_app_bucket ? aws_s3_bucket.app_bucket[0].arn : null
}

output "app_bucket_domain_name" {
  description = "Domain name of the application S3 bucket"
  value       = var.create_app_bucket ? aws_s3_bucket.app_bucket[0].bucket_domain_name : null
}

output "static_website_bucket_id" {
  description = "ID of the static website S3 bucket"
  value       = var.create_static_website_bucket ? aws_s3_bucket.static_website_bucket[0].id : null
}

output "static_website_bucket_arn" {
  description = "ARN of the static website S3 bucket"
  value       = var.create_static_website_bucket ? aws_s3_bucket.static_website_bucket[0].arn : null
}

output "static_website_bucket_domain_name" {
  description = "Domain name of the static website S3 bucket"
  value       = var.create_static_website_bucket ? aws_s3_bucket.static_website_bucket[0].bucket_domain_name : null
}

output "static_website_bucket_website_endpoint" {
  description = "Website endpoint of the static website S3 bucket"
  value       = var.create_static_website_bucket ? aws_s3_bucket_website_configuration.static_website_config[0].website_endpoint : null
}

output "backup_bucket_id" {
  description = "ID of the backup S3 bucket"
  value       = var.create_backup_bucket ? aws_s3_bucket.backup_bucket[0].id : null
}

output "backup_bucket_arn" {
  description = "ARN of the backup S3 bucket"
  value       = var.create_backup_bucket ? aws_s3_bucket.backup_bucket[0].arn : null
}

output "backup_bucket_domain_name" {
  description = "Domain name of the backup S3 bucket"
  value       = var.create_backup_bucket ? aws_s3_bucket.backup_bucket[0].bucket_domain_name : null
}
