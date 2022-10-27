module "s3_bucket" {
  source  = "registry.terraform.io/terraform-aws-modules/s3-bucket/aws"
  version = "~>3.0"

  bucket = "${var.env}-${var.namespace}-test-s3-bucket"
  acl    = "private"

  versioning = {
    enabled = false
  }
  force_destroy = true
}
