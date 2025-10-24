# Development Environment Outputs

# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

output "internet_gateway_id" {
  description = "ID of the Internet Gateway"
  value       = module.vpc.internet_gateway_id
}

# EC2 Outputs
output "ec2_instance_ids" {
  description = "IDs of the EC2 instances"
  value       = module.ec2.instance_ids
}

output "ec2_instance_public_ips" {
  description = "Public IPs of the EC2 instances"
  value       = module.ec2.instance_public_ips
}

output "ec2_instance_private_ips" {
  description = "Private IPs of the EC2 instances"
  value       = module.ec2.instance_private_ips
}

output "ec2_key_pair_name" {
  description = "Name of the EC2 key pair"
  value       = module.ec2.key_pair_name
}

output "ec2_security_group_id" {
  description = "ID of the EC2 security group"
  value       = module.ec2.security_group_id
}

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = module.ec2.alb_dns_name
}

# S3 Outputs
output "app_bucket_id" {
  description = "ID of the application S3 bucket"
  value       = module.s3.app_bucket_id
}

output "app_bucket_arn" {
  description = "ARN of the application S3 bucket"
  value       = module.s3.app_bucket_arn
}

output "static_website_bucket_id" {
  description = "ID of the static website S3 bucket"
  value       = module.s3.static_website_bucket_id
}

output "static_website_bucket_website_endpoint" {
  description = "Website endpoint of the static website S3 bucket"
  value       = module.s3.static_website_bucket_website_endpoint
}

output "backup_bucket_id" {
  description = "ID of the backup S3 bucket"
  value       = module.s3.backup_bucket_id
}

# IAM Outputs
output "ec2_role_arn" {
  description = "ARN of the EC2 instance role"
  value       = module.iam.ec2_role_arn
}

output "ec2_instance_profile_name" {
  description = "Name of the EC2 instance profile"
  value       = module.iam.ec2_instance_profile_name
}

output "rds_role_arn" {
  description = "ARN of the RDS role"
  value       = module.iam.rds_role_arn
}

output "lambda_role_arn" {
  description = "ARN of the Lambda execution role"
  value       = module.iam.lambda_role_arn
}

# RDS Outputs
output "db_instance_id" {
  description = "ID of the RDS instance"
  value       = module.rds.db_instance_id
}

output "db_instance_endpoint" {
  description = "Endpoint of the RDS instance"
  value       = module.rds.db_instance_endpoint
}

output "db_instance_address" {
  description = "Address of the RDS instance"
  value       = module.rds.db_instance_address
}

output "db_instance_port" {
  description = "Port of the RDS instance"
  value       = module.rds.db_instance_port
}

output "read_replica_endpoint" {
  description = "Endpoint of the read replica"
  value       = module.rds.read_replica_endpoint
}

# Environment Summary
output "environment_summary" {
  description = "Summary of the development environment"
  value = {
    environment = "dev"
    vpc_id      = module.vpc.vpc_id
    ec2_count   = length(module.ec2.instance_ids)
    rds_enabled = var.create_rds
    s3_buckets  = {
      app_bucket     = module.s3.app_bucket_id != null
      static_website = module.s3.static_website_bucket_id != null
      backup         = module.s3.backup_bucket_id != null
    }
  }
}
