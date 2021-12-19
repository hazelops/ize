package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/mitchellh/mapstructure"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type deployCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newDeployCmd() *deployCmd {
	cc := &deployCmd{}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Manage deployments.",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "infra",
		Short: "Deploy infrastructures.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			cc.log.Infof("infra: %s", cc.config.Infra)

			for pname, provider := range cc.config.Infra {
				switch pname {
				case "terraform":
					var tic terraformInfraConfig
					mapstructure.Decode(provider, &tic)

					//terraform init
					opts := terraform.Options{
						ContainerName: "terraform",
						Cmd:           []string{"init", "-input=true"},
						Env: []string{
							fmt.Sprintf("ENV=%v", cc.config.Env),
							fmt.Sprintf("AWS_PROFILE=%v", tic.Profile),
							fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
							fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
						},
						TerraformVersion: tic.Version,
					}

					spinner := &pterm.SpinnerPrinter{}

					if cc.log.Level < 4 {
						spinner, _ = pterm.DefaultSpinner.Start("execution terraform init")
					}

					err = terraform.Run(&cc.log, opts)
					if err != nil {
						cc.log.Errorf("terraform %s not completed", "init")
						return err
					}

					if cc.log.Level < 4 {
						spinner.Success("terrafrom init completed")
					} else {
						pterm.Success.Println("terrafrom init completed")
					}

					//terraform plan
					opts = terraform.Options{
						ContainerName: "terraform",
						Cmd:           []string{"plan"},
						Env: []string{
							fmt.Sprintf("ENV=%v", cc.config.Env),
							fmt.Sprintf("AWS_PROFILE=%v", tic.Profile),
							fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
							fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
						},
						TerraformVersion: tic.Version,
					}

					if cc.log.Level < 4 {
						spinner, _ = pterm.DefaultSpinner.Start("execution terraform plan")
					}

					err = terraform.Run(&cc.log, opts)
					if err != nil {
						cc.log.Errorf("terraform %s not completed", "plan")
						return err
					}

					if cc.log.Level < 4 {
						spinner.Success("terrafrom plan completed")
					} else {
						pterm.Success.Println("terrafrom plan completed")
					}

					//terraform apply
					opts = terraform.Options{
						ContainerName: "terraform",
						Cmd:           []string{"apply", "-auto-approve"},
						Env: []string{
							fmt.Sprintf("ENV=%v", cc.config.Env),
							fmt.Sprintf("AWS_PROFILE=%v", tic.Profile),
							fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
							fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
						},
						TerraformVersion: tic.Version,
					}

					if cc.log.Level < 4 {
						spinner, _ = pterm.DefaultSpinner.Start("execution terraform apply")
					}

					err = terraform.Run(&cc.log, opts)
					if err != nil {
						cc.log.Errorf("terraform %s not completed", "apply")
						return err
					}

					if cc.log.Level < 4 {
						spinner.Success("terrafrom apply completed")
					} else {
						pterm.Success.Println("terrafrom apply completed")
					}

					// terraform output
					outputPath := fmt.Sprintf("%s/.terraform/output.json", viper.Get("ENV_DIR"))

					fmt.Println(outputPath)

					opts = terraform.Options{
						ContainerName: "terraform",
						Cmd:           []string{"output", "-json"},
						Env: []string{
							fmt.Sprintf("ENV=%v", cc.config.Env),
							fmt.Sprintf("AWS_PROFILE=%v", tic.Profile),
							fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
							fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
						},
						TerraformVersion: tic.Version,
						OutputPath:       outputPath,
					}

					if cc.log.Level < 4 {
						spinner, _ = pterm.DefaultSpinner.Start("execution terraform output")
					}

					err = terraform.Run(&cc.log, opts)
					if err != nil {
						cc.log.Errorf("terraform %s not completed", "output")
						return err
					}

					if cc.log.Level < 4 {
						spinner.Success("terrafrom output completed")
					} else {
						pterm.Success.Println("terrafrom output completed")
					}

					sess, err := utils.GetSession(&utils.SessionConfig{
						Region:  cc.config.AwsRegion,
						Profile: cc.config.AwsProfile,
					})
					if err != nil {
						return err
					}

					name := fmt.Sprintf("/%s/terraform-output", cc.config.Env)

					outputFile, err := os.Open(outputPath)
					if err != nil {
						return err
					}

					defer outputFile.Close()

					byteValue, _ := ioutil.ReadAll(outputFile)
					sDec := base64.StdEncoding.EncodeToString(byteValue)
					if err != nil {
						return err
					}

					ssm.New(sess).PutParameter(&ssm.PutParameterInput{
						Name:      &name,
						Value:     aws.String(string(sDec)),
						Type:      aws.String(ssm.ParameterTypeSecureString),
						Overwrite: aws.Bool(true),
						Tier:      aws.String("Intelligent-Tiering"),
						DataType:  aws.String("text"),
					})

				default:
					return fmt.Errorf("provider %s is not supported", pname)
				}
			}

			return nil
		},
	})

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

type terraformInfraConfig struct {
	RootDir string `mapstructure:"root_dir,optional"`
	Version string `mapstructure:"terraform_version,optional"`
	Region  string `mapstructure:"aws_region,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
