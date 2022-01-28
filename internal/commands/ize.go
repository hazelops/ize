package commands

import (
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/commands/console"
	"github.com/hazelops/ize/internal/commands/deploy"
	"github.com/hazelops/ize/internal/commands/env"
	"github.com/hazelops/ize/internal/commands/initialize"
	"github.com/hazelops/ize/internal/commands/mfa"
	"github.com/hazelops/ize/internal/commands/secrets"
	"github.com/hazelops/ize/internal/commands/terraform"
	"github.com/hazelops/ize/internal/commands/tunnel"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
		Use: "ize",
		Long: fmt.Sprintf("%s\n%s\n%s",
			pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
			pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh/docs"),
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
		secrets.NewCmdSecrets(),
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
