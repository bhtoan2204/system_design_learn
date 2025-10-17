MY_IP := $(shell curl -s ifconfig.me)/32

init:
	terraform init

plan:
	terraform plan -var="allowed_ssh_cidr=$(MY_IP)" -out tfplan.binary

show:
	terraform show -json tfplan.binary > plan.json

cost-estimate:
	infracost breakdown --path plan.json

apply:
	terraform apply tfplan.binary

# refresh-apply:
# 	terraform apply -refresh-only tfplan.binary

destroy:
	terraform destroy -auto-approve -var="allowed_ssh_cidr=$(MY_IP)"

.PHONY: init plan apply destroy show cost-estimate