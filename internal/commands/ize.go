package commands

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/hazelops/ize/internal/commands/configure"
	"github.com/hazelops/ize/internal/commands/console"
	"github.com/hazelops/ize/internal/commands/deploy"
	"github.com/hazelops/ize/internal/commands/destroy"
	"github.com/hazelops/ize/internal/commands/env"
	"github.com/hazelops/ize/internal/commands/exec"
	"github.com/hazelops/ize/internal/commands/initialize"
	"github.com/hazelops/ize/internal/commands/logs"
	"github.com/hazelops/ize/internal/commands/mfa"
	"github.com/hazelops/ize/internal/commands/secrets"
	"github.com/hazelops/ize/internal/commands/status"
	"github.com/hazelops/ize/internal/commands/terraform"
	"github.com/hazelops/ize/internal/commands/tunnel"
	cfg "github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/version"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var deployIzeDesc = templates.LongDesc(`
Opinionated tool for infrastructure and code.

This tool is designed as a simple wrapper around popular tools, 
so they can be easily integrated in one infra: terraform, 
ECS deployment, serverless, and others.

It combines infra, build and deploy workflows in one 
and is too simple to be considered sophisticated. 
So let's not do it but rather embrace the simplicity and minimalism.
`)

var (
	rootCmd = &cobra.Command{
		Use:              "ize",
		TraverseChildren: true,
		SilenceErrors:    true,
		Long:             deployIzeDesc,
		Version:          version.FullVersionNumber(),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n%s\n%s\n\n",
				pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
				pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh/docs"),
				pterm.Sprintf("%s %s", pterm.Green("Version:"), version.FullVersionNumber()),
			)
			cmd.Help()
		},
	}
)

func Execute(args []string) {
	go version.CheckLatestRealese()

	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Println(err)
	}
}

func init() {
	initLogrus()
	customizeDefaultPtermPrefix()

	rootCmd.PersistentFlags().StringP("log-level", "l", "", "enable debug messages")
	rootCmd.PersistentFlags().Bool("plain-text", false, "enable plain text")
	rootCmd.PersistentFlags().StringP("config-file", "c", "", "set config file name")
	rootCmd.PersistentFlags().StringP("env", "e", "", "(required) set environment name (overrides value set in IZE_ENV / ENV if any of them are set)")
	rootCmd.PersistentFlags().StringP("aws-profile", "p", "", "(required) set AWS profile (overrides value in ize.toml and IZE_AWS_PROFILE / AWS_PROFILE if any of them are set)")
	rootCmd.PersistentFlags().StringP("aws-region", "r", "", "(required) set AWS region (overrides value in ize.toml and IZE_AWS_REGION / AWS_REGION if any of them are set)")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "(required) set namespace (overrides value in ize.toml and IZE_NAMESPACE / NAMESPACE if any of them are set)")
	rootCmd.PersistentFlags().String("terraform-version", "", "set terraform-version")
	rootCmd.PersistentFlags().String("prefer-runtime", "native", "set prefer runtime (native or docker)")
	rootCmd.Flags().StringP("tag", "t", "", "set tag")
	viper.BindPFlags(rootCmd.PersistentFlags())
	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(strings.ReplaceAll(f.Name, "_", "-"), rootCmd.PersistentFlags().Lookup(f.Name))
	})

	addCommands()

	cobra.OnInitialize(cfg.InitConfig)
}

func addCommands() {
	rootCmd.AddCommand(
		deploy.NewCmdDeploy(),
		destroy.NewCmdDestroy(),
		console.NewCmdConsole(),
		env.NewCmdEnv(),
		mfa.NewCmdMfa(),
		terraform.NewCmdTerraform(),
		secrets.NewCmdSecrets(),
		initialize.NewCmdInit(),
		tunnel.NewCmdTunnel(),
		exec.NewCmdExec(),
		configure.NewCmdConfig(),
		logs.NewCmdLogs(),
		status.NewDebugCmd(),
		NewGendocCmd(),
		NewVersionCmd(),
	)
}

func initLogrus() {
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
