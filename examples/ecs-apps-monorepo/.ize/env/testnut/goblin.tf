module "goblin" {
  depends_on = [
    module.ecs
  ]
  source  = "registry.terraform.io/hazelops/ecs-app/aws"
  version = "~>1.4"

  name             = "goblin"
  app_type         = "web"
  env              = var.env
  namespace        = var.namespace
  ecs_cluster_name = local.ecs_cluster_name

  # Containers
  docker_registry      = local.docker_registry
#  image_id             = local.image_id
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
#  route53_health_check_enabled = false
#  sns_service_subscription_endpoint = "nutcorp-ops@hazelops.com"
#  sns_service_subscription_endpoint_protocol = "email"
  domain_names = [
    "goblin.${var.root_domain_name}"
  ]
  https_enabled = false

  # Environment variables
  app_secrets = [
    "EXAMPLE_SECRET"
  ]
  environment = {
    EXAMPLE_API_KEY = "Api Key"
  }
}
