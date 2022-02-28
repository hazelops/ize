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
	"github.com/hazelops/ize/internal/docker/ecsdeploy"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/hazelops/ize/pkg/templates"
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
	Service          ecsdeploy.Service
	AutoApprove      bool
}

type Services map[string]*ecsdeploy.Service

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

func NewCmdDeploy() *cobra.Command {
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

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.Service.Image, "image", "", "set image name")
	cmd.Flags().StringVar(&o.Service.EcsCluster, "ecs-cluster", "", "set ECS cluster name")
	cmd.Flags().StringVar(&o.Service.Path, "path", "", "specify the path to the service")
	cmd.Flags().StringVar(&o.Service.TaskDefinitionArn, "task-definition-arn", "", "set task definition arn")
	cmd.Flags().BoolVar(&o.AutoApprove, "auto-approve", false, "approve deploy all")

	cmd.AddCommand(
		NewCmdDeployInfra(),
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

		for i := range o.Services {
			if len(o.Services[i].EcsCluster) == 0 {
				o.Services[i].EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
			}
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
		o.Config, err = config.InitializeConfig(config.WithDocker())
		viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("can`t complete options: %w", err)
		}
		o.Service.TaskDefinitionArn = viper.GetString("task-definition-arn")
		o.Service.EcsCluster = viper.GetString("ecs-cluster")
		o.Service.Path = viper.GetString("path")
		o.ServiceName = cmd.Flags().Args()[0]
		viper.UnmarshalKey(fmt.Sprintf("service.%s", o.ServiceName), &o.Service)
		if o.Service.EcsCluster == "" {
			o.Service.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
		}
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

func (o *DeployOptions) Run() error {
	if o.ServiceName == "" {
		err := deployAll(o)
		if err != nil {
			return err
		}
	} else {
		err := deployService(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func validate(o *DeployOptions) error {
	if o.Service.Image == "" {
		if o.Service.Path == "" {
			return fmt.Errorf("can't validate options: image or path must be specified")
		}
	} else {
		o.SkipBuildAndPush = true
	}

	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	if len(o.Service.EcsCluster) == 0 {
		return fmt.Errorf("can't validate options: ECS cluster must be specified")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified")
	}

	if len(o.ServiceName) == 0 {
		return fmt.Errorf("can't validate options: service name be specified")
	}

	return nil
}

func validateAll(o *DeployOptions) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified")
	}

	for sname, svc := range o.Services {
		if len(svc.Type) == 0 {
			return fmt.Errorf("can't validate options: type for service %s must be specified", sname)
		}

		if len(svc.Image) == 0 {
			if len(svc.Path) == 0 {
				return fmt.Errorf("can't validate options: image or path for service %s must be specified", sname)
			}
		}
	}

	return nil
}

func deployAll(o *DeployOptions) error {
	logrus.Infof("infra: %s", o.Infra)
	spinner := &pterm.SpinnerPrinter{}

	//terraform init run options
	opts := terraform.Options{
		ContainerName: "terraform",
		Cmd:           []string{"init", "-input=true"},
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Config.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Infra.Profile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
		TerraformVersion: o.Infra.Version,
	}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform init")
	}

	err := terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "init")
		return fmt.Errorf("can't deploy all: %w", err)
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform init completed")
	} else {
		pterm.Success.Println("terraform init completed")
	}

	//terraform plan run options
	opts.Cmd = []string{"plan"}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform plan")
	}

	err = terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "plan")
		return fmt.Errorf("can't deploy all: %w", err)
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform plan completed")
	} else {
		pterm.Success.Println("terraform plan completed")
	}

	//terraform apply run options
	opts.Cmd = []string{"apply", "-auto-approve"}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform apply")
	}

	err = terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "apply")
		return fmt.Errorf("can't deploy all: %w", err)
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform apply completed")
	} else {
		pterm.Success.Println("terraform apply completed")
	}

	//terraform output run options
	outputPath := fmt.Sprintf("%s/.terraform/output.json", viper.Get("ENV_DIR"))

	opts.Cmd = []string{"output", "-json"}
	opts.OutputPath = outputPath

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform output")
	}

	err = terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "output")
		return fmt.Errorf("can't deploy all: %w", err)
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform output completed")
	} else {
		pterm.Success.Println("terraform output completed")
	}

	name := fmt.Sprintf("/%s/terraform-output", o.Config.Env)

	outputFile, err := os.Open(outputPath)
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	defer outputFile.Close()

	byteValue, _ := ioutil.ReadAll(outputFile)
	sDec := base64.StdEncoding.EncodeToString(byteValue)
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	_, err = ssm.New(o.Config.Session).PutParameter(&ssm.PutParameterInput{
		Name:      &name,
		Value:     aws.String(string(sDec)),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(true),
		Tier:      aws.String("Intelligent-Tiering"),
		DataType:  aws.String("text"),
	})
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	logrus.Debug(o.Services)

	err = InDependencyOrder(aws.BackgroundContext(), &o.Services, func(c context.Context, name string) error {
		if logrus.GetLevel() < 4 {
			spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("deploy service %s", name))
		}

		o.Config.AwsProfile = o.Infra.Profile

		err = ecsdeploy.DeployService(
			o.Services[name],
			name,
			o.Tag,
			o.Config,
		)
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("can't deploy all: %w", err)
		}

		if logrus.GetLevel() < 4 {
			spinner.Success(fmt.Sprintf("deploy service %s completed", name))
		} else {
			pterm.Success.Printfln("deploy service %s completed", name)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func deployService(o *DeployOptions) error {
	spinner := &pterm.SpinnerPrinter{}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("deploy service %s", o.ServiceName))
	}

	err := ecsdeploy.DeployService(
		&o.Service,
		o.ServiceName,
		o.Tag,
		o.Config,
	)
	if err != nil {
		return fmt.Errorf("can't deploy: %w", err)
	}

	if logrus.GetLevel() < 4 {
		spinner.Success(fmt.Sprintf("deploy service %s completed", o.ServiceName))
	} else {
		pterm.Success.Printfln("deploy service %s completed", o.ServiceName)
	}

	return nil
}
