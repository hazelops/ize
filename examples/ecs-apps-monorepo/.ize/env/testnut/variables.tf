variable "env" {}
variable "namespace" {}
variable "aws_profile" {}
variable "aws_region" {}
variable "ssh_public_key" {}
variable "root_domain_name" {}
variable "ec2_key_pair_name" {}
variable "docker_registry" {}
variable "docker_image_tag" {}

locals {
  env                  = var.env
  namespace            = var.namespace
  public_subnets       = module.vpc.public_subnets
  private_subnets      = module.vpc.private_subnets
  key_name             = aws_key_pair.root.key_name
  iam_instance_profile = module.ec2_profile.iam_instance_profile_id
  root_domain_name     = var.root_domain_name
  zone_id              = aws_route53_zone.env_domain.id
  vpc_id               = module.vpc.vpc_id
  security_groups      = [aws_security_group.default_permissive.id]
  alb_security_groups  = [aws_security_group.default_permissive.id]
  docker_registry      = var.docker_registry
  docker_image_tag     = var.docker_image_tag
  ecs_cluster_name     = module.ecs.ecs_cluster_name
}
