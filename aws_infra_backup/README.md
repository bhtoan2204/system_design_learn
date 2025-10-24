# AWS Infrastructure với Terraform - Modules + Environments

Dự án này triển khai một infrastructure AWS hoàn chỉnh sử dụng mô hình "Modules + Environments" được khuyến nghị bởi HashiCorp.

## 🏗️ Cấu trúc Dự án

```
aws_infra/
├── modules/                    # Các module tái sử dụng
│   ├── vpc/                   # VPC, subnets, route tables, IGW, NAT Gateway
│   ├── ec2/                   # EC2 instances, security groups, ALB
│   ├── iam/                   # IAM roles, policies, instance profiles
│   ├── s3/                    # S3 buckets với các tính năng bảo mật
│   └── rds/                   # RDS instances, subnet groups, parameter groups
├── environments/               # Các môi trường triển khai
│   ├── dev/                   # Development environment
│   ├── staging/               # Staging environment
│   └── prod/                  # Production environment
└── README.md                  # Tài liệu này
```

## 📋 Yêu cầu

- **Terraform**: >= 1.5
- **AWS Provider**: >= 5.0
- **AWS CLI**: Cấu hình với credentials hợp lệ
- **SSH Key**: Public key để truy cập EC2 instances

## 🚀 Cách sử dụng

### 1. Chuẩn bị Backend Infrastructure

Trước khi triển khai, bạn cần tạo S3 bucket và DynamoDB table cho remote state:

```bash
# Sử dụng Terraform để tạo backend infrastructure
make setup-backend

# Hoặc manual:
cd backend-setup
terraform init
terraform plan
terraform apply
```

### 2. Sử dụng Makefile

```bash
# Xem tất cả commands có sẵn
make help

# Triển khai development environment
make dev-init
make dev-plan
make dev-apply

# Triển khai staging environment
make staging-init
make staging-plan
make staging-apply

# Triển khai production environment
make prod-init
make prod-plan
make prod-apply

# Validate tất cả environments
make all-validate

# Format tất cả files
make all-format
```

### 3. Triển khai Manual

```bash
# Development Environment
cd environments/dev
terraform init
terraform plan
terraform apply

# Staging Environment
cd environments/staging
terraform init
terraform plan
terraform apply

# Production Environment
cd environments/prod
terraform init
terraform plan
terraform apply
```

## 🔧 Modules Chi tiết

### VPC Module (`modules/vpc/`)

Tạo VPC với:
- Public và private subnets
- Internet Gateway
- NAT Gateway (optional)
- Route tables và associations
- Multi-AZ support

**Variables chính:**
- `vpc_cidr`: CIDR block cho VPC
- `public_subnet_cidrs`: CIDR blocks cho public subnets
- `private_subnet_cidrs`: CIDR blocks cho private subnets
- `availability_zones`: Danh sách AZ
- `enable_nat_gateway`: Bật/tắt NAT Gateway

### EC2 Module (`modules/ec2/`)

Tạo EC2 instances với:
- Multiple AMI support (Ubuntu, Amazon Linux)
- Security groups với custom rules
- Application Load Balancer (optional)
- EBS volumes với encryption
- User data scripts

**Variables chính:**
- `ami_type`: Loại AMI (ubuntu_22_04_lts, amazon_linux_2, etc.)
- `instance_type`: Instance type
- `instance_count`: Số lượng instances
- `public_key`: SSH public key
- `enable_alb`: Bật/tắt ALB

### IAM Module (`modules/iam/`)

Tạo IAM resources:
- EC2 instance roles
- RDS roles
- Lambda execution roles
- Custom policies
- Instance profiles

**Variables chính:**
- `create_ec2_role`: Tạo EC2 role
- `enable_s3_access`: Bật S3 access
- `enable_cloudwatch_logs`: Bật CloudWatch Logs
- `enable_ssm_access`: Bật SSM access

### S3 Module (`modules/s3/`)

Tạo S3 buckets với:
- Application bucket
- Static website bucket
- Backup bucket
- Versioning và lifecycle policies
- Encryption và public access blocks
- CORS configuration

**Variables chính:**
- `create_app_bucket`: Tạo app bucket
- `create_static_website_bucket`: Tạo static website bucket
- `create_backup_bucket`: Tạo backup bucket
- `enable_versioning`: Bật versioning
- `enable_lifecycle`: Bật lifecycle policies

### RDS Module (`modules/rds/`)

Tạo RDS instances với:
- Multiple engine support (MySQL, PostgreSQL, etc.)
- Subnet groups
- Parameter groups
- Security groups
- Read replicas (optional)
- Backup và monitoring

**Variables chính:**
- `create_rds`: Tạo RDS instance
- `db_engine`: Database engine
- `db_instance_class`: Instance class
- `create_read_replica`: Tạo read replica
- `db_deletion_protection`: Bật deletion protection

## 🌍 Environments

### Development (`environments/dev/`)
- **VPC CIDR**: 10.0.0.0/16
- **Instances**: 1x t3.micro
- **RDS**: Disabled
- **ALB**: Disabled
- **Features**: Basic setup cho development

### Staging (`environments/staging/`)
- **VPC CIDR**: 10.1.0.0/16
- **Instances**: 2x t3.micro
- **RDS**: Enabled (db.t3.micro)
- **ALB**: Enabled
- **Features**: Production-like setup cho testing

### Production (`environments/prod/`)
- **VPC CIDR**: 10.2.0.0/16
- **Instances**: 3x t3.small
- **RDS**: Enabled (db.t3.small) + Read Replica
- **ALB**: Enabled với deletion protection
- **Features**: High availability, monitoring, backup

## 🔐 Bảo mật

- **Encryption**: Tất cả EBS volumes và RDS storage được encrypt
- **Security Groups**: Restrictive rules với custom CIDR blocks
- **IAM**: Least privilege principle
- **S3**: Public access blocks và encryption
- **RDS**: Private subnets, không public access

## 📊 Monitoring & Logging

- **CloudWatch Logs**: EC2 instances có thể gửi logs
- **Performance Insights**: Enabled cho RDS (prod)
- **ALB Health Checks**: Configured cho load balancer
- **Backup**: Automated backups cho RDS

## 🏷️ Tagging Strategy

Tất cả resources được tag với:
- `Environment`: dev/staging/prod
- `Project`: demo
- `ManagedBy`: Terraform
- `Owner`: dev-team/staging-team/prod-team

## 🔄 CI/CD Integration

Có thể tích hợp với CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
name: Deploy Infrastructure
on:
  push:
    branches: [main]
jobs:
  deploy-dev:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Deploy Dev
        run: |
          make dev-init
          make dev-plan
          make dev-apply
```

## 🚨 Troubleshooting

### Common Issues:

1. **State Lock**: Nếu bị lock, có thể force unlock:
   ```bash
   terraform force-unlock <lock-id>
   ```

2. **Backend Configuration**: Đảm bảo S3 bucket và DynamoDB table tồn tại

3. **SSH Key**: Đảm bảo public key trong `terraform.tfvars` là hợp lệ

4. **Permissions**: Đảm bảo AWS credentials có đủ quyền

## 📈 Scaling

### Horizontal Scaling:
- Tăng `instance_count` trong `terraform.tfvars`
- ALB sẽ tự động distribute traffic

### Vertical Scaling:
- Thay đổi `instance_type` và `db_instance_class`
- Apply changes với `terraform plan` và `terraform apply`

## 💰 Cost Optimization

- **Development**: Sử dụng t3.micro instances
- **Staging**: Moderate resources với ALB
- **Production**: Right-sized instances với monitoring

## 📚 Tài liệu Tham khảo

- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Terraform Modules](https://www.terraform.io/docs/modules/index.html)
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)
- [HashiCorp Best Practices](https://www.terraform.io/docs/cloud/guides/recommended-practices/index.html)
