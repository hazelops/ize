module github.com/hazelops/ize

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.3.2
	github.com/Masterminds/semver v1.5.0
	github.com/aws/aws-sdk-go v1.40.56
	github.com/containerd/containerd v1.5.4 // indirect
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/elliotchance/sshtunnel v1.3.1
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/hashicorp/hcl/v2 v2.10.1
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/pterm/pterm v0.12.31
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/zclconf/go-cty v1.8.0
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871 // indirect
	golang.org/x/net v0.0.0-20211118161319-6a13c67c3ce4 // indirect
	golang.org/x/sys v0.0.0-20211117180635-dee7805ff2e1 // indirect
	gopkg.in/ini.v1 v1.62.0
)

replace github.com/spf13/pflag => github.com/cornfeedhobo/pflag v1.1.0
