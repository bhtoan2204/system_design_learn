# Staging Environment
# This file defines the staging environment infrastructure

terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# Provider configuration
provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Environment = "staging"
      Project     = var.project_name
      ManagedBy   = "Terraform"
      Owner       = var.owner
    }
  }
}

# Common tags
locals {
  common_tags = {
    Environment = "staging"
    Project     = var.project_name
    ManagedBy   = "Terraform"
    Owner       = var.owner
  }
}

# VPC Module
module "vpc" {
  source = "../../modules/vpc"
  
  environment           = "staging"
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  availability_zones   = var.availability_zones
  enable_nat_gateway   = var.enable_nat_gateway
  
  tags = local.common_tags
}

# IAM Module
module "iam" {
  source = "../../modules/iam"
  
  environment              = "dev"
  create_ec2_role         = var.create_ec2_role
  create_rds_role         = var.create_rds_role
  create_lambda_role      = var.create_lambda_role
  enable_s3_access        = var.enable_s3_access
  enable_cloudwatch_logs  = var.enable_cloudwatch_logs
  enable_ssm_access       = var.enable_ssm_access
  lambda_custom_policies  = var.lambda_custom_policies
  
  tags = local.common_tags
}

# S3 Module
module "s3" {
  source = "../../modules/s3"
  
  environment                    = "dev"
  app_name                      = var.app_name
  create_app_bucket            = var.create_app_bucket
  create_static_website_bucket  = var.create_static_website_bucket
  create_backup_bucket          = var.create_backup_bucket
  enable_versioning            = var.enable_versioning
  enable_lifecycle             = var.enable_lifecycle
  lifecycle_expiration_days    = var.lifecycle_expiration_days
  lifecycle_noncurrent_expiration_days = var.lifecycle_noncurrent_expiration_days
  
  tags = local.common_tags
}

# EC2 Module
module "ec2" {
  source = "../../modules/ec2"
  
  environment                = "dev"
  vpc_id                    = module.vpc.vpc_id
  subnet_ids                = module.vpc.public_subnet_ids
  ami_type                  = var.ami_type
  instance_type             = var.instance_type
  instance_count            = var.instance_count
  public_key                = var.public_key
  allowed_ssh_cidrs         = var.allowed_ssh_cidrs
  custom_ingress_rules      = var.custom_ingress_rules
  associate_public_ip       = var.associate_public_ip
  user_data                 = var.user_data
  root_volume_type          = var.root_volume_type
  root_volume_size          = var.root_volume_size
  encrypt_root_volume       = var.encrypt_root_volume
  additional_ebs_volumes     = var.additional_ebs_volumes
  enable_alb                = var.enable_alb
  alb_internal              = var.alb_internal
  alb_deletion_protection   = var.alb_deletion_protection
  alb_target_port           = var.alb_target_port
  alb_health_check_path     = var.alb_health_check_path
  
  tags = local.common_tags
}

# RDS Module
module "rds" {
  source = "../../modules/rds"
  
  environment                        = "dev"
  vpc_id                           = module.vpc.vpc_id
  subnet_ids                       = module.vpc.private_subnet_ids
  create_rds                       = var.create_rds
  db_identifier                    = var.db_identifier
  db_engine                        = var.db_engine
  db_engine_version                = var.db_engine_version
  db_instance_class                = var.db_instance_class
  db_allocated_storage             = var.db_allocated_storage
  db_max_allocated_storage         = var.db_max_allocated_storage
  db_storage_type                  = var.db_storage_type
  db_storage_encrypted             = var.db_storage_encrypted
  db_name                          = var.db_name
  db_username                      = var.db_username
  db_password                      = var.db_password
  db_port                          = var.db_port
  db_publicly_accessible           = var.db_publicly_accessible
  db_backup_retention_period      = var.db_backup_retention_period
  db_backup_window                 = var.db_backup_window
  db_maintenance_window            = var.db_maintenance_window
  db_monitoring_interval           = var.db_monitoring_interval
  db_monitoring_role_arn          = var.db_monitoring_role_arn
  db_parameter_group_family        = var.db_parameter_group_family
  db_parameters                    = var.db_parameters
  db_deletion_protection           = var.db_deletion_protection
  db_skip_final_snapshot           = var.db_skip_final_snapshot
  db_performance_insights_enabled  = var.db_performance_insights_enabled
  db_performance_insights_retention_period = var.db_performance_insights_retention_period
  create_read_replica              = var.create_read_replica
  read_replica_instance_class      = var.read_replica_instance_class
  read_replica_publicly_accessible = var.read_replica_publicly_accessible
  allowed_cidr_blocks              = var.allowed_cidr_blocks
  allowed_security_group_ids       = [module.ec2.security_group_id]
  
  tags = local.common_tags
}
