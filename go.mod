module github.com/hazelops/ize

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.3.3
	github.com/Masterminds/semver v1.5.0
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/Microsoft/hcsshim v0.9.2 // indirect
	github.com/aws/aws-sdk-go v1.43.22
	github.com/bgentry/speakeasy v0.1.0
	github.com/briandowns/spinner v1.18.1
	github.com/containerd/cgroups v1.0.3 // indirect
	github.com/containerd/console v1.0.3
	github.com/containerd/containerd v1.6.1 // indirect
	github.com/docker/cli v20.10.13+incompatible
	github.com/docker/distribution v2.8.1+incompatible
	github.com/docker/docker v20.10.13+incompatible
	github.com/fatih/color v1.13.0
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gookit/color v1.5.0 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/hashicorp/hcl/v2 v2.11.1
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lab47/vterm v0.0.0-20211107042118-80c3d2849f9c
	github.com/mattn/go-isatty v0.0.14
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/go-glint v0.0.0-20210722152315-6515ceb4a127
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.3
	github.com/moby/buildkit v0.10.0 // indirect
	github.com/moby/sys/mount v0.3.1 // indirect
	github.com/moby/sys/mountinfo v0.6.0 // indirect
	github.com/moby/sys/symlink v0.2.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0
	github.com/oklog/ulid v1.3.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.1.0 // indirect
	github.com/pterm/pterm v0.12.39
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.1
	github.com/zclconf/go-cty v1.10.0
	go.opentelemetry.io/contrib v1.5.0 // indirect
	go.opentelemetry.io/otel/oteltest v1.0.0-RC3 // indirect
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
	google.golang.org/genproto v0.0.0-20220322021311-435b647f9ef2 // indirect
	google.golang.org/grpc v1.45.0 // indirect
	gopkg.in/ini.v1 v1.66.4
)

replace github.com/spf13/pflag => github.com/cornfeedhobo/pflag v1.1.0
