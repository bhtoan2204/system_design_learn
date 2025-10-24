# EC2 Module Variables

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for EC2 instances"
  type        = list(string)
}

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

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
