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
  version = "~> 3.0"

  name = "${var.env}-vpc"
  cidr = "10.0.0.0/16"

  azs = [
    "us-east-1a"
  ]
  public_subnets = [
    "10.0.1.0/24"
  ]

  private_subnets = [
    "10.0.2.0/24"
  ]

  manage_default_network_acl = true
  default_network_acl_name   = "${var.env}-${var.namespace}"

}

module "nat_instance" {
  source                 = "hazelops/ec2-nat/aws"
  version                = "~> 2.0"
  enabled                = true
  env                    = var.env
  vpc_id                 = module.vpc.vpc_id
  allowed_cidr_blocks    = [module.vpc.vpc_cidr_block]
  public_subnets         = module.vpc.public_subnets
  private_route_table_id = module.vpc.private_route_table_ids[0]
  ec2_key_pair_name      = local.key_name
}

resource "aws_security_group" "default_permissive" {
  name        = "${var.env}-default-permissive"
  vpc_id      = module.vpc.vpc_id
  description = "Managed by Terraform"

  ingress {
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Terraform = "true"
    Env       = var.env
    Name      = "${var.env}-default-permissive"
  }
}

module "ec2_profile" {
  source      = "terraform-aws-modules/ecs/aws//modules/ecs-instance-profile"
  version     = "~> 2.0"
  name        = "${var.env}-${var.namespace}"
  include_ssm = true
}

module "bastion" {
  source  = "hazelops/ec2-openvpn-connector/aws"
  version = "~>0.2"

  vpn_enabled         = false
  env                 = var.env
  vpc_id              = module.vpc.vpc_id
  allowed_cidr_blocks = [module.vpc.vpc_cidr_block]
  private_subnets     = module.vpc.private_subnets
  ec2_key_pair_name   = local.key_name
  ssh_forward_rules = [
    "LocalForward 32084 info.cern.ch:80"
  ]
}
