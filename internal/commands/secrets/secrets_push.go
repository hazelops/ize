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
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type SecretsPushOptions struct {
	Config      *config.Project
	AppName     string
	Backend     string
	FilePath    string
	SecretsPath string
	Force       bool
}

var secretsPushExample = templates.Examples(`
	# Push secrets:

    # This will push secrets for "squibby" app
    ize secrets push squibby
    
    # This will push secrets for "squibby" app from a "example-service.json" file to the AWS SSM storage with force option (values will be overwritten if exist)
	ize secrets push squibby --backend ssm --file example-service.json --force
`)

func NewSecretsPushFlags() *SecretsPushOptions {
	return &SecretsPushOptions{}
}

func NewCmdSecretsPush() *cobra.Command {
	o := NewSecretsPushFlags()

	cmd := &cobra.Command{
		Use:               "push <app>",
		Example:           secretsPushExample,
		Short:             "Push secrets to a key-value storage (like SSM)",
		Long:              "This command pushes secrets from a local file to a key-value storage (like SSM)",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			err := o.Complete(cmd)
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

func (o *SecretsPushOptions) Complete(cmd *cobra.Command) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg
	o.AppName = cmd.Flags().Args()[0]

	if o.FilePath == "" {
		o.FilePath = fmt.Sprintf("%s/%s/%s.json", o.Config.EnvDir, "secrets", o.AppName)
	}

	if o.SecretsPath == "" {
		o.SecretsPath = fmt.Sprintf("/%s/%s", o.Config.Env, o.AppName)
	}

	return nil
}

func (o *SecretsPushOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *SecretsPushOptions) Run() error {
	s, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Pushing secrets for %s...", o.AppName))
	if o.Backend == "ssm" {
		err := o.push(s)
		if err != nil {
			return fmt.Errorf("can't push secrets: %w", err)
		}
	} else {
		return fmt.Errorf("backend with type %s not found or not supported", o.Backend)
	}

	s.Success("Pushing secrets complete!")

	return nil
}

func (o *SecretsPushOptions) push(s *pterm.SpinnerPrinter) error {
	s.UpdateText("Reading secrets from file...")
	values, err := getKeyValuePairs(o.FilePath)
	if err != nil {
		return err
	}

	s.UpdateText(fmt.Sprintf("Pushing secrets to %s://%s...", o.Backend, o.SecretsPath))

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

	return nil
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
