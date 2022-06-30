package config

type Project struct {
	TerraformVersion string `mapstructure:"terraform_version,omitempty"`
	AwsRegion        string `mapstructure:"aws_region,omitempty"`
	AwsProfile       string `mapstructure:"aws_profile,omitempty"`
	Namespace        string `mapstructure:",omitempty"`
	Env              string `mapstructure:",omitempty"`
	LogLevel         string `mapstructure:"log_level,omitempty"`
	PlainText        bool   `mapstructure:"plain_text,omitempty"`
	PreferRuntime    string `mapstructure:"prefer_runtime,omitempty"`
	Tag              string `mapstructure:",omitempty"`

	home     string `mapstructure:",omitempty"`
	rootDir  string `mapstructure:"root_dir,omitempty"`
	infraDir string `mapstructure:"infra_dir,omitempty"`
	envDir   string `mapstructure:"env_dir,omitempty"`
	appsPath string `mapstructure:"apps_path,omitempty"`

	Infra *Infra          `mapstructure:",omitempty"`
	App   *map[string]App `mapstructure:",omitempty"`

	Terraform  *map[string]Terraform  `mapstructure:",omitempty"`
	Ecs        *map[string]Ecs        `mapstructure:",omitempty"`
	Serverless *map[string]Serverless `mapstructure:",omitempty"`
}
