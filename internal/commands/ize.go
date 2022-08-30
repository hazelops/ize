package commands

import (
	"bytes"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/version"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

var izeDescTpl = templates.LongDesc(`
{{ .Message}}
{{ .Docs}}
{{ .Version}}

Opinionated tool for infrastructure and code.

This tool is designed as a simple wrapper around popular tools, 
so they can be easily integrated in one infra: terraform, 
ECS deployment, serverless, and others.

It combines infra, build and deploy workflows in one 
and is too simple to be considered sophisticated. 
So let's not do it but rather embrace the simplicity and minimalism.
`)

func newRootCmd(project *config.Project) *cobra.Command {
	var izeLongDesc bytes.Buffer
	err := template.Must(template.New("desc").Parse(izeDescTpl)).Execute(&izeLongDesc, struct {
		Message string
		Docs    string
		Version string
	}{
		Message: pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
		Docs:    pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh/docs"),
		Version: pterm.Sprintf("%s %s", pterm.Green("Version:"), version.FullVersionNumber()),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:              "ize",
		TraverseChildren: true,
		SilenceErrors:    true,
		Long:             izeLongDesc.String(),
		Version:          version.FullVersionNumber(),
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.PersistentFlags().StringP("log-level", "l", "", "set log level. Possible levels: info, debug, trace, panic, warn, error, fatal(default)")
	cmd.PersistentFlags().Bool("plain-text", false, "enable plain text")
	cmd.PersistentFlags().StringP("config-file", "c", "", "set config file name")
	cmd.PersistentFlags().StringP("env", "e", "", "(required) set environment name (overrides value set in IZE_ENV / ENV if any of them are set)")
	cmd.PersistentFlags().StringP("aws-profile", "p", "", "(required) set AWS profile (overrides value in ize.toml and IZE_AWS_PROFILE / AWS_PROFILE if any of them are set)")
	cmd.PersistentFlags().StringP("aws-region", "r", "", "(required) set AWS region (overrides value in ize.toml and IZE_AWS_REGION / AWS_REGION if any of them are set)")
	cmd.PersistentFlags().StringP("namespace", "n", "", "(required) set namespace (overrides value in ize.toml and IZE_NAMESPACE / NAMESPACE if any of them are set)")
	cmd.PersistentFlags().String("terraform-version", "", "set terraform-version")
	cmd.PersistentFlags().String("prefer-runtime", "native", "set prefer runtime (native or docker)")
	cmd.Flags().StringP("tag", "t", "", "set tag")

	cmd.AddCommand(
		NewCmdBuild(project),
		NewCmdDeploy(project),
		NewCmdDown(project),
		NewCmdConsole(project),
		NewCmdTerraform(project),
		NewCmdSecrets(project),
		NewCmdInit(),
		NewCmdTunnel(project),
		NewCmdExec(project),
		NewCmdConfig(),
		NewCmdLogs(project),
		NewDebugCmd(project),
		NewCmdGen(project),
		NewCmdPush(project),
		NewCmdUp(project),
		NewValidateCmd(),
		NewVersionCmd())

	return cmd
}

func Execute() {
	cfg := new(config.Project)
	cmd := newRootCmd(cfg)

	cobra.OnInitialize(func() {
		config.InitConfig()

		cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			if len(f.Value.String()) != 0 {
				_ = viper.BindPFlag(strings.ReplaceAll(f.Name, "-", "_"), cmd.PersistentFlags().Lookup(f.Name))
			}
		})

		if !(slices.Contains(os.Args, "aws-profile") ||
			slices.Contains(os.Args, "doc") ||
			slices.Contains(os.Args, "completion") ||
			slices.Contains(os.Args, "version") ||
			slices.Contains(os.Args, "init") ||
			slices.Contains(os.Args, "validate") ||
			slices.Contains(os.Args, "config")) {
			err := cfg.GetConfig()
			if err != nil {
				pterm.Error.Println(err)
				os.Exit(1)
			}
		}
	})

	if err := cmd.Execute(); err != nil {
		fmt.Println()
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

func init() {
	initLogger()
	customizeDefaultPtermPrefix()
}

func initLogger() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		PadLevelText:     true,
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
}

func customizeDefaultPtermPrefix() {
	pterm.Info.Prefix = pterm.Prefix{
		Text:  "ℹ",
		Style: pterm.NewStyle(pterm.FgBlue),
	}

	pterm.Success.Prefix = pterm.Prefix{
		Text:  "✓",
		Style: pterm.NewStyle(pterm.FgGreen),
	}

	pterm.Error.Prefix = pterm.Prefix{
		Text:  "✗",
		Style: pterm.NewStyle(pterm.FgRed),
	}

	pterm.Warning.Prefix = pterm.Prefix{
		Text:  "⚠",
		Style: pterm.NewStyle(pterm.FgYellow),
	}

	pterm.DefaultSpinner.Sequence = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
}
