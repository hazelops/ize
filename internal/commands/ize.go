package commands

import (
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/commands/console"
	"github.com/hazelops/ize/internal/commands/deploy"
	"github.com/hazelops/ize/internal/commands/env"
	"github.com/hazelops/ize/internal/commands/initialize"
	"github.com/hazelops/ize/internal/commands/mfa"
	"github.com/hazelops/ize/internal/commands/secret"
	"github.com/hazelops/ize/internal/commands/terraform"
	"github.com/hazelops/ize/internal/commands/tunnel"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Response struct {
	Err error

	Cmd *cobra.Command
}

func Execute(args []string) error {
	app, err := newApp()
	if err != nil {
		return err
	}
	return app.Execute()
}

func newApp() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use: "ize",
		Long: fmt.Sprintf("%s\n%s\n%s",
			pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
			pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh"),
			pterm.Sprintf("%s %s", pterm.Green("Version:"), Version),
		),
		TraverseChildren: true,
		SilenceUsage:     true,
		SilenceErrors:    true,
	}

	rootCmd.AddCommand(
		deploy.NewCmdDeploy(),
		console.NewCmdConsole(),
		env.NewCmdEnv(),
		mfa.NewCmdMfa(),
		terraform.NewCmdTerraform(),
		secret.NewCmdSecret(),
		initialize.NewCmdInit(),
		tunnel.NewCmdTunnel(),
		NewGendocCmd(),
		NewVersionCmd(),
	)

	rootCmd.PersistentFlags().StringP("log-level", "l", "", "enable debug messages")
	rootCmd.PersistentFlags().StringP("config-file", "c", "", "set config file name")
	rootCmd.PersistentFlags().StringP("env", "e", "", "set enviroment name")
	rootCmd.PersistentFlags().StringP("aws-profile", "p", "", "set AWS profile")
	rootCmd.PersistentFlags().StringP("aws-region", "r", "", "set AWS region")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "set namespace")

	rootCmd.Flags().StringP("tag", "t", "", "set tag")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())

	return rootCmd, nil
}
