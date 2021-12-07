package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type secretCmd struct {
	*baseBuilderCmd

	vaultType string
	filePath  string
	path      string
}

func (b *commandsBuilder) newSecretCmd() *secretCmd {
	cc := &secretCmd{}

	cmd := &cobra.Command{
		Use:              "secret",
		Short:            "manage secret",
		RunE:             nil,
		TraverseChildren: true,
	}

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Set secrets to storage",
		Long:  "This command set sercrets to storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			pterm.DefaultSection.Printfln("Starting config setting")

			basename := filepath.Base(cc.filePath)
			svc := strings.TrimSuffix(basename, filepath.Ext(basename))

			if cc.vaultType == "ssm" {

				err = Set(
					utils.SessionConfig{
						Region:  cc.config.AwsRegion,
						Profile: cc.config.AwsProfile,
					},
					cc.filePath,
					fmt.Sprintf("/%s/%s", cc.config.Env, svc),
					svc,
				)
				if err != nil {
					pterm.DefaultSection.Println("Config setting not completed")
					return err
				}
			} else {
				pterm.DefaultSection.Println("Config setting not completed")
				return fmt.Errorf("vault with type %s not found or not supported", cc.vaultType)
			}

			pterm.DefaultSection.Printfln("Config setting completed")

			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove secrets from storage",
		Long:  "This command remove sercrets from storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			pterm.DefaultSection.Printfln("Starting remove secrets")

			if cc.vaultType == "ssm" {
				err = Remove(
					utils.SessionConfig{
						Region:  cc.config.AwsRegion,
						Profile: cc.config.AwsProfile,
					},
					cc.path,
				)
				if err != nil {
					pterm.DefaultSection.Println("Remove secrets not completed")
					return err
				}
			} else {
				pterm.DefaultSection.Println("Remove secrets not completed")
				return fmt.Errorf("vault with type %s not found or not supported", cc.vaultType)
			}

			pterm.DefaultSection.Printfln("Remove secrets completed")

			return nil
		},
	}

	removeCmd.Flags().StringVar(&cc.vaultType, "type", "", "vault type")
	removeCmd.Flags().StringVar(&cc.path, "path", "", "path to secrets")
	removeCmd.MarkFlagRequired("type")
	removeCmd.MarkFlagRequired("path")

	setCmd.Flags().StringVar(&cc.vaultType, "type", "", "vault type")
	setCmd.Flags().StringVar(&cc.filePath, "file", "", "file with sercrets")
	setCmd.MarkFlagRequired("type")
	setCmd.MarkFlagRequired("file")

	cmd.AddCommand(setCmd, removeCmd)

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func Remove(sessCfg utils.SessionConfig, path string) error {
	if path == "" {
		pterm.Info.Printfln("Path were not set")
		return nil
	}

	sess, err := utils.GetSession(&sessCfg)
	if err != nil {
		return err
	}
	pterm.Success.Printfln("Geting AWS session")

	ssmSvc := ssm.New(sess)

	out, err := ssmSvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path: &path,
	})
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Getting secrets from path")

	if len(out.Parameters) == 0 {
		pterm.Info.Printfln("No values ​​found along the path")
		pterm.Success.Printfln("Deleting secrets from path")
		return nil
	}

	var names []*string

	for _, p := range out.Parameters {
		names = append(names, p.Name)
	}

	_, err = ssmSvc.DeleteParameters(&ssm.DeleteParametersInput{
		Names: names,
	})

	if err != nil {
		return err
	}

	pterm.Success.Printfln("Deleting secrets from path")

	return nil
}

func Set(sessCfg utils.SessionConfig, file string, path string, svc string) error {
	sess, err := utils.GetSession(&sessCfg)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Geting AWS session")

	values, err := getKeyValuePairs(file)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Reading secrets from file")

	ssmSvc := ssm.New(sess)

	for key, value := range values {
		name := fmt.Sprintf("%s/%s", path, key)

		_, err := ssmSvc.PutParameter(&ssm.PutParameterInput{
			Name:  &name,
			Value: aws.String(value),
			Type:  aws.String(ssm.ParameterTypeSecureString),
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

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ParameterAlreadyExists":
				return fmt.Errorf("secret already exists, you can use --force to overwrite it")
			default:
				return err
			}
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
