package secret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type SecretSetOptions struct {
	Config *config.Config
	Type   string
	Path   string
	Force  bool
}

func NewSecretSetFlags() *SecretSetOptions {
	return &SecretSetOptions{}
}

func NewCmdSecretSet() *cobra.Command {
	o := NewSecretSetFlags()

	cmd := &cobra.Command{
		Use:   "set",
		Short: "set secrets to storage",
		Long:  "This command set sercrets to storage.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete(cmd, args)
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

	cmd.Flags().StringVar(&o.Type, "type", "", "vault type")
	cmd.Flags().StringVar(&o.Path, "file", "", "file with sercrets")
	cmd.Flags().BoolVar(&o.Force, "force", false, "allow overwrite of parameters")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("file")

	return cmd
}

func (o *SecretSetOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	return nil
}

func (o *SecretSetOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *SecretSetOptions) Run() error {
	pterm.DefaultSection.Printfln("Starting config setting")

	if o.Type == "ssm" {

		err := set(o)
		if err != nil {
			pterm.DefaultSection.Println("Config setting not completed")
			return err
		}
	} else {
		pterm.DefaultSection.Println("Config setting not completed")
		return fmt.Errorf("vault with type %s not found or not supported", o.Type)
	}

	pterm.DefaultSection.Printfln("Config setting completed")

	return nil
}

func set(o *SecretSetOptions) error {
	basename := filepath.Base(o.Path)
	svc := strings.TrimSuffix(basename, filepath.Ext(basename))
	path := fmt.Sprintf("/%s/%s", o.Config.Env, svc)

	sess, err := utils.GetSession(
		&utils.SessionConfig{
			Region:  o.Config.AwsRegion,
			Profile: o.Config.AwsProfile,
		},
	)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Geting AWS session")

	values, err := getKeyValuePairs(o.Path)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Reading secrets from file")

	ssmSvc := ssm.New(sess)

	for key, value := range values {
		name := fmt.Sprintf("%s/%s", path, key)

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
					Value: &svc,
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

	pterm.Success.Printfln("Putting secrets in SSM")

	return err
}

func getKeyValuePairs(file string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(wd + "/" + file)
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
