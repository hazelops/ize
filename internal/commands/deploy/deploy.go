package deploy

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/hazelops/ize/internal/services"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeployOptions struct {
	Config           *config.Config
	ServiceName      string
	Tag              string
	SkipBuildAndPush bool
	Services         Services
	Infra            Infra
	App              services.App
	AutoApprove      bool
}

type Services map[string]*services.App

type Infra struct {
	Version string `mapstructure:"terraform_version"`
	Region  string `mapstructure:"aws_region"`
	Profile string `mapstructure:"aws_profile"`
}

var deployLongDesc = templates.LongDesc(`
	Deploy infraftructure or sevice.
           Service name must be specified for a service deploy. 
	The infrastructure for the service must be prepared in advance.
`)

var deployExample = templates.Examples(`
	# Deploy all (config file required)
	ize deploy

	# Deploy service (config file required)
	ize deploy <service name>

	# Deploy service via config file
	ize --config-file (or -c) /path/to/config deploy <service name>

	# Deploy service via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize deploy <service name>
`)

func NewDeployFlags() *DeployOptions {
	return &DeployOptions{}
}

func NewCmdDeploy(ui terminal.UI) *cobra.Command {
	o := NewDeployFlags()

	cmd := &cobra.Command{
		Use:     "deploy [flags] <service name>",
		Example: deployExample,
		Short:   "manage deployments",
		Long:    deployLongDesc,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if len(args) == 0 && !o.AutoApprove {
				pterm.Warning.Println("please set flag --auto-approve")
				return nil
			}

			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run(ui)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&o.AutoApprove, "auto-approve", false, "approve deploy all")

	cmd.AddCommand(
		NewCmdDeployInfra(ui),
	)

	return cmd
}

func (o *DeployOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		o.Config, err = config.InitializeConfig(config.WithConfigFile())
		viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("can`t complete options: %w", err)
		}

		viper.UnmarshalKey("service", &o.Services)
		viper.UnmarshalKey("infra.terraform", &o.Infra)

		for _, v := range o.Services {
			fmt.Println(*v)
		}

		if len(o.Infra.Profile) == 0 {
			o.Infra.Profile = o.Config.AwsProfile
		}

		if len(o.Infra.Region) == 0 {
			o.Infra.Region = o.Config.AwsRegion
		}

		if len(o.Infra.Version) == 0 {
			o.Infra.Version = viper.GetString("terraform_version")
		}
	} else {
		o.Config, err = config.InitializeConfig(config.WithDocker(), config.WithConfigFile())
		viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("can`t complete options: %w", err)
		}
		o.ServiceName = cmd.Flags().Args()[0]
		viper.UnmarshalKey(fmt.Sprintf("service.%s", o.ServiceName), &o.App)
	}

	o.Tag = viper.GetString("tag")

	return nil
}

func (o *DeployOptions) Validate() error {
	if o.ServiceName == "" {
		err := validateAll(o)
		if err != nil {
			return err
		}
	} else {
		err := validate(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *DeployOptions) Run(ui terminal.UI) error {
	if o.ServiceName == "" {
		err := deployAll(ui, o)
		if err != nil {
			return err
		}
	} else {
		err := deployService(ui, o)
		if err != nil {
			return err
		}
	}

	return nil
}

func validate(o *DeployOptions) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified\n")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified\n")
	}

	if len(o.ServiceName) == 0 {
		return fmt.Errorf("can't validate options: service name be specified\n")
	}

	return nil
}

func validateAll(o *DeployOptions) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified\n")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified\n")
	}

	for sname, svc := range o.Services {
		if len(svc.Type) == 0 {
			return fmt.Errorf("can't validate options: type for service %s must be specified\n", sname)
		}
	}

	return nil
}

func deployAll(ui terminal.UI, o *DeployOptions) error {
	ui.Output("Running deploy infra...", terminal.WithHeaderStyle())

	logrus.Infof("infra: %s", o.Infra)

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", o.Config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", o.Infra.Profile),
		fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
		fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	//terraform init run options
	opts := terraform.Options{
		ContainerName:    "terraform",
		Cmd:              []string{"init", "-input=true"},
		Env:              env,
		TerraformVersion: o.Infra.Version,
	}

	ui.Output("Execution terraform init...", terminal.WithHeaderStyle())

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	ui.Output("Execution terraform plan...", terminal.WithHeaderStyle())

	//terraform plan run options
	opts.Cmd = []string{"plan"}

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return err
	}

	//terraform apply run options
	opts.Cmd = []string{"apply", "-auto-approve"}

	ui.Output("Execution terraform apply...", terminal.WithHeaderStyle())

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return err
	}

	//terraform output run options
	outputPath := fmt.Sprintf("%s/.terraform/output.json", viper.Get("ENV_DIR"))

	opts.Cmd = []string{"output", "-json"}
	opts.OutputPath = outputPath

	ui.Output("Execution terraform output...", terminal.WithHeaderStyle())

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("/%s/terraform-output", o.Config.Env)

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

	ssm.New(o.Config.Session).PutParameter(&ssm.PutParameterInput{
		Name:      &name,
		Value:     aws.String(string(sDec)),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(true),
		Tier:      aws.String("Intelligent-Tiering"),
		DataType:  aws.String("text"),
	})

	logrus.Debug(o.Services)

	ui.Output("Deploying apps...", terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	err = InDependencyOrder(aws.BackgroundContext(), &o.Services, func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Infra.Profile

		o.Services[name].Name = name
		err := o.Services[name].Deploy(sg, ui)
		if err != nil {
			return fmt.Errorf("can't deploy all: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	s := sg.Add("Deploy all completed!")
	s.Done()

	return nil
}

func deployService(ui terminal.UI, o *DeployOptions) error {
	ui.Output("Deploying %s app...", o.ServiceName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	o.App.Name = o.ServiceName
	err := o.App.Deploy(sg, ui)
	if err != nil {
		return fmt.Errorf("can't deploy: %w", err)
	}

	ui.Output("deploy service %s completed\n", o.ServiceName, terminal.WithSuccessStyle())

	return nil
}
