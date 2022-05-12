package secrets

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SecretsEditOptions struct {
	Config   *config.Config
	AppName  string
	FilePath string
}

func NewSecretsEditFlags() *SecretsEditOptions {
	return &SecretsEditOptions{}
}

func NewCmdSecretsEdit() *cobra.Command {
	o := NewSecretsEditFlags()

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit secrets file",
		Long:  "This command open secrets file in default text editor",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.FilePath, "file", "", "file with secrets")

	return cmd
}

func (o *SecretsEditOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg
	o.AppName = cmd.Flags().Args()[0]

	if o.FilePath == "" {
		o.FilePath = fmt.Sprintf("%s/%s/%s.json", viper.GetString("ENV_DIR"), "secrets", o.AppName)
	}

	return nil
}

func (o *SecretsEditOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	return nil
}

func (o *SecretsEditOptions) Run() error {
	absPath, err := filepath.Abs(o.FilePath)
	if err != nil {
		return fmt.Errorf("can't secrets edit: %w", err)
	}

	checkSecretFolder()

	f, err := os.OpenFile(absPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("can't secrets edit: %w", err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("can't secrets edit: %w", err)
	}

	text := string(b)

	err = survey.AskOne(
		&survey.Editor{
			Message:       fmt.Sprintf("Edit %s secrets file", o.AppName),
			Default:       text,
			AppendDefault: true,
		},
		&text,
	)
	if err != nil {
		return fmt.Errorf("can't secrets edit: %w", err)
	}

	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("can't secrets edit: %w", err)
	}

	_, err = io.WriteString(f, text)
	if err != nil {
		return fmt.Errorf("can't secrets edit: %w", err)
	}

	return nil
}

func checkSecretFolder() {
	secretsFolder := filepath.Join(viper.GetString("ENV_DIR"), "secrets")
	_, err := os.Stat(secretsFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(secretsFolder, 0775)
	}
}
