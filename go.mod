module github.com/hazelops/ize

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.3.2
	github.com/Masterminds/semver v1.5.0
	github.com/aws/aws-sdk-go v1.42.11
	github.com/containerd/containerd v1.5.8 // indirect
	github.com/docker/docker v20.10.11+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/elliotchance/sshtunnel v1.3.1
	github.com/hashicorp/hcl/v2 v2.10.1
	github.com/mitchellh/mapstructure v1.4.2
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/pterm/pterm v0.12.33
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/zclconf/go-cty v1.10.0
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871 // indirect
	golang.org/x/net v0.0.0-20211118161319-6a13c67c3ce4 // indirect
	golang.org/x/sys v0.0.0-20211123173158-ef496fb156ab // indirect
	gopkg.in/ini.v1 v1.64.0
)

replace github.com/spf13/pflag => github.com/cornfeedhobo/pflag v1.1.0
