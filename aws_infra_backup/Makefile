# Makefile for AWS Infrastructure - Modules + Environments
# This Makefile provides commands for managing multiple environments

MY_IP := $(shell curl -s ifconfig.me)/32
ENVIRONMENT ?= dev

# Default environment commands
.PHONY: help
help: ## Show this help message
	@echo "AWS Infrastructure Management"
	@echo "============================="
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Environment-specific commands:"
	@echo "  make dev-init     - Initialize dev environment"
	@echo "  make dev-plan     - Plan dev environment"
	@echo "  make dev-apply    - Apply dev environment"
	@echo "  make dev-destroy  - Destroy dev environment"
	@echo ""
	@echo "  make staging-init     - Initialize staging environment"
	@echo "  make staging-plan     - Plan staging environment"
	@echo "  make staging-apply    - Apply staging environment"
	@echo "  make staging-destroy  - Destroy staging environment"
	@echo ""
	@echo "  make prod-init     - Initialize prod environment"
	@echo "  make prod-plan     - Plan prod environment"
	@echo "  make prod-apply    - Apply prod environment"
	@echo "  make prod-destroy  - Destroy prod environment"
	@echo ""
	@echo "Usage: make <command> ENVIRONMENT=<env>"
	@echo "Example: make plan ENVIRONMENT=dev"

# Generic environment commands
init: ## Initialize Terraform for specified environment
	@echo "Initializing $(ENVIRONMENT) environment..."
	cd environments/$(ENVIRONMENT) && terraform init

plan: ## Plan Terraform for specified environment
	@echo "Planning $(ENVIRONMENT) environment..."
	cd environments/$(ENVIRONMENT) && terraform plan -var="allowed_ssh_cidrs=[\"$(MY_IP)\"]" -out tfplan.binary

apply: ## Apply Terraform for specified environment
	@echo "Applying $(ENVIRONMENT) environment..."
	cd environments/$(ENVIRONMENT) && terraform apply tfplan.binary

destroy: ## Destroy Terraform for specified environment
	@echo "Destroying $(ENVIRONMENT) environment..."
	cd environments/$(ENVIRONMENT) && terraform destroy -auto-approve -var="allowed_ssh_cidrs=[\"$(MY_IP)\"]"

show: ## Show Terraform plan for specified environment
	@echo "Showing $(ENVIRONMENT) environment plan..."
	cd environments/$(ENVIRONMENT) && terraform show -json tfplan.binary > plan.json

output: ## Show Terraform outputs for specified environment
	@echo "Showing $(ENVIRONMENT) environment outputs..."
	cd environments/$(ENVIRONMENT) && terraform output

validate: ## Validate Terraform configuration for specified environment
	@echo "Validating $(ENVIRONMENT) environment..."
	cd environments/$(ENVIRONMENT) && terraform validate

format: ## Format Terraform files for specified environment
	@echo "Formatting $(ENVIRONMENT) environment..."
	cd environments/$(ENVIRONMENT) && terraform fmt -recursive

# Development environment commands
dev-init: ## Initialize dev environment
	@$(MAKE) init ENVIRONMENT=dev

dev-plan: ## Plan dev environment
	@$(MAKE) plan ENVIRONMENT=dev

dev-apply: ## Apply dev environment
	@$(MAKE) apply ENVIRONMENT=dev

dev-destroy: ## Destroy dev environment
	@$(MAKE) destroy ENVIRONMENT=dev

dev-output: ## Show dev environment outputs
	@$(MAKE) output ENVIRONMENT=dev

dev-validate: ## Validate dev environment
	@$(MAKE) validate ENVIRONMENT=dev

dev-format: ## Format dev environment
	@$(MAKE) format ENVIRONMENT=dev

# Staging environment commands
staging-init: ## Initialize staging environment
	@$(MAKE) init ENVIRONMENT=staging

staging-plan: ## Plan staging environment
	@$(MAKE) plan ENVIRONMENT=staging

staging-apply: ## Apply staging environment
	@$(MAKE) apply ENVIRONMENT=staging

staging-destroy: ## Destroy staging environment
	@$(MAKE) destroy ENVIRONMENT=staging

staging-output: ## Show staging environment outputs
	@$(MAKE) output ENVIRONMENT=staging

staging-validate: ## Validate staging environment
	@$(MAKE) validate ENVIRONMENT=staging

staging-format: ## Format staging environment
	@$(MAKE) format ENVIRONMENT=staging

# Production environment commands
prod-init: ## Initialize prod environment
	@$(MAKE) init ENVIRONMENT=prod

prod-plan: ## Plan prod environment
	@$(MAKE) plan ENVIRONMENT=prod

prod-apply: ## Apply prod environment
	@$(MAKE) apply ENVIRONMENT=prod

prod-destroy: ## Destroy prod environment
	@$(MAKE) destroy ENVIRONMENT=prod

prod-output: ## Show prod environment outputs
	@$(MAKE) output ENVIRONMENT=prod

prod-validate: ## Validate prod environment
	@$(MAKE) validate ENVIRONMENT=prod

prod-format: ## Format prod environment
	@$(MAKE) format ENVIRONMENT=prod

# All environments commands
all-init: ## Initialize all environments
	@echo "Initializing all environments..."
	@$(MAKE) dev-init
	@$(MAKE) staging-init
	@$(MAKE) prod-init

all-plan: ## Plan all environments
	@echo "Planning all environments..."
	@$(MAKE) dev-plan
	@$(MAKE) staging-plan
	@$(MAKE) prod-plan

all-validate: ## Validate all environments
	@echo "Validating all environments..."
	@$(MAKE) dev-validate
	@$(MAKE) staging-validate
	@$(MAKE) prod-validate

all-format: ## Format all environments
	@echo "Formatting all environments..."
	@$(MAKE) dev-format
	@$(MAKE) staging-format
	@$(MAKE) prod-format

# Module commands
module-validate: ## Validate all modules
	@echo "Validating all modules..."
	cd modules/vpc && terraform validate
	cd modules/ec2 && terraform validate
	cd modules/iam && terraform validate
	cd modules/s3 && terraform validate
	cd modules/rds && terraform validate

module-format: ## Format all modules
	@echo "Formatting all modules..."
	cd modules/vpc && terraform fmt -recursive
	cd modules/ec2 && terraform fmt -recursive
	cd modules/iam && terraform fmt -recursive
	cd modules/s3 && terraform fmt -recursive
	cd modules/rds && terraform fmt -recursive

# Cleanup commands
clean: ## Clean up temporary files
	@echo "Cleaning up temporary files..."
	find . -name "*.tfstate*" -not -path "./environments/*" -delete
	find . -name "tfplan.binary" -delete
	find . -name "plan.json" -delete
	find . -name ".terraform" -type d -exec rm -rf {} + 2>/dev/null || true
	find . -name ".terraform.lock.hcl" -delete

# Security commands
check-security: ## Run security checks on Terraform files
	@echo "Running security checks..."
	@if command -v tfsec >/dev/null 2>&1; then \
		tfsec .; \
	else \
		echo "tfsec not installed. Install with: go install github.com/aquasecurity/tfsec/cmd/tfsec@latest"; \
	fi

# Cost estimation (requires infracost)
cost-estimate: ## Estimate costs for specified environment
	@echo "Estimating costs for $(ENVIRONMENT) environment..."
	@if command -v infracost >/dev/null 2>&1; then \
		cd environments/$(ENVIRONMENT) && infracost breakdown --path .; \
	else \
		echo "infracost not installed. Install with: brew install infracost"; \
	fi

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v terraform-docs >/dev/null 2>&1; then \
		find . -name "*.tf" -not -path "./environments/*" -exec terraform-docs markdown table {} \; > DOCS.md; \
	else \
		echo "terraform-docs not installed. Install with: go install github.com/terraform-docs/terraform-docs@latest"; \
	fi