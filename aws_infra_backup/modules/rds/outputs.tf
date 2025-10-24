# RDS Module Outputs

output "db_instance_id" {
  description = "ID of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].id : null
}

output "db_instance_arn" {
  description = "ARN of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].arn : null
}

output "db_instance_endpoint" {
  description = "Endpoint of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].endpoint : null
}

output "db_instance_hosted_zone_id" {
  description = "Hosted zone ID of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].hosted_zone_id : null
}

output "db_instance_address" {
  description = "Address of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].address : null
}

output "db_instance_port" {
  description = "Port of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].port : null
}

output "db_instance_name" {
  description = "Name of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].name : null
}

output "db_instance_username" {
  description = "Username of the RDS instance"
  value       = var.create_rds ? aws_db_instance.main[0].username : null
}

output "db_subnet_group_id" {
  description = "ID of the DB subnet group"
  value       = var.create_rds ? aws_db_subnet_group.main[0].id : null
}

output "db_subnet_group_arn" {
  description = "ARN of the DB subnet group"
  value       = var.create_rds ? aws_db_subnet_group.main[0].arn : null
}

output "db_parameter_group_id" {
  description = "ID of the DB parameter group"
  value       = var.create_rds ? aws_db_parameter_group.main[0].id : null
}

output "db_parameter_group_arn" {
  description = "ARN of the DB parameter group"
  value       = var.create_rds ? aws_db_parameter_group.main[0].arn : null
}

output "db_security_group_id" {
  description = "ID of the RDS security group"
  value       = var.create_rds ? aws_security_group.rds_security_group[0].id : null
}

output "read_replica_id" {
  description = "ID of the read replica"
  value       = var.create_rds && var.create_read_replica ? aws_db_instance.read_replica[0].id : null
}

output "read_replica_arn" {
  description = "ARN of the read replica"
  value       = var.create_rds && var.create_read_replica ? aws_db_instance.read_replica[0].arn : null
}

output "read_replica_endpoint" {
  description = "Endpoint of the read replica"
  value       = var.create_rds && var.create_read_replica ? aws_db_instance.read_replica[0].endpoint : null
}
