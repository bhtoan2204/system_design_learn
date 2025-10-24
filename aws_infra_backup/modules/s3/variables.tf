# S3 Module Variables

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

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

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
