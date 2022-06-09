package secrets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SecretsPullOptions struct {
	Config      *config.Config
	AppName     string
	Backend     string
	FilePath    string
	SecretsPath string
	Force       bool
}

func NewSecretsPullFlags() *SecretsPullOptions {
	return &SecretsPullOptions{}
}

func NewCmdSecretsPull() *cobra.Command {
	o := NewSecretsPullFlags()

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull secrets to a a local file (like SSM)",
		Long:  "This command pulles secrets from a key-value storage to a local file (like SSM)",
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

	cmd.Flags().StringVar(&o.Backend, "backend", "ssm", "backend type (default=ssm)")
	cmd.Flags().StringVar(&o.FilePath, "file", "", "file with secrets")
	cmd.Flags().StringVar(&o.SecretsPath, "path", "", "path where to store secrets (/<env>/<app> by default)")
	cmd.Flags().BoolVar(&o.Force, "force", false, "allow values overwrite")

	return cmd
}

func (o *SecretsPullOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg
	o.AppName = cmd.Flags().Args()[0]

	if o.FilePath == "" {
		o.FilePath = fmt.Sprintf("%s/%s/%s.json", viper.GetString("ENV_DIR"), "secrets", o.AppName)
	}

	if o.SecretsPath == "" {
		o.SecretsPath = fmt.Sprintf("/%s/%s", o.Config.Env, o.AppName)
	}

	return nil
}

func (o *SecretsPullOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *SecretsPullOptions) Run() error {
	s, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Pulling secrets for %s...", o.AppName))
	if o.Backend == "ssm" {
		err := o.pull(s)
		if err != nil {
			return fmt.Errorf("can't pull secrets: %w", err)
		}
	} else {
		return fmt.Errorf("backend with type %s not found or not supported", o.Backend)
	}

	s.Success("Pulling secrets complete!")

	return nil
}

func (o *SecretsPullOptions) pull(s *pterm.SpinnerPrinter) error {
	s.UpdateText(fmt.Sprintf("Pulling secrets from %s://%s...", o.Backend, o.SecretsPath))

	ssmSvc := ssm.New(o.Config.Session)

	params, err := ssmSvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path: aws.String(o.SecretsPath),
	})
	if err != nil {
		return err
	}

	values := make(map[string]interface{})

	for _, param := range params.Parameters {
		values[*param.Name] = *param.Value
	}

	b, err := json.MarshalIndent(values, "", "")
	if err != nil {
		return err
	}

	if _, err := os.Stat(o.FilePath); os.IsNotExist(err) || o.Force {
		err := ioutil.WriteFile(o.FilePath, b, 0644)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("file %s already exists. Please use --force to overwrite.", o.FilePath)
	}

	return nil
}
