# Basic usage example from https://github.com/hazelops/terraform-aws-ecs-app?tab=readme-ov-file#usage
module "squibby" {
  source  = "registry.terraform.io/hazelops/ecs-app/aws"
  version = "~>2.0"
#   name    = "${var.namespace}-squibby"
  name    = "squibby"

  env              = var.env
  ecs_cluster_name = module.ecs.ecs_cluster_name
  vpc_id           = module.vpc.vpc_id
  public_subnets   = module.vpc.public_subnets
  private_subnets  = module.vpc.private_subnets
  security_groups     = [aws_security_group.default_permissive.id]
  alb_security_groups = [aws_security_group.default_permissive.id]
  root_domain_name = var.root_domain_name
  zone_id          = aws_route53_zone.env_domain.id
  ecr_repo_create  = true

  https_enabled = false # Disabled for easier testing
  ecr_force_delete = true # Enabled for easier testing

  environment = {
    API_KEY   = "00000000000000000000000000000000"
    JWT_TOKEN = "99999999999999999999999999999999"
  }
}
