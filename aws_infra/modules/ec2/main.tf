# EC2 Module
# This module creates EC2 instances, security groups, and key pairs

terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

# Data sources for AMIs
data "aws_ami" "ubuntu_22_04_lts" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "aws_ami" "ubuntu_20_04_lts" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# AMI mapping
locals {
  ami_mapping = {
    ubuntu_22_04_lts  = data.aws_ami.ubuntu_22_04_lts.id
    ubuntu_20_04_lts  = data.aws_ami.ubuntu_20_04_lts.id
    amazon_linux_2    = data.aws_ami.amazon_linux_2.id
    amazon_linux_2023 = data.aws_ami.amazon_linux_2023.id
  }
  
  selected_ami_id = local.ami_mapping[var.ami_type]
}

# Key Pair
resource "aws_key_pair" "ec2_key_pair" {
  key_name   = "${var.environment}-ec2-key-pair"
  public_key = var.public_key
  
  tags = merge(var.tags, {
    Name = "${var.environment}-ec2-key-pair"
  })
}

# Security Group for EC2
resource "aws_security_group" "ec2_security_group" {
  name        = "${var.environment}-ec2-security-group"
  description = "Security group for EC2 instances"
  vpc_id      = var.vpc_id
  
  # SSH access
  dynamic "ingress" {
    for_each = var.allowed_ssh_cidrs
    content {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
    }
  }
  
  # HTTP access
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  # HTTPS access
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  # Custom ports
  dynamic "ingress" {
    for_each = var.custom_ingress_rules
    content {
      from_port   = ingress.value.from_port
      to_port     = ingress.value.to_port
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.cidr_blocks
    }
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = merge(var.tags, {
    Name = "${var.environment}-ec2-security-group"
  })
}

# EC2 Instances
resource "aws_instance" "ec2" {
  count = var.instance_count
  
  ami                    = local.selected_ami_id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.ec2_key_pair.key_name
  subnet_id              = var.subnet_ids[count.index % length(var.subnet_ids)]
  vpc_security_group_ids = [aws_security_group.ec2_security_group.id]
  
  associate_public_ip_address = var.associate_public_ip
  
  # User data script
  user_data = var.user_data
  
  # Root volume
  root_block_device {
    volume_type           = var.root_volume_type
    volume_size          = var.root_volume_size
    delete_on_termination = true
    encrypted            = var.encrypt_root_volume
    
    tags = merge(var.tags, {
      Name = "${var.environment}-ec2-${count.index + 1}-root"
    })
  }
  
  # Additional EBS volumes
  dynamic "ebs_block_device" {
    for_each = var.additional_ebs_volumes
    content {
      device_name           = ebs_block_device.value.device_name
      volume_type           = ebs_block_device.value.volume_type
      volume_size          = ebs_block_device.value.volume_size
      delete_on_termination = ebs_block_device.value.delete_on_termination
      encrypted            = ebs_block_device.value.encrypted
    }
  }
  
  tags = merge(var.tags, {
    Name = "${var.environment}-ec2-${count.index + 1}"
    Type = "EC2"
  })
}

# Application Load Balancer (optional)
resource "aws_lb" "alb" {
  count = var.enable_alb ? 1 : 0
  
  name               = "${var.environment}-alb"
  internal           = var.alb_internal
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_security_group[0].id]
  subnets            = var.subnet_ids
  
  enable_deletion_protection = var.alb_deletion_protection
  
  tags = merge(var.tags, {
    Name = "${var.environment}-alb"
  })
}

# ALB Security Group
resource "aws_security_group" "alb_security_group" {
  count = var.enable_alb ? 1 : 0
  
  name        = "${var.environment}-alb-security-group"
  description = "Security group for Application Load Balancer"
  vpc_id      = var.vpc_id
  
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = merge(var.tags, {
    Name = "${var.environment}-alb-security-group"
  })
}

# ALB Target Group
resource "aws_lb_target_group" "alb_target_group" {
  count = var.enable_alb ? 1 : 0
  
  name     = "${var.environment}-tg"
  port     = var.alb_target_port
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  
  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = var.alb_health_check_path
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }
  
  tags = merge(var.tags, {
    Name = "${var.environment}-target-group"
  })
}

# ALB Target Group Attachment
resource "aws_lb_target_group_attachment" "alb_target_group_attachment" {
  count = var.enable_alb ? var.instance_count : 0
  
  target_group_arn = aws_lb_target_group.alb_target_group[0].arn
  target_id        = aws_instance.ec2[count.index].id
  port             = var.alb_target_port
}

# ALB Listener
resource "aws_lb_listener" "alb_listener" {
  count = var.enable_alb ? 1 : 0
  
  load_balancer_arn = aws_lb.alb[0].arn
  port              = "80"
  protocol          = "HTTP"
  
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.alb_target_group[0].arn
  }
}
