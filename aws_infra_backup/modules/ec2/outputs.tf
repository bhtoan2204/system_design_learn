# EC2 Module Outputs

output "instance_ids" {
  description = "IDs of the EC2 instances"
  value       = aws_instance.ec2[*].id
}

output "instance_public_ips" {
  description = "Public IPs of the EC2 instances"
  value       = aws_instance.ec2[*].public_ip
}

output "instance_private_ips" {
  description = "Private IPs of the EC2 instances"
  value       = aws_instance.ec2[*].private_ip
}

output "instance_arns" {
  description = "ARNs of the EC2 instances"
  value       = aws_instance.ec2[*].arn
}

output "key_pair_name" {
  description = "Name of the EC2 key pair"
  value       = aws_key_pair.ec2_key_pair.key_name
}

output "security_group_id" {
  description = "ID of the EC2 security group"
  value       = aws_security_group.ec2_security_group.id
}

output "alb_arn" {
  description = "ARN of the Application Load Balancer"
  value       = var.enable_alb ? aws_lb.alb[0].arn : null
}

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = var.enable_alb ? aws_lb.alb[0].dns_name : null
}

output "alb_zone_id" {
  description = "Zone ID of the Application Load Balancer"
  value       = var.enable_alb ? aws_lb.alb[0].zone_id : null
}

output "target_group_arn" {
  description = "ARN of the target group"
  value       = var.enable_alb ? aws_lb_target_group.alb_target_group[0].arn : null
}
