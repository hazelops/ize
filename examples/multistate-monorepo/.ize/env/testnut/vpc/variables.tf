variable "env" {}
variable "namespace" {}
variable "aws_profile" {}
variable "aws_region" {}
variable "ssh_public_key" {}
variable "ec2_key_pair_name" {}

locals {
  env                  = var.env
  namespace            = var.namespace
}
