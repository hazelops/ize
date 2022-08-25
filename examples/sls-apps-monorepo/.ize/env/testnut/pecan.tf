module "pecan" {
  source     = "registry.terraform.io/hazelops/sls-app/aws"
  version    = "~> 0.2"
  name       = "pecan"
  parameters = {
    ROOT_DOMAIN_NAME = var.root_domain_name
  }
}
