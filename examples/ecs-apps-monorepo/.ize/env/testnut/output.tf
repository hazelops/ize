output "vpc_public_subnets" {
  value = module.vpc.public_subnets
}

output "vpc_private_subnets" {
  value = module.vpc.private_subnets
}

output "security_groups" {
  value = module.vpc.default_security_group_id
}

output "subnets" {
  value = concat ([module.vpc.public_subnets], [module.vpc.private_subnets])
}
