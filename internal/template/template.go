package template

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pterm/pterm"
	"github.com/zclconf/go-cty/cty"
)

const (
	backend = "backend.tf"
	vars    = "terraform.tfvars"
	ize     = "ize.hcl"
)

func GenerateVarsTf(opts VarsOpts, path string) error {
	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()

	rootBody.SetAttributeValue("env", cty.StringVal(opts.ENV))
	rootBody.SetAttributeValue("aws_profile", cty.StringVal(opts.AWS_PROFILE))
	rootBody.SetAttributeValue("aws_region", cty.StringVal(opts.AWS_REGION))
	rootBody.SetAttributeValue("ec2_key_pair_name", cty.StringVal(opts.EC2_KEY_PAIR_NAME))
	rootBody.SetAttributeValue("docker_image_tag", cty.StringVal(opts.TAG))
	rootBody.SetAttributeValue("ssh_public_key", cty.StringVal(opts.SSH_PUBLIC_KEY))
	rootBody.SetAttributeValue("docker_registry", cty.StringVal(opts.DOCKER_REGISTRY))
	rootBody.SetAttributeValue("namespace", cty.StringVal(opts.NAMESPACE))
	if len(opts.ROOT_DOMAIN_NAME) > 0 {
		rootBody.SetAttributeValue("root_domain_name", cty.StringVal(opts.ROOT_DOMAIN_NAME))
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", path, vars))
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = f.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

func GenerateConfigFile(opts ConfigOpts, path string) error {
	if !filepath.IsAbs(path) {
		if path == "" {
			path += ize
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		path = filepath.Join(wd, path)
	}

	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()

	rootBody.SetAttributeValue("env", cty.StringVal(opts.ENV))
	rootBody.SetAttributeValue("aws_profile", cty.StringVal(opts.AWS_PROFILE))
	rootBody.SetAttributeValue("aws_region", cty.StringVal(opts.AWS_REGION))
	rootBody.SetAttributeValue("terraform_version", cty.StringVal(opts.TERRAFORM_VERSION))
	rootBody.SetAttributeValue("namespace", cty.StringVal(opts.NAMESPACE))

	var owr bool = false

	_, err := os.Stat(path)
	if err == nil {
		var qs = []*survey.Question{
			{
				Prompt: &survey.Confirm{
					Message: " The file already exists. Overwrite?",
				},
				Validate: survey.Required,
				Name:     "owr",
			},
		}

		err = survey.Ask(qs, &owr, survey.WithIcons(func(is *survey.IconSet) {
			is.Question.Text = " ??"
			is.Question.Format = "black:green"
			is.Error.Format = "black:red"
		}))
		if err != nil {
			return err
		}

		if !owr {
			return nil
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = f.WriteTo(file)
	if err != nil {
		return err
	}

	if owr {
		pterm.Success.Println("Config file overwritten successfully")
	}

	pterm.Success.Println("Config file created successfully")

	return nil
}

func GenerateBackendTf(opts BackendOpts, path string) error {
	f := hclwrite.NewEmptyFile()

	if strings.Contains(opts.ENV, "localstack") || strings.Contains(opts.ENV, "local") {
		rootBody := f.Body()
		// AWS Provider block
		providerBlock := rootBody.AppendNewBlock("provider", []string{"aws"})
		providerBlock.Body().SetAttributeTraversal("profile", hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "aws_profile"},
		})
		providerBlock.Body().SetAttributeTraversal("region", hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "aws_region"},
		})
		providerBlock.Body().SetAttributeValue("s3_force_path_style", cty.True)
		providerBlock.Body().SetAttributeValue("secret_key", cty.StringVal("mock_secret_key"))
		providerBlock.Body().SetAttributeValue("skip_credentials_validation", cty.True)
		providerBlock.Body().SetAttributeValue("skip_metadata_api_check", cty.True)
		providerBlock.Body().SetAttributeValue("skip_requesting_account_id", cty.True)
		rootBody.AppendNewline()

		// Endpoints
		endpointBlock := rootBody.AppendNewBlock("endpoints", []string{})
		endpointBlock.Body().SetAttributeValue("apigateway", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("acm", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("cloudformation", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("cloudwatch", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("ec2", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("dynamodb", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("es", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("firehose", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("iam", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("kinesis", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("lambda", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("route53", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("redshift", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("s3", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("secretsmanager", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("ses", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("sns", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("sqs", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("ssm", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("stepfunctions", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("sts", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("ecs", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
		endpointBlock.Body().SetAttributeValue("ecr", cty.StringVal(opts.LOCALSTACK_ENDPOINT))
	} else {
		rootBody := f.Body()
		// AWS Provider block
		providerBlock := rootBody.AppendNewBlock("provider", []string{"aws"})
		providerBlock.Body().SetAttributeTraversal("profile", hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "aws_profile"},
		})
		providerBlock.Body().SetAttributeTraversal("region", hcl.Traversal{
			hcl.TraverseRoot{Name: "var"},
			hcl.TraverseAttr{Name: "aws_region"},
		})
		defaultTagsBlock := providerBlock.Body().AppendNewBlock("default_tags", nil)
		defaultTagsBlock.Body().SetAttributeValue("tags", cty.ObjectVal(map[string]cty.Value{
			"terraform": cty.StringVal("true"),
			"env":       cty.StringVal(opts.ENV),
			"namespace": cty.StringVal(opts.NAMESPACE),
		}))

		rootBody.AppendNewline()

		// Terraform block
		terraformBlock := f.Body().AppendNewBlock("terraform", []string{})
		// backend s3 block
		backendBlock := terraformBlock.Body().AppendNewBlock("backend", []string{"s3"})
		backendBlock.Body().SetAttributeValue("bucket", cty.StringVal(opts.TERRAFORM_STATE_BUCKET_NAME))
		backendBlock.Body().SetAttributeValue("key", cty.StringVal(opts.TERRAFORM_STATE_KEY))
		backendBlock.Body().SetAttributeValue("region", cty.StringVal(opts.TERRAFORM_STATE_REGION))
		backendBlock.Body().SetAttributeValue("profile", cty.StringVal(opts.TERRAFORM_STATE_PROFILE))
		backendBlock.Body().SetAttributeValue("dynamodb_table", cty.StringVal(opts.TERRAFORM_STATE_DYNAMODB_TABLE))
	}

	backendPath := filepath.Join(path, backend)

	_, err := os.Stat(backendPath)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(backendPath)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = f.WriteTo(file)
		if err != nil {
			return err
		}

		return nil
	}

	newHash := md5.Sum(f.Bytes())
	oldFile, err := os.ReadFile(backendPath)
	if err != nil {
		return err
	}

	oldHash := md5.Sum(oldFile)

	if !reflect.DeepEqual(newHash, oldHash) {
		file, err := os.Create(backendPath)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = f.WriteTo(file)
		if err != nil {
			return err
		}
	}

	return nil
}

type VarsOpts struct {
	ENV               string
	AWS_PROFILE       string
	AWS_REGION        string
	EC2_KEY_PAIR_NAME string
	ROOT_DOMAIN_NAME  string
	TAG               string
	SSH_PUBLIC_KEY    string
	DOCKER_REGISTRY   string
	NAMESPACE         string
}

type BackendOpts struct {
	NAMESPACE                      string
	ENV                            string
	LOCALSTACK_ENDPOINT            string
	TERRAFORM_STATE_BUCKET_NAME    string
	TERRAFORM_STATE_KEY            string
	TERRAFORM_STATE_REGION         string
	TERRAFORM_STATE_PROFILE        string
	TERRAFORM_STATE_DYNAMODB_TABLE string
	TERRAFORM_AWS_PROVIDER_VERSION string
}

type ConfigOpts struct {
	ENV               string
	AWS_PROFILE       string
	AWS_REGION        string
	TERRAFORM_VERSION string
	NAMESPACE         string
}
