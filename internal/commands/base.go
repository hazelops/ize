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
	"go.uber.org/zap/zapcore"
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
	)

	return b
}

func (b *commandsBuilder) newBuilderBasicCdm(cmd *cobra.Command) *baseBuilderCmd {
	bcmd := &baseBuilderCmd{baseCmd: &baseCmd{cmd: cmd}, commandsBuilder: b}
	return bcmd
}

func (b *commandsBuilder) newIzeCmd() *izeCmd {
	cc := &izeCmd{}

	cc.baseBuilderCmd = b.newBuilderCmd(&cobra.Command{
		Use:     "ize",
		Version: GetVersionNumber(),
		Short:   "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cc.cmd.PersistentFlags().StringVarP(&cc.ll, "log-level", "l", "", "enable debug message")
	cc.cmd.PersistentFlags().StringVarP(&cc.cfgFile, "config-file", "c", "", "set config file name")
			cc.cmd.PersistentFlags().Parse(args)

	var logLevel zapcore.Level

	switch cc.ll {
	case "info":
		logLevel = zapcore.InfoLevel
	case "debug":
		logLevel = zapcore.DebugLevel
	default:
		logLevel = zapcore.WarnLevel
	}

	cc.log = logger.NewSugaredLogger(logLevel)
		},
	})

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
	cfgFile string
	ll      string

	config *config.Config
	log    logger.StandartLogger
}

func (cc *izeBuilderCommon) Init() error {
	config, err := cc.initConfig(cc.cfgFile)
	if err != nil {
		return err
	}

	cc.config = config

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
	}

	// for _, s := range config.Service {
	// 	switch serviceType := s.Type; serviceType {
	// 	case "ecs":
	// 		ecsCfg := Ecs{}
	// 		if s.Body != nil {
	// 			if diag := gohcl.DecodeBody(s.Body, &hcl.EvalContext{}, &ecsCfg); diag.HasErrors() {
	// 				return fmt.Errorf("error: %w", diag)
	// 			}
	// 		}
	// 		app = ecsCfg
	// 	}

	// }

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

	//Check Docker and SSM Agent
	_, err = CheckCommand("docker", []string{"info"})
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

type Ecs struct {
	TerraformStateBucketName string `hcl:"terraform_state_bucket_name"`
}
