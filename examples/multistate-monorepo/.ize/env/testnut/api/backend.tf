provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "testnut"
      namespace = "testnut"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "testnut-tf-state"
    key            = "testnut/api.tfstate"
    region         = "us-east-1"
    profile        = "default"
    dynamodb_table = "tf-state-lock"
  }
}
