package commands

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

var (
	ll string
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
	b.addCommands(b.newTerraformCmd())

	return b
}

func (b *commandsBuilder) newBuilderBasicCdm(cmd *cobra.Command) *baseBuilderCmd {
	bcmd := &baseBuilderCmd{baseCmd: &baseCmd{cmd: cmd}, commandsBuilder: b}
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

	cc.cmd.PersistentFlags().StringVarP(&ll, "log-level", "l", "warn", "enable debug message")
	cc.cmd.PersistentFlags().StringVarP(&cc.cfgFile, "config-file", "c", "", "set config file name")

	var logLevel zapcore.Level

	// TODO: Fix
	switch ll {
	case "info":
		logLevel = zapcore.InfoLevel

	case "debug":
		logLevel = zapcore.DebugLevel
	default:
		logLevel = zapcore.WarnLevel
	}

	cc.log = logger.NewSugaredLogger(logLevel)

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
	cfg     *config.Config
	log     logger.StandartLogger
}

func (cc *izeBuilderCommon) Init() error {
	cfg, err := cc.initConfig(cc.cfgFile)
	if err != nil {
		return err
	}

	cc.cfg = cfg

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
	}

	var app interface{}

	for _, s := range cfg.Service {
		switch serviceType := s.Type; serviceType {
		case "ecs":
			ecsCfg := Ecs{}
			if s.Body != nil {
				if diag := gohcl.DecodeBody(s.Body, &hcl.EvalContext{}, &ecsCfg); diag.HasErrors() {
					return fmt.Errorf("error: %w", diag)
				}
			}
			app = ecsCfg
		}

	}

	fmt.Println(app)
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AutomaticEnv() // read in environment variables that match

	//TODO ensure values of the variables are checked for nil before passing down to docker.

	// Global
	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.infra", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, cc.cfg.Env))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	viper.SetDefault("TF_LOG", fmt.Sprintf(""))
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt", viper.Get("ENV_DIR")))

	return nil
}

type Ecs struct {
	TerraformStateBucketName string `hcl:"terraform_state_bucket_name"`
}
