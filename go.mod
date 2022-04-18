module github.com/hazelops/ize

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.3.4
	github.com/Masterminds/semver v1.5.0
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/aws/aws-sdk-go v1.43.41
	github.com/bgentry/speakeasy v0.1.0
	github.com/briandowns/spinner v1.18.1
	github.com/containerd/console v1.0.3
	github.com/docker/cli v20.10.14+incompatible
	github.com/docker/distribution v2.8.1+incompatible
	github.com/docker/docker v20.10.14+incompatible
	github.com/fatih/color v1.13.0
	github.com/hashicorp/go-version v1.4.0
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
	github.com/morikuni/aec v1.0.0
	github.com/oklog/ulid v1.3.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pterm/pterm v0.12.41
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/zclconf/go-cty v1.10.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
	gopkg.in/ini.v1 v1.66.4
)

replace github.com/spf13/pflag => github.com/cornfeedhobo/pflag v1.1.0
