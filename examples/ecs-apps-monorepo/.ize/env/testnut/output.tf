output "bastion_instance_id" {
  description = "Bastion EC2 instance ID"
  value       = module.openvpn_instance.instance_id
}

output "cmd" {
  description = "Map of useful commands"
  value = merge({},{
    tunnel = module.openvpn_instance.cmd
  })
}

output "ssh_forward_config" {
  value = module.openvpn_instance.ssh_config
}

output "vpc_public_subnets" {
  value = module.vpc.public_subnets
}

output "vpc_private_subnets" {
  value = module.vpc.private_subnets
}
