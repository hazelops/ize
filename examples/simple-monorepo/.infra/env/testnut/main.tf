resource "aws_key_pair" "root" {
  key_name = var.ec2_key_pair_name
  public_key = var.ssh_public_key
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.env}-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["us-east-1a"]
  private_subnets = ["10.0.1.0/24"]
  public_subnets  = ["10.0.101.0/24"]

  enable_nat_gateway = true

  tags = {
    Terraform = "true"
    Environment = var.env
  }
}

module "bastion" {
  source = "hazelops/ec2-bastion/aws"
  version = "~> 2.0"
  env = var.env
  aws_profile = var.aws_profile
  vpc_id = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnets
  ec2_key_pair_name = local.key_name
  ssh_forward_rules = [
    "LocalForward 31022 127.0.0.53:53"
  ]
}

output "cmd" {
  description = "Map of useful commands"
  value = {
    tunnel = module.bastion.cmd
  }
}

output "bastion_instance_id" {
  description = "Bastion EC2 instance ID"
  value       = module.bastion.instance_id
}

output "ssh_forward_config" {
  value = module.bastion.ssh_config
}
