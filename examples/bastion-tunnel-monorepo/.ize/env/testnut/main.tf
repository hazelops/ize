resource "aws_key_pair" "root" {
  key_name   = var.ec2_key_pair_name
  public_key = var.ssh_public_key

  lifecycle {
    ignore_changes = [
      public_key
    ]
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "${var.env}-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["${var.aws_region}a"]
  public_subnets  = ["10.0.1.0/24"]
  private_subnets = ["10.0.2.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = true
}

module "bastion" {
  source  = "hazelops/ec2-openvpn-connector/aws"
  version = "~>0.4.1"

  vpn_enabled         = false
  env                 = var.env
  vpc_id              = module.vpc.vpc_id
  allowed_cidr_blocks = [module.vpc.vpc_cidr_block]
  private_subnets     = module.vpc.private_subnets

  ec2_key_pair_name = var.ec2_key_pair_name
  ssh_forward_rules = [
    "LocalForward 32084 info.cern.ch:80"
  ]
}
