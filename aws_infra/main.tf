provider "aws" {
  region = "ap-southeast-1"
}

resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name = "demo-vpc"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "demo-igw"
  }
}

resource "aws_subnet" "public_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "ap-southeast-1a"
  map_public_ip_on_launch = true
  tags = {
    Name = "demo-public-a"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
  tags = {
    Name = "demo-public-rt"
  }
}

resource "aws_route_table_association" "public_a" {
  subnet_id      = aws_subnet.public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_key_pair" "ec2_key_pair" {
  key_name = "ec2-key-pair"
  public_key = file("~/.ssh/id_rsa.pub")
}

resource "aws_instance" "ec2" {
  ami = "ami-0ffd8e96d1336b6ac" // amazon linux 2023
  instance_type = "t3.micro"
  key_name = aws_key_pair.ec2_key_pair.key_name
  subnet_id = aws_subnet.public_a.id
  vpc_security_group_ids = [aws_security_group.ec2_security_group.id]
  associate_public_ip_address = true
  tags = {
    Name = "ec2"
  }
}

resource "aws_security_group" "ec2_security_group" {
  name = "ec2-security-group"
  description = "EC2 security group"
  vpc_id = aws_vpc.main.id
  
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = [var.allowed_ssh_cidr]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

variable "allowed_ssh_cidr" {
  description = "CIDR allowed to SSH (e.g. 203.0.113.5/32)"
  type        = string
}

output "ec2_public_ip" {
  value = aws_instance.ec2.public_ip
}