# Development Environment Variables

# General Configuration
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "demo"
}

variable "owner" {
  description = "Owner of the resources"
  type        = string
  default     = "dev-team"
}

# VPC Configuration
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.20.0/24"]
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["ap-southeast-1a", "ap-southeast-1b"]
}

variable "enable_nat_gateway" {
  description = "Enable NAT Gateway for private subnets"
  type        = bool
  default     = true
}

# IAM Configuration
variable "create_ec2_role" {
  description = "Create EC2 instance role"
  type        = bool
  default     = true
}

variable "create_rds_role" {
  description = "Create RDS role"
  type        = bool
  default     = false
}

variable "create_lambda_role" {
  description = "Create Lambda execution role"
  type        = bool
  default     = false
}

variable "enable_s3_access" {
  description = "Enable S3 access for EC2 role"
  type        = bool
  default     = true
}

variable "enable_cloudwatch_logs" {
  description = "Enable CloudWatch Logs access for EC2 role"
  type        = bool
  default     = true
}

variable "enable_ssm_access" {
  description = "Enable SSM access for EC2 role"
  type        = bool
  default     = true
}

variable "lambda_custom_policies" {
  description = "Custom policies for Lambda role"
  type = list(object({
    Effect   = string
    Action   = list(string)
    Resource = list(string)
  }))
  default = []
}

# S3 Configuration
variable "app_name" {
  description = "Application name"
  type        = string
  default     = "demo"
}

variable "create_app_bucket" {
  description = "Create application S3 bucket"
  type        = bool
  default     = true
}

variable "create_static_website_bucket" {
  description = "Create static website S3 bucket"
  type        = bool
  default     = false
}

variable "create_backup_bucket" {
  description = "Create backup S3 bucket"
  type        = bool
  default     = false
}

variable "enable_versioning" {
  description = "Enable versioning for app bucket"
  type        = bool
  default     = true
}

variable "enable_lifecycle" {
  description = "Enable lifecycle configuration"
  type        = bool
  default     = true
}

variable "lifecycle_expiration_days" {
  description = "Days after which objects expire"
  type        = number
  default     = 90
}

variable "lifecycle_noncurrent_expiration_days" {
  description = "Days after which non-current versions expire"
  type        = number
  default     = 30
}

# EC2 Configuration
variable "ami_type" {
  description = "Type of AMI to use for EC2 instance"
  type        = string
  default     = "ubuntu_22_04_lts"
  
  validation {
    condition = contains([
      "ubuntu_22_04_lts",
      "ubuntu_20_04_lts",
      "amazon_linux_2",
      "amazon_linux_2023"
    ], var.ami_type)
    error_message = "AMI type must be one of the supported operating systems."
  }
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "instance_count" {
  description = "Number of EC2 instances to create"
  type        = number
  default     = 1
}

variable "public_key" {
  description = "Public key for EC2 key pair"
  type        = string
}

variable "allowed_ssh_cidrs" {
  description = "List of CIDR blocks allowed to SSH"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "custom_ingress_rules" {
  description = "Custom ingress rules for security group"
  type = list(object({
    from_port   = number
    to_port     = number
    protocol    = string
    cidr_blocks = list(string)
  }))
  default = []
}

variable "associate_public_ip" {
  description = "Associate public IP address"
  type        = bool
  default     = true
}

variable "user_data" {
  description = "User data script for EC2 instances"
  type        = string
  default     = ""
}

variable "root_volume_type" {
  description = "Type of root volume"
  type        = string
  default     = "gp3"
}

variable "root_volume_size" {
  description = "Size of root volume in GB"
  type        = number
  default     = 20
}

variable "encrypt_root_volume" {
  description = "Encrypt root volume"
  type        = bool
  default     = true
}

variable "additional_ebs_volumes" {
  description = "Additional EBS volumes"
  type = list(object({
    device_name           = string
    volume_type           = string
    volume_size          = number
    delete_on_termination = bool
    encrypted            = bool
  }))
  default = []
}

variable "enable_alb" {
  description = "Enable Application Load Balancer"
  type        = bool
  default     = false
}

variable "alb_internal" {
  description = "Internal ALB"
  type        = bool
  default     = false
}

variable "alb_deletion_protection" {
  description = "Enable ALB deletion protection"
  type        = bool
  default     = false
}

variable "alb_target_port" {
  description = "Target port for ALB"
  type        = number
  default     = 80
}

variable "alb_health_check_path" {
  description = "Health check path for ALB"
  type        = string
  default     = "/"
}

# RDS Configuration
variable "create_rds" {
  description = "Create RDS instance"
  type        = bool
  default     = false
}

variable "db_identifier" {
  description = "Identifier for the RDS instance"
  type        = string
  default     = "database"
}

variable "db_engine" {
  description = "Database engine"
  type        = string
  default     = "mysql"
}

variable "db_engine_version" {
  description = "Database engine version"
  type        = string
  default     = "8.0"
}

variable "db_instance_class" {
  description = "Database instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
  default     = 20
}

variable "db_max_allocated_storage" {
  description = "Maximum allocated storage in GB"
  type        = number
  default     = 100
}

variable "db_storage_type" {
  description = "Storage type"
  type        = string
  default     = "gp2"
}

variable "db_storage_encrypted" {
  description = "Encrypt storage"
  type        = bool
  default     = true
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "demo"
}

variable "db_username" {
  description = "Database username"
  type        = string
  default     = "admin"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "db_port" {
  description = "Database port"
  type        = number
  default     = 3306
}

variable "db_publicly_accessible" {
  description = "Make database publicly accessible"
  type        = bool
  default     = false
}

variable "db_backup_retention_period" {
  description = "Backup retention period in days"
  type        = number
  default     = 7
}

variable "db_backup_window" {
  description = "Backup window"
  type        = string
  default     = "03:00-04:00"
}

variable "db_maintenance_window" {
  description = "Maintenance window"
  type        = string
  default     = "sun:04:00-sun:05:00"
}

variable "db_monitoring_interval" {
  description = "Monitoring interval in seconds"
  type        = number
  default     = 0
}

variable "db_monitoring_role_arn" {
  description = "Monitoring role ARN"
  type        = string
  default     = null
}

variable "db_parameter_group_family" {
  description = "Parameter group family"
  type        = string
  default     = "mysql8.0"
}

variable "db_parameters" {
  description = "Database parameters"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "db_deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = false
}

variable "db_skip_final_snapshot" {
  description = "Skip final snapshot"
  type        = bool
  default     = true
}

variable "db_performance_insights_enabled" {
  description = "Enable Performance Insights"
  type        = bool
  default     = false
}

variable "db_performance_insights_retention_period" {
  description = "Performance Insights retention period in days"
  type        = number
  default     = 7
}

variable "create_read_replica" {
  description = "Create read replica"
  type        = bool
  default     = false
}

variable "read_replica_instance_class" {
  description = "Read replica instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "read_replica_publicly_accessible" {
  description = "Make read replica publicly accessible"
  type        = bool
  default     = false
}

variable "allowed_cidr_blocks" {
  description = "CIDR blocks allowed to access database"
  type        = list(string)
  default     = []
}
