package config

type Infra struct {
	Terraform Terraform `mapstructure:"infra.terraform,omitempty"`
	Tunnel    Tunnel    `mapstructure:"infra.tunnel,omitempty"`
}

type Terraform struct {
	Version           string `mapstructure:",omitempty"`
	StateBucketRegion string `mapstructure:"state_bucket_region,omitempty"`
	StateBucketName   string `mapstructure:"state_bucket_name,omitempty"`
	StateName         string `mapstructure:"state_name,omitempty"`
	RootDomainName    string `mapstructure:"root_domain_name,omitempty"`
	AwsRegion         string `mapstructure:"aws_region,omitempty"`
	AwsProfile        string `mapstructure:"aws_profile,omitempty"`
}

type Tunnel struct {
	BastionInstanceID string   `mapstructure:"bastion_instance_id,omitempty"`
	ForwardHost       []string `mapstructure:"forward_host,omitempty"`
}
