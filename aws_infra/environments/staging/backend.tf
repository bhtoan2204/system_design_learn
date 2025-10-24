# Backend Configuration for Staging Environment
# This file configures the remote backend for state management

terraform {
  backend "s3" {
    # S3 bucket for storing Terraform state
    bucket = "terraform-state-staging-demo"
    
    # Key path for the state file
    key = "staging/terraform.tfstate"
    
    # AWS region
    region = "ap-southeast-1"
    
    # DynamoDB table for state locking
    dynamodb_table = "terraform-state-locks-staging"
    
    # Enable encryption
    encrypt = true
    
    # Optional: AWS profile to use
    # profile = "default"
  }
}
