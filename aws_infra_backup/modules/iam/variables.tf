# IAM Module Variables

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

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

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
