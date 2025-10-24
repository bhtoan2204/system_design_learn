# RDS Module
# This module creates RDS instances and related resources

terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  count = var.create_rds ? 1 : 0
  
  name       = "${var.environment}-db-subnet-group"
  subnet_ids = var.subnet_ids
  
  tags = merge(var.tags, {
    Name = "${var.environment}-db-subnet-group"
  })
}

# DB Parameter Group
resource "aws_db_parameter_group" "main" {
  count = var.create_rds ? 1 : 0
  
  family = var.db_parameter_group_family
  name   = "${var.environment}-db-parameter-group"
  
  dynamic "parameter" {
    for_each = var.db_parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }
  
  tags = merge(var.tags, {
    Name = "${var.environment}-db-parameter-group"
  })
}

# DB Security Group
resource "aws_security_group" "rds_security_group" {
  count = var.create_rds ? 1 : 0
  
  name        = "${var.environment}-rds-security-group"
  description = "Security group for RDS instance"
  vpc_id      = var.vpc_id
  
  # Database access from EC2 instances
  dynamic "ingress" {
    for_each = var.allowed_cidr_blocks
    content {
      from_port   = var.db_port
      to_port     = var.db_port
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
    }
  }
  
  # Database access from security groups
  dynamic "ingress" {
    for_each = var.allowed_security_group_ids
    content {
      from_port       = var.db_port
      to_port         = var.db_port
      protocol        = "tcp"
      security_groups = [ingress.value]
    }
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = merge(var.tags, {
    Name = "${var.environment}-rds-security-group"
  })
}

# RDS Instance
resource "aws_db_instance" "main" {
  count = var.create_rds ? 1 : 0
  
  identifier = "${var.environment}-${var.db_identifier}"
  
  # Engine configuration
  engine         = var.db_engine
  engine_version = var.db_engine_version
  instance_class = var.db_instance_class
  
  # Storage configuration
  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = var.db_max_allocated_storage
  storage_type          = var.db_storage_type
  storage_encrypted     = var.db_storage_encrypted
  
  # Database configuration
  db_name  = var.db_name
  username = var.db_username
  password = var.db_password
  port     = var.db_port
  
  # Network configuration
  db_subnet_group_name   = aws_db_subnet_group.main[0].name
  vpc_security_group_ids = [aws_security_group.rds_security_group[0].id]
  publicly_accessible    = var.db_publicly_accessible
  
  # Backup configuration
  backup_retention_period = var.db_backup_retention_period
  backup_window          = var.db_backup_window
  maintenance_window     = var.db_maintenance_window
  
  # Monitoring configuration
  monitoring_interval = var.db_monitoring_interval
  monitoring_role_arn = var.db_monitoring_role_arn
  
  # Parameter group
  parameter_group_name = aws_db_parameter_group.main[0].name
  
  # Deletion protection
  deletion_protection = var.db_deletion_protection
  skip_final_snapshot = var.db_skip_final_snapshot
  final_snapshot_identifier = var.db_skip_final_snapshot ? null : "${var.environment}-${var.db_identifier}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  
  # Performance Insights
  performance_insights_enabled = var.db_performance_insights_enabled
  performance_insights_retention_period = var.db_performance_insights_enabled ? var.db_performance_insights_retention_period : null
  
  tags = merge(var.tags, {
    Name = "${var.environment}-${var.db_identifier}"
  })
}

# RDS Read Replica (optional)
resource "aws_db_instance" "read_replica" {
  count = var.create_rds && var.create_read_replica ? 1 : 0
  
  identifier = "${var.environment}-${var.db_identifier}-replica"
  
  # Replica configuration
  replicate_source_db = aws_db_instance.main[0].identifier
  instance_class      = var.read_replica_instance_class
  
  # Storage configuration
  storage_encrypted = var.db_storage_encrypted
  
  # Network configuration
  vpc_security_group_ids = [aws_security_group.rds_security_group[0].id]
  publicly_accessible    = var.read_replica_publicly_accessible
  
  # Monitoring configuration
  monitoring_interval = var.db_monitoring_interval
  monitoring_role_arn = var.db_monitoring_role_arn
  
  # Performance Insights
  performance_insights_enabled = var.db_performance_insights_enabled
  performance_insights_retention_period = var.db_performance_insights_enabled ? var.db_performance_insights_retention_period : null
  
  tags = merge(var.tags, {
    Name = "${var.environment}-${var.db_identifier}-replica"
  })
}
