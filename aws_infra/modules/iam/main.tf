# IAM Module
# This module creates IAM roles, policies, and instance profiles

terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# EC2 Instance Role
resource "aws_iam_role" "ec2_role" {
  count = var.create_ec2_role ? 1 : 0
  
  name = "${var.environment}-ec2-role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
  
  tags = merge(var.tags, {
    Name = "${var.environment}-ec2-role"
  })
}

# EC2 Instance Profile
resource "aws_iam_instance_profile" "ec2_profile" {
  count = var.create_ec2_role ? 1 : 0
  
  name = "${var.environment}-ec2-profile"
  role = aws_iam_role.ec2_role[0].name
  
  tags = merge(var.tags, {
    Name = "${var.environment}-ec2-profile"
  })
}

# S3 Access Policy for EC2
resource "aws_iam_policy" "ec2_s3_policy" {
  count = var.create_ec2_role && var.enable_s3_access ? 1 : 0
  
  name        = "${var.environment}-ec2-s3-policy"
  description = "Policy for EC2 to access S3"
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::${var.environment}-*",
          "arn:aws:s3:::${var.environment}-*/*"
        ]
      }
    ]
  })
  
  tags = merge(var.tags, {
    Name = "${var.environment}-ec2-s3-policy"
  })
}

# CloudWatch Logs Policy for EC2
resource "aws_iam_policy" "ec2_cloudwatch_policy" {
  count = var.create_ec2_role && var.enable_cloudwatch_logs ? 1 : 0
  
  name        = "${var.environment}-ec2-cloudwatch-policy"
  description = "Policy for EC2 to write to CloudWatch Logs"
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams"
        ]
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })
  
  tags = merge(var.tags, {
    Name = "${var.environment}-ec2-cloudwatch-policy"
  })
}

# Attach policies to EC2 role
resource "aws_iam_role_policy_attachment" "ec2_s3_policy_attachment" {
  count = var.create_ec2_role && var.enable_s3_access ? 1 : 0
  
  role       = aws_iam_role.ec2_role[0].name
  policy_arn = aws_iam_policy.ec2_s3_policy[0].arn
}

resource "aws_iam_role_policy_attachment" "ec2_cloudwatch_policy_attachment" {
  count = var.create_ec2_role && var.enable_cloudwatch_logs ? 1 : 0
  
  role       = aws_iam_role.ec2_role[0].name
  policy_arn = aws_iam_policy.ec2_cloudwatch_policy[0].arn
}

# Attach AWS managed policies
resource "aws_iam_role_policy_attachment" "ec2_ssm_policy_attachment" {
  count = var.create_ec2_role && var.enable_ssm_access ? 1 : 0
  
  role       = aws_iam_role.ec2_role[0].name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# RDS Access Role
resource "aws_iam_role" "rds_role" {
  count = var.create_rds_role ? 1 : 0
  
  name = "${var.environment}-rds-role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "rds.amazonaws.com"
        }
      }
    ]
  })
  
  tags = merge(var.tags, {
    Name = "${var.environment}-rds-role"
  })
}

# Lambda Execution Role
resource "aws_iam_role" "lambda_role" {
  count = var.create_lambda_role ? 1 : 0
  
  name = "${var.environment}-lambda-role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
  
  tags = merge(var.tags, {
    Name = "${var.environment}-lambda-role"
  })
}

# Attach basic Lambda execution policy
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  count = var.create_lambda_role ? 1 : 0
  
  role       = aws_iam_role.lambda_role[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Custom Lambda policies
resource "aws_iam_policy" "lambda_custom_policy" {
  count = var.create_lambda_role && length(var.lambda_custom_policies) > 0 ? 1 : 0
  
  name        = "${var.environment}-lambda-custom-policy"
  description = "Custom policy for Lambda function"
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = var.lambda_custom_policies
  })
  
  tags = merge(var.tags, {
    Name = "${var.environment}-lambda-custom-policy"
  })
}

resource "aws_iam_role_policy_attachment" "lambda_custom_policy_attachment" {
  count = var.create_lambda_role && length(var.lambda_custom_policies) > 0 ? 1 : 0
  
  role       = aws_iam_role.lambda_role[0].name
  policy_arn = aws_iam_policy.lambda_custom_policy[0].arn
}
