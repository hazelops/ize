package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type Project struct {
	TerraformVersion string `mapstructure:"terraform_version,omitempty"`
	AwsRegion        string `mapstructure:"aws_region,omitempty"`
	AwsProfile       string `mapstructure:"aws_profile,omitempty"`
	Namespace        string `mapstructure:",omitempty"`
	Env              string `mapstructure:",omitempty"`
	LogLevel         string `mapstructure:"log_level,omitempty"`
	PlainText        bool   `mapstructure:"plain_text,omitempty"`
	CustomPrompt     bool   `mapstructure:"custom_prompt,omitempty"`
	PreferRuntime    string `mapstructure:"prefer_runtime,omitempty"`
	Tag              string `mapstructure:",omitempty"`
	DockerRegistry   string `mapstructure:"docker_registry,omitempty"`

	Home      string `mapstructure:",omitempty"`
	RootDir   string `mapstructure:"root_dir,omitempty"`
	InfraDir  string `mapstructure:"infra_dir,omitempty"`
	EnvDir    string `mapstructure:"env_dir,omitempty"`
	AppsPath  string `mapstructure:"apps_path,omitempty"`
	TFLog     string `mapstructure:"tf_log,omitempty"`
	TFLogPath string `mapstructure:"tf_log_path,omitempty"`

	Session *session.Session

	Tunnel     *Tunnel                `mapstructure:",omitempty"`
	Terraform  map[string]*Terraform  `mapstructure:",omitempty"`
	Ecs        map[string]*Ecs        `mapstructure:",omitempty"`
	Serverless map[string]*Serverless `mapstructure:",omitempty"`
	Alias      map[string]*Alias      `mapstructure:",omitempty"`
}

func (p *Project) GetApps() map[string]*interface{} {
	apps := map[string]*interface{}{}

	for name, body := range p.Ecs {
		var v interface{}
		v = map[string]interface{}{
			"depends_on": body.DependsOn,
		}
		apps[name] = &v
	}

	for name, body := range p.Serverless {
		var v interface{}
		v = map[string]interface{}{
			"depends_on": body.DependsOn,
		}
		apps[name] = &v
	}

	for name, body := range p.Alias {
		var v interface{}
		v = map[string]interface{}{
			"depends_on": body.DependsOn,
		}
		apps[name] = &v
	}

	return apps
}
