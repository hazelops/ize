package config

type App struct {
	Type string `mapstructure:",omitempty"`
	Path string `mapstructure:",omitempty"`
}
