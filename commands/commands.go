package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type baseCmd struct {
	cmd *cobra.Command
}

type baseBuilderCmd struct {
	*baseCmd
	*commandsBuilder
}

type izeCmd struct {
	*baseBuilderCmd

	//Need to get state app who build
	c *commandeer
}

type commandsBuilder struct {
	izeBuilderCommon

	commands []cmder
}

func newCommandBuilder() *commandsBuilder {
	return &commandsBuilder{}
}

func (b *commandsBuilder) newBuilderCmd(cmd *cobra.Command) *baseBuilderCmd {
	bcmd := &baseBuilderCmd{commandsBuilder: b, baseCmd: &baseCmd{cmd: cmd}}
	bcmd.izeBuilderCommon.handleFlags(cmd)
	return bcmd
}

func (b *commandsBuilder) addCommands(commands ...cmder) *commandsBuilder {
	b.commands = append(b.commands, commands...)
	return b
}

func (b *commandsBuilder) addAll() *commandsBuilder {
	b.addCommands(b.newTerraformCmd())

	return b
}

func (b *commandsBuilder) newBuilderBasicCdm(cmd *cobra.Command) *baseBuilderCmd {
	bcmd := &baseBuilderCmd{baseCmd: &baseCmd{cmd: cmd}, commandsBuilder: b}
	bcmd.izeBuilderCommon.handleCommonBuilderFlags(cmd)
	return bcmd
}

func (b *commandsBuilder) newIzeCmd() *izeCmd {
	cc := &izeCmd{}

	cc.baseBuilderCmd = b.newBuilderCmd(&cobra.Command{
		Use:   "ize",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	})

	if cc.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cc.cfgFile)
	} else {
		//// Find home directory.
		//home, err := os.UserHomeDir()
		//cobra.CheckErr(err)

		// Search config in home directory with name ".ize" (without extension).
		viper.AddConfigPath(".")

		viper.SetConfigName("ize")
		viper.SetConfigType("yaml")
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AutomaticEnv() // read in environment variables that match

	//TODO ensure values of the variables are checked for nil before passing down to docker.

	// Global
	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.infra", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.Get("ENV")))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	viper.SetDefault("TF_LOG", fmt.Sprintf(""))
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt", viper.Get("ENV_DIR")))
	viper.SetDefault("TERRAFORM_VERSION", fmt.Sprintf("0.12.29"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	cc.cmd.PersistentFlags().StringVarP(&cc.cfgFile, "config", "c", "", "config file (default is $HOME/.ize.yaml)")
	cc.cmd.PersistentFlags().BoolVarP(&cc.logging, "", "v", false, "enable debug message")

	return cc
}

func addCommands(root *cobra.Command, commands ...cmder) {
	for _, command := range commands {
		cmd := command.getCommand()
		if cmd == nil {
			continue
		}
		root.AddCommand(cmd)
	}
}

func (c *baseCmd) getCommand() *cobra.Command {
	return c.cmd
}

func (b *commandsBuilder) build() *izeCmd {
	i := b.newIzeCmd()
	addCommands(i.getCommand(), b.commands...)
	return i
}

type izeBuilderCommon struct {
	cfgFile string
	logging bool
}

func (cc *izeBuilderCommon) handleCommonBuilderFlags(cmd *cobra.Command) {
}

func (cc *izeBuilderCommon) handleFlags(cmd *cobra.Command) {
	cc.handleCommonBuilderFlags(cmd)
}
