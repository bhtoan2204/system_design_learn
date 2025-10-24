# IAM Module Outputs

output "ec2_role_arn" {
  description = "ARN of the EC2 instance role"
  value       = var.create_ec2_role ? aws_iam_role.ec2_role[0].arn : null
}

output "ec2_role_name" {
  description = "Name of the EC2 instance role"
  value       = var.create_ec2_role ? aws_iam_role.ec2_role[0].name : null
}

output "ec2_instance_profile_arn" {
  description = "ARN of the EC2 instance profile"
  value       = var.create_ec2_role ? aws_iam_instance_profile.ec2_profile[0].arn : null
}

output "ec2_instance_profile_name" {
  description = "Name of the EC2 instance profile"
  value       = var.create_ec2_role ? aws_iam_instance_profile.ec2_profile[0].name : null
}

output "rds_role_arn" {
  description = "ARN of the RDS role"
  value       = var.create_rds_role ? aws_iam_role.rds_role[0].arn : null
}

output "rds_role_name" {
  description = "Name of the RDS role"
  value       = var.create_rds_role ? aws_iam_role.rds_role[0].name : null
}

output "lambda_role_arn" {
  description = "ARN of the Lambda execution role"
  value       = var.create_lambda_role ? aws_iam_role.lambda_role[0].arn : null
}

output "lambda_role_name" {
  description = "Name of the Lambda execution role"
  value       = var.create_lambda_role ? aws_iam_role.lambda_role[0].name : null
}
