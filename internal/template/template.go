package template

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/hairyhenderson/gomplate/conv"
)

const (
	backend = "backend.tf"
	vars    = "terraform.tfvars"
)

func GenereateBackendTf(opts BackendOpts, path string) error {
	tmpl, err := template.New(backend).Funcs(CreateStringFuncs(context.Background())).Parse(backendTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", path, backend))
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, opts)
	if err != nil {
		return err
	}

	file.Close()

	return nil
}

func GenerateVarsTf(opts VarsOpts, path string) error {
	tmpl, err := template.New(vars).Funcs(CreateStringFuncs(context.Background())).Parse(varsTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", path, vars))
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, opts)
	if err != nil {
		return err
	}

	file.Close()

	return nil
}

type VarsOpts struct {
	ENV               string
	AWS_PROFILE       string
	AWS_REGION        string
	EC2_KEY_PAIR_NAME string
	TAG               string
	SSH_PUBLIC_KEY    string
	DOCKER_REGISTRY   string
	NAMESPACE         string
}

type BackendOpts struct {
	ENV                            string
	LOCALSTACK_ENDPOINT            string
	TERRAFORM_STATE_BUCKET_NAME    string
	TERRAFORM_STATE_KEY            string
	TERRAFORM_STATE_REGION         string
	TERRAFORM_STATE_PROFILE        string
	TERRAFORM_STATE_DYNAMODB_TABLE string
	TERRAFORM_AWS_PROVIDER_VERSION string
}

func CreateStringFuncs(ctx context.Context) map[string]interface{} {
	f := map[string]interface{}{}

	ns := &StringFuncs{ctx}
	f["strings"] = func() interface{} { return ns }

	f["contains"] = strings.Contains

	return f
}

// Contains -
func (StringFuncs) Contains(substr string, s interface{}) bool {
	return strings.Contains(conv.ToString(s), substr)
}

type StringFuncs struct {
	ctx context.Context
}
