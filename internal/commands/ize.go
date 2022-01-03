package commands

import (
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/commands/console"
	"github.com/hazelops/ize/internal/commands/deploy"
	"github.com/hazelops/ize/internal/commands/env"
	"github.com/hazelops/ize/internal/commands/mfa"
	"github.com/hazelops/ize/internal/commands/terraform"
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

var (
	rootCmd = &cobra.Command{
		Use: "ize",
		Long: fmt.Sprintf("%s\n%s\n%s",
			pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
			pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh"),
			pterm.Sprintf("%s %s", pterm.Green("Version:"), Version),
		),
		Version:          Version,
		TraverseChildren: true,
	}
)

func newApp() (*cobra.Command, error) {
	rootCmd = &cobra.Command{
		Use: "ize",
		Long: fmt.Sprintf("%s\n%s\n%s",
			pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
			pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh"),
			pterm.Sprintf("%s %s", pterm.Green("Version:"), Version),
		),
		Version:          Version,
		TraverseChildren: true,
	}

	rootCmd.AddCommand(
func Execute(args []string) Response {
	izeCmd := newCommandBuilder().addAll().build()
	cmd := izeCmd.getCommand()
	cmd.SetArgs(args)
		NewVersionCmd(),
	)

	rootCmd.PersistentFlags().StringP("log-level", "l", "", "enable debug messages")
	rootCmd.PersistentFlags().StringP("config-file", "c", "", "set config file name")

	rootCmd.Flags().StringP("env", "e", "", "set enviroment name")
	rootCmd.Flags().StringP("aws-profile", "p", "", "set AWS profile")
	rootCmd.Flags().StringP("aws-region", "r", "", "set AWS region")
	rootCmd.Flags().StringP("namespace", "n", "", "set namespace")
	rootCmd.Flags().StringP("tag", "t", "", "set tag")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())

	return rootCmd, nil
}
