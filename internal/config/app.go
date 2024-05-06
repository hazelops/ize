package config

type Ecs struct {
	Name                   string   `mapstructure:",omitempty"`
	Path                   string   `mapstructure:",omitempty"`
	Image                  string   `mapstructure:",omitempty"`
	Cluster                string   `mapstructure:",omitempty"`
	TaskDefinitionRevision string   `mapstructure:"task_definition_revision"`
	DockerRegistry         string   `mapstructure:"docker_registry,omitempty"`
	Timeout                int      `mapstructure:",omitempty"`
	Unsafe                 bool     `mapstructure:",omitempty"`
	SkipDeploy             bool     `mapstructure:"skip_deploy,omitempty"`
	Icon                   string   `mapstructure:"icon,omitempty"`
	AwsProfile             string   `mapstructure:"aws_profile,omitempty"`
	AwsRegion              string   `mapstructure:"aws_region,omitempty"`
	DependsOn              []string `mapstructure:"depends_on,omitempty"`
}

type Helm struct {
	Name           string   `mapstructure:",omitempty"`
	Path           string   `mapstructure:",omitempty"`
	Image          string   `mapstructure:",omitempty"`
	Namespace      string   `mapstructure:",omitempty"`
	HelmRelease    string   `mapstructure:"helm_release,omitempty"`
	DockerRegistry string   `mapstructure:"docker_registry,omitempty"`
	Timeout        int      `mapstructure:",omitempty"`
	SkipDeploy     bool     `mapstructure:"skip_deploy,omitempty"`
	Force          bool     `mapstructure:"force"`
	Icon           string   `mapstructure:"icon,omitempty"`
	AwsProfile     string   `mapstructure:"aws_profile,omitempty"`
	AwsRegion      string   `mapstructure:"aws_region,omitempty"`
	DependsOn      []string `mapstructure:"depends_on,omitempty"`
}

type Serverless struct {
	Name                    string   `mapstructure:",omitempty"`
	File                    string   `mapstructure:",omitempty"`
	NodeVersion             string   `mapstructure:"node_version"`
	ServerlessVersion       string   `mapstructure:"serverless_version"`
	Path                    string   `mapstructure:",omitempty"`
	SLSNodeModuleCacheMount string   `mapstructure:",omitempty"`
	CreateDomain            bool     `mapstructure:"create_domain"`
	Force                   bool     `mapstructure:"force"`
	Env                     []string `mapstructure:",omitempty"`
	Icon                    string   `mapstructure:"icon,omitempty"`
	UseYarn                 bool     `mapstructure:"use_yarn,omitempty"`
	AwsProfile              string   `mapstructure:"aws_profile,omitempty"`
	AwsRegion               string   `mapstructure:"aws_region,omitempty"`
	DependsOn               []string `mapstructure:"depends_on,omitempty"`
}

type Alias struct {
	Name      string   `mapstructure:",omitempty"`
	Icon      string   `mapstructure:"icon,omitempty"`
	DependsOn []string `mapstructure:"depends_on"`
}
