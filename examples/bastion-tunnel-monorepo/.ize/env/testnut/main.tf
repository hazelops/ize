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
  version = "~> 2.0"

  name = "${var.env}-vpc"
  cidr = "10.30.0.0/16"

  azs = [
    "us-east-1a"
  ]
  public_subnets = [
    "10.30.10.0/23"
  ]

  private_subnets = [
    "10.30.20.0/23"
  ]

  enable_nat_gateway = var.nat_gateway_enabled
  //  single_nat_gateway = var.env == "prod" ? false : true
  single_nat_gateway                  = true
  enable_vpn_gateway                  = false
  enable_s3_endpoint                  = true
  enable_ecr_api_endpoint             = false
  ecr_api_endpoint_security_group_ids = [aws_security_group.default_permissive.id]
  enable_ecr_dkr_endpoint             = false
  ecr_dkr_endpoint_security_group_ids = [aws_security_group.default_permissive.id]
  enable_dns_hostnames                = true
  enable_dns_support                  = true
  manage_default_network_acl          = true
  default_network_acl_name            = "${var.env}-${var.namespace}"
  tags = {
    Terraform = "true"
    Env       = var.env
  }
}

# nat-instance - use it when you want to save costs on AWS NAT Gateway. Use it only in test needs.
module "nat_instance" {
  source  = "hazelops/ec2-nat/aws"
  version = "~> 2.0"

  enabled = var.nat_gateway_enabled ? false : true

  env                    = var.env
  vpc_id                 = module.vpc.vpc_id
  allowed_cidr_blocks    = [module.vpc.vpc_cidr_block]
  public_subnets         = module.vpc.public_subnets
  private_route_table_id = module.vpc.private_route_table_ids[0]
  ec2_key_pair_name      = local.key_name
}
###### End of the nat-instance block

data "aws_route53_zone" "root" {
  name         = "${var.root_domain_name}."
  private_zone = false
}

resource "aws_route53_record" "env_ns_record" {
  zone_id = data.aws_route53_zone.root.id
  name    = "${var.env}.${var.root_domain_name}"
  type    = "NS"
  // Fast TTL for tests
  ttl     = "60"
  records = aws_route53_zone.env_domain.name_servers
}


resource "aws_route53_zone" "env_domain" {
  name = "${var.env}.${var.root_domain_name}"
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

module "openvpn_instance" {
  source  = "hazelops/ec2-openvpn-connector/aws"
  version = "~>0.2"

  vpn_enabled     = false
  bastion_enabled = true

  env                 = var.env
  vpc_id              = module.vpc.vpc_id
  allowed_cidr_blocks = [module.vpc.vpc_cidr_block]
  private_subnets     = module.vpc.private_subnets
  ec2_key_pair_name   = local.key_name
  ssh_forward_rules = [
    // Test outputs
    "LocalForward 32084 info.cern.ch:80"
  ]
}
