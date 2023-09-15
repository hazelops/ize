package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var explainSecretsPullTmpl = `
aws ssm get-parameters-by-path \
	--path "/{{.Env}}/{{svc}}" \
	--with-decryption \
	--recursive \
	--parameter-filters "Key=Type,Values=SecureString" \
	--output json | jq '.Parameters | [.[] | {(.Name|capture(".*/(?<a>.*)").a): .Value}]|reduce .[] as $item ({}; . + $item)' > {{.EnvDir}}/secrets/{{svc}}.json
`

type SecretsPullOptions struct {
	Config      *config.Project
	AppName     string
	Backend     string
	FilePath    string
	SecretsPath string
	Force       bool
	Explain     bool
}

func NewSecretsPullFlags(project *config.Project) *SecretsPullOptions {
	return &SecretsPullOptions{
		Config: project,
	}
}

func NewCmdSecretsPull(project *config.Project) *cobra.Command {
	o := NewSecretsPullFlags(project)

	cmd := &cobra.Command{
		Use:               "pull",
		Short:             "Pull secrets to a a local file (like SSM)",
		Long:              "This command pulls secrets from a key-value storage to a local file (like SSM)",
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

func (o *SecretsPullOptions) Complete(cmd *cobra.Command) error {
	o.AppName = cmd.Flags().Args()[0]

	if o.FilePath == "" {
		o.FilePath = fmt.Sprintf("%s/%s/%s.json", o.Config.EnvDir, "secrets", o.AppName)
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
	if o.Explain {
		err := o.Config.Generate(explainSecretsPullTmpl, template.FuncMap{
			"svc": func() string {
				return o.AppName
			},
		})
		if err != nil {
			return err
		}

		return nil
	}

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

	values := make(map[string]interface{})

	params, err := o.Config.AWSClient.SSMClient.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(o.SecretsPath),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
		ParameterFilters: []*ssm.ParameterStringFilter{
			{
				Key:    aws.String("Type"),
				Values: aws.StringSlice([]string{"SecureString"}),
			},
		},
	})
	if err != nil {
		return err
	}

	for _, param := range params.Parameters {
		p := strings.Split(*param.Name, "/")
		values[p[len(p)-1]] = *param.Value
	}

	for {
		if params.NextToken == nil {
			break
		}

		params, err = o.Config.AWSClient.SSMClient.GetParametersByPath(&ssm.GetParametersByPathInput{
			Path:           aws.String(o.SecretsPath),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
			NextToken:      params.NextToken,
			ParameterFilters: []*ssm.ParameterStringFilter{
				{
					Key:    aws.String("Type"),
					Values: aws.StringSlice([]string{"SecureString"}),
				},
			},
		})
		if err != nil {
			return err
		}

		for _, param := range params.Parameters {
			p := strings.Split(*param.Name, "/")
			values[p[len(p)-1]] = *param.Value
		}
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
		return fmt.Errorf("file %s already exists. Please use --force to overwrite", o.FilePath)
	}

	return nil
}
