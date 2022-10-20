package template

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const (
	vars = "terraform.tfvars"
	ize  = "ize.hcl"
)

func GenerateVarsTf(opts VarsOpts, path string) error {
	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()

	rootBody.SetAttributeValue("env", cty.StringVal(opts.Env))
	rootBody.SetAttributeValue("aws_profile", cty.StringVal(opts.AwsProfile))
	rootBody.SetAttributeValue("aws_region", cty.StringVal(opts.AwsRegion))
	rootBody.SetAttributeValue("ec2_key_pair_name", cty.StringVal(opts.EC2KeyPairName))
	if len(opts.Tag) != 0 {
		rootBody.SetAttributeValue("docker_image_tag", cty.StringVal(opts.Tag))
	}
	rootBody.SetAttributeValue("ssh_public_key", cty.StringVal(opts.SSHPublicKey))
	if len(opts.DockerRegistry) != 0 {
		rootBody.SetAttributeValue("docker_registry", cty.StringVal(opts.DockerRegistry))
	}
	rootBody.SetAttributeValue("namespace", cty.StringVal(opts.Namespace))
	if len(opts.RootDomainName) > 0 {
		rootBody.SetAttributeValue("root_domain_name", cty.StringVal(opts.RootDomainName))
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

func GenerateBackendTf(opts BackendOpts, path string) error {
	f := hclwrite.NewEmptyFile()

	if strings.Contains(opts.Env, "localstack") || strings.Contains(opts.Env, "local") {
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
		endpointBlock.Body().SetAttributeValue("apigateway", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("acm", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("cloudformation", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("cloudwatch", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("ec2", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("dynamodb", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("es", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("firehose", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("iam", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("kinesis", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("lambda", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("route53", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("redshift", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("s3", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("secretsmanager", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("ses", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("sns", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("sqs", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("ssm", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("stepfunctions", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("sts", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("ecs", cty.StringVal(opts.LocalstackEndpoint))
		endpointBlock.Body().SetAttributeValue("ecr", cty.StringVal(opts.LocalstackEndpoint))
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
			"env":       cty.StringVal(opts.Env),
			"namespace": cty.StringVal(opts.Namespace),
		}))

		rootBody.AppendNewline()

		// Terraform block
		terraformBlock := f.Body().AppendNewBlock("terraform", []string{})
		// backend s3 block
		backendBlock := terraformBlock.Body().AppendNewBlock("backend", []string{"s3"})
		backendBlock.Body().SetAttributeValue("bucket", cty.StringVal(opts.TerraformStateBucketName))
		backendBlock.Body().SetAttributeValue("key", cty.StringVal(opts.TerraformStateKey))
		backendBlock.Body().SetAttributeValue("region", cty.StringVal(opts.TerraformStateRegion))
		backendBlock.Body().SetAttributeValue("profile", cty.StringVal(opts.TerraformStateProfile))
		backendBlock.Body().SetAttributeValue("dynamodb_table", cty.StringVal(opts.TerraformStateDynamodbTable))
	}

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
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
	oldFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	oldHash := md5.Sum(oldFile)

	if !reflect.DeepEqual(newHash, oldHash) {
		file, err := os.Create(path)
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
	Env            string
	AwsProfile     string
	AwsRegion      string
	EC2KeyPairName string
	RootDomainName string
	Tag            string
	SSHPublicKey   string
	DockerRegistry string
	Namespace      string
}

type BackendOpts struct {
	Namespace                   string
	Env                         string
	LocalstackEndpoint          string
	TerraformStateBucketName    string
	TerraformStateKey           string
	TerraformStateRegion        string
	TerraformStateProfile       string
	TerraformStateDynamodbTable string
	TerraformAwsProviderVersion string
}
