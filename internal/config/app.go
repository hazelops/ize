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
	DependsOn              []string `mapstructure:"depends_on,omitempty"`
}

type Serverless struct {
	Name                    string   `mapstructure:",omitempty"`
	File                    string   `mapstructure:",omitempty"`
	NodeVersion             string   `mapstructure:"node_version"`
	Path                    string   `mapstructure:",omitempty"`
	SLSNodeModuleCacheMount string   `mapstructure:",omitempty"`
	CreateDomain            bool     `mapstructure:"create_domain"`
	Env                     []string `mapstructure:",omitempty"`
	DependsOn               []string `mapstructure:"depends_on,omitempty"`
}

type Alias struct {
	Name      string
	DependsOn string `mapstructure:"depends_on"`
}
