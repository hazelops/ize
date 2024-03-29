package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

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
	Explain     bool
}

var explainSecretsPushTmpl = `
SERVICE_SECRETS_FILE={{.EnvDir}}/secrets/{{svc}}.json
SERVICE_SECRETS=$(cat $SERVICE_SECRETS_FILE | jq -e -r '. | keys[]')
for item in $(echo $SERVICE_SECRETS); do 
    aws --profile={{.AwsProfile}} ssm put-parameter --name="/{{.Env}}/{{svc}}/${item}" --value="$(cat $SERVICE_SECRETS_FILE | jq -r .$item )" --type SecureString --overwrite && \
    aws --profile={{.AwsProfile}} ssm add-tags-to-resource --resource-type "Parameter" --resource-id "/{{.Env}}/{{svc}}/${item}" \
    --tags "Key=Application,Value={{svc}}" "Key=EnvVarName,Value=${item}"
done
`

var secretsPushExample = templates.Examples(`
	# Push secrets:

    # This will push secrets for "squibby" app
    ize secrets push squibby
    
    # This will push secrets for "squibby" app from a "example-service.json" file to the AWS SSM storage with force option (values will be overwritten if exist)
	ize secrets push squibby --backend ssm --file example-service.json --force
`)

func NewSecretsPushFlags(project *config.Project) *SecretsPushOptions {
	return &SecretsPushOptions{
		Config: project,
	}
}

func NewCmdSecretsPush(project *config.Project) *cobra.Command {
	o := NewSecretsPushFlags(project)

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
	cmd.Flags().BoolVar(&o.Explain, "explain", false, "bash alternative shown")
	cmd.Flags().BoolVar(&o.Force, "force", false, "allow values overwrite")

	return cmd
}

func (o *SecretsPushOptions) Complete(cmd *cobra.Command) error {
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
	if o.Explain {
		err := o.Config.Generate(explainSecretsPushTmpl, template.FuncMap{
			"svc": func() string {
				return o.AppName
			},
		})
		if err != nil {
			return err
		}

		return nil
	}

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

	for key, value := range values {
		name := fmt.Sprintf("%s/%s", o.SecretsPath, key)

		_, err := o.Config.AWSClient.SSMClient.PutParameter(&ssm.PutParameterInput{
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

		_, err = o.Config.AWSClient.SSMClient.AddTagsToResource(&ssm.AddTagsToResourceInput{
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

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

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
