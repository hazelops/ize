package config

type Ecs struct {
	Name                   string   `mapstructure:",omitempty"`
	Unsafe                 bool     `mapstructure:",omitempty"`
	Path                   string   `mapstructure:",omitempty"`
	Image                  string   `mapstructure:",omitempty"`
	Cluster                string   `mapstructure:",omitempty"`
	TaskDefinitionRevision string   `mapstructure:"task_definition_revision"`
	Timeout                int      `mapstructure:",omitempty"`
	DependsOn              []string `mapstructure:"depends_on,omitempty"`
}

type Serverless struct {
	Name                    string   `mapstructure:",omitempty"`
	File                    string   `mapstructure:",omitempty"`
	NodeVersion             string   `mapstructure:"node_version"`
	Env                     []string `mapstructure:",omitempty"`
	Path                    string   `mapstructure:",omitempty"`
	SLSNodeModuleCacheMount string   `mapstructure:",omitempty"`
	CreateDomain            bool     `mapstructure:"create_domain"`
	DependsOn               []string `mapstructure:"depends_on,omitempty"`
}

type Alias struct {
	Name      string
	DependsOn string `mapstructure:"depends_on"`
}
