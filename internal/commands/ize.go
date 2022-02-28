package commands

import (
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/commands/config"
	"github.com/hazelops/ize/internal/commands/console"
	"github.com/hazelops/ize/internal/commands/deploy"
	"github.com/hazelops/ize/internal/commands/env"
	"github.com/hazelops/ize/internal/commands/exec"
	"github.com/hazelops/ize/internal/commands/initialize"
	"github.com/hazelops/ize/internal/commands/logs"
	"github.com/hazelops/ize/internal/commands/mfa"
	"github.com/hazelops/ize/internal/commands/secrets"
	"github.com/hazelops/ize/internal/commands/terraform"
	"github.com/hazelops/ize/internal/commands/tunnel"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
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

type Response struct {
	Err error

	Cmd *cobra.Command
}

func Execute(args []string) error {
	go CheckLatestRealese()

	app, err := newApp()
	if err != nil {
		return err
	}
	return app.Execute()
}

func newApp() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:              "ize",
		TraverseChildren: true,
		SilenceErrors:    true,
		Long:             deployIzeDesc,
		Version:          Version,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n%s\n%s\n\n",
				pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
				pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh/docs"),
				pterm.Sprintf("%s %s", pterm.Green("Version:"), Version),
			)
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		deploy.NewCmdDeploy(),
		console.NewCmdConsole(),
		env.NewCmdEnv(),
		mfa.NewCmdMfa(),
		terraform.NewCmdTerraform(),
		secrets.NewCmdSecrets(),
		initialize.NewCmdInit(),
		tunnel.NewCmdTunnel(),
		exec.NewCmdExec(),
		config.NewCmdConfig(),
		logs.NewCmdLogs(),
		NewGendocCmd(),
		NewVersionCmd(),
	)

	rootCmd.PersistentFlags().StringP("log-level", "l", "", "enable debug messages")
	rootCmd.PersistentFlags().StringP("config-file", "c", "", "set config file name")
	rootCmd.PersistentFlags().StringP("env", "e", "", "(required) set environment name (overrides value set in ENV / IZE_ENV if any of them are set)")
	rootCmd.PersistentFlags().StringP("aws-profile", "p", "", "(required) set AWS profile (overrides value in ize.toml and IZE_AWS_PROFILE / AWS_PROFILE if any of them are set)")
	rootCmd.PersistentFlags().StringP("aws-region", "r", "", "(required) set AWS region (overrides value in ize.toml and IZE_AWS_REGION / AWS_REGION if any of them are set)")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "set namespace")
	rootCmd.PersistentFlags().String("terraform-version", "", "set terraform-version")

	rootCmd.Flags().StringP("tag", "t", "", "set tag")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	BindFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.PersistentFlags())

	return rootCmd, nil
}

func BindFlags(flags *pflag.FlagSet) {
	replacer := strings.NewReplacer("-", "_")

	flags.VisitAll(func(flag *pflag.Flag) {
		if err := viper.BindPFlag(replacer.Replace(flag.Name), flag); err != nil {
			panic("unable to bind flag " + flag.Name + ": " + err.Error())
		}
	})
}
