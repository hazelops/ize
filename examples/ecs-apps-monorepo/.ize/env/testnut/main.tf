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
  cidr = "10.0.0.0/16"

  azs = [
    "us-east-1a",
    "us-east-1b",

  ]
  public_subnets = [
    "10.0.1.0/24",
    "10.0.2.0/24"
  ]

  private_subnets = [
    "10.0.3.0/24",
    "10.0.4.0/24"
  ]

  enable_nat_gateway = true
  single_nat_gateway                  = true
  enable_s3_endpoint                  = true
  enable_ecr_api_endpoint             = true
  ecr_api_endpoint_security_group_ids = [aws_security_group.default_permissive.id]
  enable_ecr_dkr_endpoint             = true
  ecr_dkr_endpoint_security_group_ids = [aws_security_group.default_permissive.id]
  enable_dns_hostnames                = true
  manage_default_network_acl          = true
  default_network_acl_name            = "${var.env}-${var.namespace}"
  tags = {
    Terraform = "true"
    Env       = var.env
  }
}

data "aws_route53_zone" "root" {
  name         = "${var.root_domain_name}."
  private_zone = false
}

resource "aws_route53_record" "env_ns_record" {
  zone_id = data.aws_route53_zone.root.id
  name    = "${var.env}.${var.root_domain_name}"
  type    = "NS"
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

module "ecs" {
  source  = "terraform-aws-modules/ecs/aws"
  version = "~> 2.0"
  name    = "${var.env}-${var.namespace}"

}

module "ec2_profile" {
  source      = "terraform-aws-modules/ecs/aws//modules/ecs-instance-profile"
  version     = "~> 2.0"
  name        = "${var.env}-${var.namespace}"
  include_ssm = true
}

