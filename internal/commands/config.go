package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type configCmd struct {
	*baseBuilderCmd

	vaultType string
	filePath  string
}

func (b *commandsBuilder) newConfigCmd() *configCmd {
	cc := &configCmd{}

	cmd := &cobra.Command{
		Use:              "config",
		Short:            "",
		Long:             "",
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

			if cc.vaultType == "ssm" {
				err = Set(cc.config.AwsRegion, cc.filePath, fmt.Sprintf("/%s/%s", cc.config.Env, cc.config.Namespace))
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

	setCmd.Flags().StringVar(&cc.vaultType, "type", "", "vault type")
	setCmd.Flags().StringVar(&cc.filePath, "file", "", "file with sercrets")

	cmd.AddCommand(setCmd)

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func Set(region string, file string, path string) error {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: region,
	})
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
