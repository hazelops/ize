module "cronjob" {
  depends_on = [
    module.ecs
  ]
  source  = "registry.terraform.io/hazelops/ecs-app/aws"
  version = "~>1.4"

  name             = "cronjob"
  app_type         = "worker"
  env              = var.env
  namespace        = var.namespace
  ecs_cluster_name = local.ecs_cluster_name

  # Containers
  docker_registry      = local.docker_registry
  docker_image_tag     = local.docker_image_tag
  iam_instance_profile = local.iam_instance_profile
  key_name             = local.key_name

  # Network
  vpc_id          = local.vpc_id
  public_subnets  = local.public_subnets
  private_subnets = local.private_subnets
  security_groups = local.security_groups

  # Environment variables
  environment = {
    APP_NAME = "cronjob"
  }
}

# EventBridge rule to schedule the cron task
resource "aws_cloudwatch_event_rule" "cronjob" {
  name                = "${var.env}-cronjob-cron"
  description         = "Schedule for cronjob ECS task"
  schedule_expression = "rate(1 hour)"
  is_enabled          = false # disabled by default for testing
}

# IAM role for EventBridge to run ECS tasks
resource "aws_iam_role" "cronjob_events" {
  name = "${var.env}-cronjob-events"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "events.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "cronjob_events" {
  name = "${var.env}-cronjob-events"
  role = aws_iam_role.cronjob_events.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecs:RunTask"
        ]
        Resource = ["*"]
      },
      {
        Effect = "Allow"
        Action = [
          "iam:PassRole"
        ]
        Resource = ["*"]
      }
    ]
  })
}

# EventBridge target pointing to the ECS task
resource "aws_cloudwatch_event_target" "cronjob" {
  rule     = aws_cloudwatch_event_rule.cronjob.name
  arn      = module.ecs.ecs_cluster_arn
  role_arn = aws_iam_role.cronjob_events.arn

  ecs_target {
    task_definition_arn = module.cronjob.task_definition_arn
    task_count          = 1
    launch_type         = "FARGATE"

    network_configuration {
      subnets = local.private_subnets
    }
  }
}
