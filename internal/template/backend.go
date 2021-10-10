package template

var backendTemplate = `{{if or (.ENV | strings.Contains "localstack") (.ENV | strings.Contains "local") }}terraform {
  backend "local" {}
}

provider "aws" {
  profile                     = var.aws_profile
  region                      = var.aws_region
  s3_force_path_style         = true
  secret_key                  = "mock_secret_key"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway     = "{{ .LOCALSTACK_ENDPOINT }}"
    acm            = "{{ .LOCALSTACK_ENDPOINT }}"
    cloudformation = "{{ .LOCALSTACK_ENDPOINT }}"
    cloudwatch     = "{{ .LOCALSTACK_ENDPOINT }}"
    ec2            = "{{ .LOCALSTACK_ENDPOINT }}"
    dynamodb       = "{{ .LOCALSTACK_ENDPOINT }}"
    es             = "{{ .LOCALSTACK_ENDPOINT }}"
    firehose       = "{{ .LOCALSTACK_ENDPOINT }}"
    iam            = "{{ .LOCALSTACK_ENDPOINT }}"
    kinesis        = "{{ .LOCALSTACK_ENDPOINT }}"
    lambda         = "{{ .LOCALSTACK_ENDPOINT }}"
    route53        = "{{ .LOCALSTACK_ENDPOINT }}"
    redshift       = "{{ .LOCALSTACK_ENDPOINT }}"
    s3             = "{{ .LOCALSTACK_ENDPOINT }}"
    secretsmanager = "{{ .LOCALSTACK_ENDPOINT }}"
    ses            = "{{ .LOCALSTACK_ENDPOINT }}"
    sns            = "{{ .LOCALSTACK_ENDPOINT }}"
    sqs            = "{{ .LOCALSTACK_ENDPOINT }}"
    ssm            = "{{ .LOCALSTACK_ENDPOINT }}"
    stepfunctions  = "{{ .LOCALSTACK_ENDPOINT }}"
    sts            = "{{ .LOCALSTACK_ENDPOINT }}"
    ecs            = "{{ .LOCALSTACK_ENDPOINT }}"
    ecr            = "{{ .LOCALSTACK_ENDPOINT }}"
  }
}{{else}}provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

terraform {
  backend "s3" {
    bucket         = "{{if .TERRAFORM_STATE_BUCKET_NAME}}{{ .TERRAFORM_STATE_BUCKET_NAME }}{{else}}nutcorpnet-tf-state{{end}}"
    key            = "{{if .TERRAFORM_STATE_KEY}}{{ .TERRAFORM_STATE_KEY }}{{else}}foo/terraform.tfstate{{end}}"
    region         = "{{if .TERRAFORM_STATE_REGION}}{{ .TERRAFORM_STATE_REGION }}{{else}}us-east-1{{end}}"
    profile        = "{{if .TERRAFORM_STATE_PROFILE}}{{ .TERRAFORM_STATE_PROFILE }}{{else}}nutcorp-dev{{end}}"
    dynamodb_table = "{{if .TERRAFORM_STATE_DYNAMODB_TABLE}}{{ .TERRAFORM_STATE_DYNAMODB_TABLE }}{{else}}tf-state-lock{{end}}"
  }
}{{end}}`
