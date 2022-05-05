module "squibby" {
  depends_on = [
    module.ecs
  ]
  source  = "hazelops/ecs-app/aws"
  version = "~>1.1"

  name             = "squibby"
  app_type         = "web"
  env              = var.env
  namespace        = var.namespace
  ecs_cluster_name = local.ecs_cluster_name

  # Containers
  docker_registry      = local.docker_registry
  image_id             = local.image_id
  docker_image_tag     = local.docker_image_tag
  iam_instance_profile = local.iam_instance_profile
  key_name             = local.key_name

  # Load Balancer
  public                = true
  alb_health_check_path = "/"
  alb_security_groups   = local.alb_security_groups

  # Network
  vpc_id                       = local.vpc_id
  public_subnets               = local.public_subnets
  private_subnets              = local.private_subnets
  security_groups              = local.security_groups
  root_domain_name             = var.root_domain_name
  zone_id                      = local.zone_id
  route53_health_check_enabled = false
  sns_service_subscription_endpoint = "nutcorp-ops@hazelops.com"
  sns_service_subscription_endpoint_protocol = "email"
  domain_names = [
    "squibby.${var.root_domain_name}"
  ]

  # Environment variables
  app_secrets = [
  ]
  environment = {
    ENV      = var.env
    APP_NAME = "Craben"
  }
}
