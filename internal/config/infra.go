package config

type Infra struct {
	Terraform Terraform `mapstructure:"infra.terraform,omitempty"`
	Tunnel    Tunnel    `mapstructure:"infra.tunnel,omitempty"`
}

type Terraform struct {
	Version             string   `mapstructure:",omitempty"`
	StateBucketRegion   string   `mapstructure:"state_bucket_region,omitempty"`
	StateBucketName     string   `mapstructure:"state_bucket_name,omitempty"`
	StateName           string   `mapstructure:"state_name,omitempty"`
	RootDomainName      string   `mapstructure:"root_domain_name,omitempty"`
	TerraformConfigFile string   `mapstructure:"terraform_config_file,omitempty"`
	AwsRegion           string   `mapstructure:"aws_region,omitempty"`
	AwsProfile          string   `mapstructure:"aws_profile,omitempty"`
	DependsOn           []string `mapstructure:"depends_on,omitempty"`
}

type Tunnel struct {
	BastionInstanceID string   `mapstructure:"bastion_instance_id,omitempty"`
	ForwardHost       []string `mapstructure:"forward_host,omitempty"`
	SSHPublicKey      string   `mapstructure:"ssh_public_key,omitempty"`
	SSHPrivateKey     string   `mapstructure:"ssh_private_key,omitempty"`
}
