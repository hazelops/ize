# TODO: encapsulate everything into modules, like "api" or "job" which would use terraform locals

variable "env" {}

variable "namespace" {
  default = "testnut"
}

variable "aws_profile" {}

variable "aws_region" {
  default = "us-east-1"
}

variable "ssh_public_key" {}

variable "root_domain_name" {
  default = "examples.ize.sh"
}

variable "ec2_key_pair_name" {
  default = "dev-nutcorp"
}

variable "docker_registry" {}

variable "docker_image_tag" {}

variable "nat_gateway_enabled" {
  default = false
  description = "Set it to true to enable NAT Gateway, otherwise nat-instance module will be used"
}

variable "monitor_enabled" {
  default = true
}

# These are generic defaults. Feel free to reuse.
locals {
  env                  = var.env
  namespace            = var.namespace
  public_subnets       = module.vpc.public_subnets
  private_subnets      = module.vpc.private_subnets
  key_name             = aws_key_pair.root.key_name
  iam_instance_profile = module.ec2_profile.this_iam_instance_profile_id
  image_id             = data.aws_ami.amazon_linux_ecs_generic.id
  root_domain_name     = var.root_domain_name
  zone_id              = aws_route53_zone.env_domain.id
  vpc_id               = module.vpc.vpc_id
  security_groups      = [aws_security_group.default_permissive.id]
  alb_security_groups  = [aws_security_group.default_permissive.id]
  docker_registry      = var.docker_registry
  docker_image_tag     = var.docker_image_tag
#  tls_cert_arn         = length(module.env_acm.this_acm_certificate_arn) > 0 ? module.env_acm.this_acm_certificate_arn : null
  ecs_cluster_name     = module.ecs.this_ecs_cluster_name
}
