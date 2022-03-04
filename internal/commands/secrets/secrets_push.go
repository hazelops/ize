package secrets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SecretsPushOptions struct {
	Config      *config.Config
	AppName     string
	Backend     string
	FilePath    string
	SecretsPath string
	Force       bool
}

func NewSecretsPushFlags() *SecretsPushOptions {
	return &SecretsPushOptions{}
}

func NewCmdSecretsPush(ui terminal.UI) *cobra.Command {
	o := NewSecretsPushFlags()

	cmd := &cobra.Command{
		Use:   "push",
		Short: "push secrets to a key-value storage (like SSM)",
		Long:  "This command pushes secrets from a local file to a key-value storage (like SSM).",
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

			err = o.Run(ui)
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

func (o *SecretsPushOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
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

func (o *SecretsPushOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	return nil
}

func (o *SecretsPushOptions) Run(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("pushing secrets for %s...", o.AppName)
	defer func() { s.Abort() }()
	if o.Backend == "ssm" {
		err := o.push(s)
		if err != nil {
			return fmt.Errorf("can't push secrets: %w", err)
		}
	} else {
		return fmt.Errorf("backend with type %s not found or not supported", o.Backend)
	}

	s.Done()
	ui.Output("pushing secrets complete!\n", terminal.WithSuccessStyle())

	return nil
}

func (o *SecretsPushOptions) push(s terminal.Step) error {
	fmt.Fprintf(s.TermOutput(), "reading secrets from file...\n")
	values, err := getKeyValuePairs(o.FilePath)
	if err != nil {
		return err
	}

	fmt.Fprintf(s.TermOutput(), "pushing secrets to %s://%s...\n", o.Backend, o.SecretsPath)

	ssmSvc := ssm.New(o.Config.Session)

	for key, value := range values {
		name := fmt.Sprintf("%s/%s", o.SecretsPath, key)

		_, err := ssmSvc.PutParameter(&ssm.PutParameterInput{
			Name:      &name,
			Value:     aws.String(value),
			Type:      aws.String(ssm.ParameterTypeSecureString),
			Overwrite: &o.Force,
		})

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ParameterAlreadyExists":
				return fmt.Errorf("secret already exists, you can use --force to overwrite it")
			default:
				return err
			}
		}

		_, err = ssmSvc.AddTagsToResource(&ssm.AddTagsToResourceInput{
			ResourceId:   &name,
			ResourceType: aws.String("Parameter"),
			Tags: []*ssm.Tag{
				{
					Key:   aws.String("Application"),
					Value: &o.AppName,
				},
				{
					Key:   aws.String("EnvVarName"),
					Value: &key,
				},
			},
		})

		if err != nil {
			return err
		}
	}

	return err
}

func getKeyValuePairs(filePath string) (map[string]string, error) {
	if !filepath.IsAbs(filePath) {
		var err error
		wd, _ := os.Getwd()
		filePath, err = filepath.Abs(wd + "/" + filePath)
		if err != nil {
			return nil, err
		}

	}

	if _, err := os.Stat(filePath); err != nil {
		pterm.Fatal.Sprintfln("%s does not exist", filePath)
		return nil, err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var result map[string]string

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
