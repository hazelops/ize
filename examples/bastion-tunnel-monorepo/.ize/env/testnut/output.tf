output "bastion_instance_id" {
  description = "Bastion EC2 instance ID"
  value       = module.bastion.instance_id
}

output "cmd" {
  description = "Map of useful commands"
  value = merge({}, {
    tunnel = module.bastion.cmd
  })
}

output "ssh_forward_config" {
  value = module.bastion.ssh_config
}

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
