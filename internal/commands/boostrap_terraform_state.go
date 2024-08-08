package commands

import (
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/spf13/cobra"
	"text/template"
)

type BootstrapTerraformStateOptions struct {
	Config  *config.Project
	AppName string
	Explain bool
}

var boostrapTerraformStateExplainTmpl = `
SERVICE_SECRETS_FILE={{.EnvDir}}/secrets/{{svc}}.json
SERVICE_SECRETS=$(cat $SERVICE_SECRETS_FILE | jq -e -r '. | keys[]')
for item in $(echo $SERVICE_SECRETS); do 
    aws --profile={{.AwsProfile}} ssm put-parameter --name="/{{.Env}}/{{svc}}/${item}" --value="$(cat $SERVICE_SECRETS_FILE | jq -r .$item )" --type SecureString --overwrite && \
    aws --profile={{.AwsProfile}} ssm add-tags-to-resource --resource-type "Parameter" --resource-id "/{{.Env}}/{{svc}}/${item}" \
    --tags "Key=Application,Value={{svc}}" "Key=EnvVarName,Value=${item}"
done
`

var boostrapTerraformStateExample = templates.Examples(`
	# Boostrap Terraform State:

    TBD
`)

func NewBoostrapTerraformStateFlags(project *config.Project) *BootstrapTerraformStateOptions {
	return &BootstrapTerraformStateOptions{
		Config: project,
	}
}

func NewBoostrapTerraformState(project *config.Project) *cobra.Command {
	o := NewBoostrapTerraformStateFlags(project)

	cmd := &cobra.Command{
		Use:     "terraform-state",
		Example: boostrapTerraformStateExample,
		Short:   "Boostrap Terraform State",
		Long:    "This command creates Terraform State bucket and DynamoDB table based on ize name convention",
		//Args:              cobra.MinimumNArgs(1),
		//ValidArgsFunction: config.GetApps,
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

	return cmd
}

func (o *BootstrapTerraformStateOptions) Complete(cmd *cobra.Command) error {
	//o.AppName = cmd.Flags().Args()[0]

	println("Done")
	return nil
}

func (o *BootstrapTerraformStateOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}
	println("Valid")

	return nil
}

func (o *BootstrapTerraformStateOptions) Run() error {
	if o.Explain {
		err := o.Config.Generate(boostrapTerraformStateExplainTmpl, template.FuncMap{
			"svc": func() string {
				return o.AppName
			},
		})
		if err != nil {
			return err
		}

		return nil
	}
	println("Running")
	// TODO: Create bucket via S3 Client with forcePAth ON
	//name := "test"
	//
	//_, err := o.Config.AWSClient.S3Client().CreateBucket(&s3.CreateBucketInput{
	//	Bucket: &name,
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}
