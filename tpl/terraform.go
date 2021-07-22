package tpl

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

func GenerateBackendTf()  {
	type TfVars struct {
		ENV string
		AWS_PROFILE string
		AWS_REGION string
		EC2_KEY_PAIR_NAME string
		TAG string
		SSH_PUBLIC_KEY string
		DOCKER_REGISTRY string
		NAMESPACE string
	}
	type Backend struct {
		LOCALSTACK_ENDPOINT string
		TERRAFORM_STATE_BUCKET_NAME string
		TERRAFORM_STATE_KEY string
		TERRAFORM_STATE_REGION string
		TERRAFORM_STATE_PROFILE string
		TERRAFORM_STATE_DYNAMODB_TABLE string
		TERRAFORM_AWS_PROVIDER_VERSION string
	}

	tfvars := TfVars{
		fmt.Sprintf("%v",viper.Get("ENV")),
		fmt.Sprintf("%v",viper.Get("AWS_PROFILE")),
		fmt.Sprintf("%v",viper.Get("AWS_REGION")),
		fmt.Sprintf("%v",viper.Get("EC2_KEY_PAIR_NAME")),
		fmt.Sprintf("%v",viper.Get("TAG")),
		fmt.Sprintf("%v",viper.Get("SSH_PUBLIC_KEY")),
		fmt.Sprintf("%v",viper.Get("DOCKER_REGISTRY")),
		fmt.Sprintf("%v",viper.Get("NAMESPACE")),
	}

	tmpl, err := template.New("backend.tf").Parse(BackendTfTemplate)
	cobra.CheckErr(err)

	err = tmpl.Execute(os.Stdout, tfvars)
	cobra.CheckErr(err)

}