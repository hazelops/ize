package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/template"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TfenvOptions struct {
	Config                   *config.Project
	TerraformStateBucketName string
}

var tfenvLongDesc = templates.LongDesc(`
	tfenv generates backend.tf and variable.tfvars files.
`)

var tfenvExample = templates.Examples(`
	# Generate files
	ize tfenv

	# Generate files via config file
	ize --config-file /path/to/config tfenv

	# Generate files via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize tfenv
`)

func NewTfenvFlags(project *config.Project) *TfenvOptions {
	return &TfenvOptions{
		Config: project,
	}
}

func NewCmdTfenv(project *config.Project) *cobra.Command {
	o := NewTfenvFlags(project)

	cmd := &cobra.Command{
		Use:     "tfenv",
		Short:   "Generate terraform files",
		Long:    tfenvLongDesc,
		Example: tfenvExample,
		Hidden:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			err := o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.TerraformStateBucketName, "terraform-state-bucket-name", "", "set terraform state bucket name (default <NAMESPACE>-tf-state)")

	return cmd
}

func (o *TfenvOptions) Run() error {
	return GenerateTerraformFiles("infra", o.TerraformStateBucketName, o.Config)

}

func GenerateTerraformFiles(name string, terraformStateBucketName string, project *config.Project) error {
	var tf *config.Terraform

	if b, ok := project.Terraform[name]; ok {
		tf = b
	} else {
		return fmt.Errorf("stack '%s' not found in config", name)
	}

	if len(terraformStateBucketName) != 0 {
		tf.StateBucketName = terraformStateBucketName
	}

	if len(tf.StateBucketName) == 0 {
		legacyBucketExists := checkTFStateBucket(project, fmt.Sprintf("%s-tf-state", project.Namespace))
		// If we found an existing bucket that conforms with the legacy format use it.
		if legacyBucketExists {
			tf.StateBucketName = fmt.Sprintf("%s-tf-state", project.Namespace)
		} else {
			resp, err := project.AWSClient.STSClient.GetCallerIdentity(
				&sts.GetCallerIdentityInput{},
			)
			if err != nil {
				return err
			}

			// If we haven't found an existing legacy format state bucket use a <NAMESPACE>-<AWS_ACCOUNT>-tf-state bucket as default (unless overridden with other parameters).
			tf.StateBucketName = fmt.Sprintf("%s-%s-tf-state", project.Namespace, *resp.Account)
		}
	}

	stateKey := fmt.Sprintf("%v/%v.tfstate", project.Env, name)
	if len(tf.StateName) != 0 {
		stateKey = fmt.Sprintf("%v/%v.tfstate", project.Env, tf.StateName)
	}

	if name == "infra" {
		if checkTFStateKey(project, tf.StateBucketName, filepath.Join(project.Env, "terraform.tfstate")) {
			stateKey = filepath.Join(project.Env, "terraform.tfstate")
			pterm.Warning.Printfln("%s/terraform.tfstate location is deprecated, please move to %s/infra.tfstate", project.Env, project.Env)
		} else {
			stateKey = filepath.Join(project.Env, "infra.tfstate")
		}
	}

	if len(tf.StateBucketRegion) == 0 {
		tf.StateBucketRegion = project.AwsRegion
	}

	backendOpts := template.BackendOpts{
		Env:                         project.Env,
		LocalstackEndpoint:          "",
		TerraformStateBucketName:    tf.StateBucketName,
		TerraformStateKey:           stateKey,
		TerraformStateRegion:        tf.StateBucketRegion,
		TerraformStateProfile:       project.AwsProfile,
		TerraformStateDynamodbTable: "tf-state-lock",
		TerraformAwsProviderVersion: "",
		Namespace:                   project.Namespace,
	}

	stackPath := filepath.Join(project.EnvDir, name)
	if name == "infra" {
		stackPath = project.EnvDir
	}

	if len(tf.TerraformConfigFile) == 0 {
		tf.TerraformConfigFile = "backend.tf"
	}

	logrus.Debugf("backend opts: %s", backendOpts)
	logrus.Debugf("state dir path: %s", stackPath)
	logrus.Debugf("config file name: %s", tf.TerraformConfigFile)

	err := template.GenerateBackendTf(
		backendOpts,
		filepath.Join(stackPath, tf.TerraformConfigFile),
	)
	if err != nil {
		pterm.Error.Printfln("Generate terraform file for \"%s\" not completed", name)
		return fmt.Errorf("can't generate backent.tf: %s", err)
	}

	home, _ := os.UserHomeDir()
	key, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", home))
	if err != nil {
		pterm.Error.Printfln("Generate terraform file for \"%s\" not completed", name)
		return fmt.Errorf("can't read public ssh key: %s", err)

	}

	// rootDomain := tf.RootDomainName

	varsOpts := template.VarsOpts{
		Env:            project.Env,
		AwsProfile:     project.AwsProfile,
		AwsRegion:      project.AwsRegion,
		EC2KeyPairName: fmt.Sprintf("%v-%v", project.Env, project.Namespace),
		RootDomainName: tf.RootDomainName,
		SSHPublicKey:   string(key)[:len(string(key))-1],
		Namespace:      project.Namespace,
	}

	if len(project.Ecs) != 0 {
		varsOpts.Tag = project.Tag
		varsOpts.DockerRegistry = project.DockerRegistry
	}

	logrus.Debugf("vars opts: %s", varsOpts)
	logrus.Debugf("state dir path: %s", stackPath)

	err = template.GenerateVarsTf(
		varsOpts,
		stackPath,
	)
	if err != nil {
		pterm.Error.Printfln("Generate terraform file for \"%s\" not completed", name)
		return fmt.Errorf("can't generate tfvars: %s", err)
	}

	pterm.Success.Printfln("Generate terraform file for \"%s\" completed", name)

	return nil
}

func checkTFStateBucket(project *config.Project, name string) bool {
	_, err := project.AWSClient.S3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return false
			default:
				return false
			}
		}
	}

	return true
}

func checkTFStateKey(project *config.Project, bucket, key string) bool {
	_, err := project.AWSClient.S3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return false
			default:
				return false
			}
		}
	}

	return true
}
