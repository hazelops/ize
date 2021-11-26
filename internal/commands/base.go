package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/logger"
	"github.com/pterm/pterm"
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
	return bcmd
}

func (b *commandsBuilder) addCommands(commands ...cmder) *commandsBuilder {
	b.commands = append(b.commands, commands...)
	return b
}

func (b *commandsBuilder) addAll() *commandsBuilder {
	b.addCommands(
		b.newTerraformCmd(),
		b.newConfigCmd(),
		b.newEnvCmd(),
		b.newTunnelCmd(),
		b.newMfaCmd(),
		b.newSSHCmd(),
		b.newInitCmd(),
		b.newDeployCmd(),
	)

	return b
}

func (b *commandsBuilder) newBuilderBasicCdm(cmd *cobra.Command) *baseBuilderCmd {
	bcmd := &baseBuilderCmd{baseCmd: &baseCmd{cmd: cmd}, commandsBuilder: b}
	return bcmd
}

var (
	rootCmd = &cobra.Command{
		Use: "ize",
		Long: fmt.Sprintf("%s\n%s\n%s",
			pterm.White(pterm.Bold.Sprint("Welcome to IZE")),
			pterm.Sprintf("%s %s", pterm.Blue("Docs:"), "https://ize.sh"),
			pterm.Sprintf("%s %s", pterm.Green("Version:"), Version),
		),
		TraverseChildren: true,
	}
)

func init() {
	rootCmd.PersistentFlags().StringP("log-level", "l", "", "enable debug messages")
	rootCmd.PersistentFlags().StringP("config-file", "c", "", "set config file name")

	rootCmd.Flags().StringP("env", "e", "", "set environment name")
	rootCmd.Flags().StringP("aws-profile", "p", "", "set AWS profile")
	rootCmd.Flags().StringP("aws-region", "r", "", "set AWS region")

	rootCmd.Flags().StringP("namespace", "n", "", "set namespace")

	//Bind viper key to a flag (required for flags/parameters that are more than 1 word)
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("config_file", rootCmd.PersistentFlags().Lookup("config-file"))
	viper.BindPFlag("aws_profile", rootCmd.Flags().Lookup("aws-profile"))
	viper.BindPFlag("aws_region", rootCmd.Flags().Lookup("aws-region"))

	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func (b *commandsBuilder) newIzeCmd() *izeCmd {
	cc := &izeCmd{}

	cc.baseBuilderCmd = b.newBuilderCmd(rootCmd)

	cc.baseCmd.cmd.SilenceErrors = true
	cc.cmd.SilenceUsage = true

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
	config *config.Config
	log    logger.StandartLogger
}

func (cc *izeBuilderCommon) Init() error {
	viper.SetEnvPrefix("IZE")
	viper.AutomaticEnv()

	config, err := cc.initConfig(viper.GetString("config_file"))
	if err != nil {
		return err
	}

	cc.config = config

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
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, cc.config.Env))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	viper.SetDefault("TF_LOG", fmt.Sprintf(""))
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt", viper.Get("ENV_DIR")))

	if err = CheckRequirements(); err != nil {
		return err
	}

	return nil
}

func CheckRequirements() error {
	//Check Docker and SSM Agent
	_, err := CheckCommand("docker", []string{"info"})
	if err != nil {
		return errors.New("docker is not running or is not installed (visit https://www.docker.com/get-started)")
	}

	_, err = CheckCommand("session-manager-plugin", []string{})
	if err != nil {
		pterm.Warning.Println("SSM Agent plugin is not installed. Trying to install SSM Agent plugin")

		var pyVersion string

		pyVersion, err = CheckCommand("python3", []string{"--version"})
		if err != nil {
			pyVersion, err = CheckCommand("python", []string{"--version"})
			if err != nil {
				return errors.New("python is not installed")
			}

			c, err := semver.NewConstraint("<= 2.6.5")
			if err != nil {
				return err
			}

			v, err := semver.NewVersion(strings.TrimSpace(strings.Split(pyVersion, " ")[1]))
			if err != nil {
				return err
			}

			if c.Check(v) {
				return fmt.Errorf("python version %s below required %s", v.String(), "2.6.5")
			}
			return errors.New("python is not installed")
		}

		c, err := semver.NewConstraint("<= 3.3.0")
		if err != nil {
			return err
		}

		v, err := semver.NewVersion(strings.TrimSpace(strings.Split(pyVersion, " ")[1]))
		if err != nil {
			return err
		}

		if c.Check(v) {
			return fmt.Errorf("python version %s below required %s", v.String(), "3.3.0")
		}

		pterm.DefaultSection.Println("Installing SSM Agent plugin")

		err = DownloadSSMAgentPlugin()
		if err != nil {
			return fmt.Errorf("download SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Downloading SSM Agent plugin")

		err = InstallSSMAgent()
		if err != nil {
			return fmt.Errorf("install SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Installing SSM Agent plugin")

		err = CleanupSSMAgent()
		if err != nil {
			return fmt.Errorf("cleanup SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Cleanup Session Manager plugin installation package")

		_, err = CheckCommand("session-manager-plugin", []string{})
		if err != nil {
			return fmt.Errorf("check SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}
	}

	return nil
}
