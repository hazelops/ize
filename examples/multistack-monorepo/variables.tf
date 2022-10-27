variable "env" {}
variable "namespace" {}
variable "aws_region" {}

locals {
  env        = var.env
  namespace  = var.namespace
  aws_region = var.aws_region
}
