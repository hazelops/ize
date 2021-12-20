variable "aws_region" {
  default = "us-east-1"
}

variable "aws_profile" {
  default = "default"
}

variable "env" {
  default = "testnut"
}

variable "ec2_key_pair_name" {}

variable "ssh_public_key" {}

locals {
  key_name             = aws_key_pair.root.key_name
}
