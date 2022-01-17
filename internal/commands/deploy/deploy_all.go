package deploy

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeployAllOptions struct {
	Version          string
	Env              string
	Namespace        string
	Profile          string
	Region           string
	Tag              string
	SkipBuildAndPush bool
	Services         []Service
	Infra            Infra
}

type Service struct {
	Name              string
	Type              string
	Path              string
	Image             string
	EcsCluster        string
	TaskDefinitionArn string
}

type Infra struct {
	Version string `mapstructure:"terraform_version"`
	Region  string `mapstructure:"aws_region"`
	Profile string `mapstructure:"aws_profile"`
}

var deployAllLongDesc = templates.LongDesc(`
	Deploy all infraftructures and sevices.
	Deployment is recommended through the configuration file.
`)

var deployAllExample = templates.Examples(`
	# Deploy all with set infra option
	ize deploy all -- infra.terraform.terraform_version=1.0.10 infra.terraform.aws_region=us-east-1 infra.terraform.aws_profile=default

	# Deploy all with set service option
	ize deploy all -- service.goblin.path=path/to/service service.goblin.type=ecs

	# Deploy all via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize deploy all
`)

func NewDeployAllFlags() *DeployAllOptions {
	return &DeployAllOptions{}
}

func NewCmdDeployAll() *cobra.Command {
	o := NewDeployAllFlags()

	cmd := &cobra.Command{
		Use:     "all",
		Example: deployAllExample,
		Short:   "manage deployments",
		Long:    deployAllLongDesc,
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

	cmd.AddCommand(NewCmdDeployInfra())

	return cmd
}

func setServiceFromArgs(args []string) error {
	for _, a := range args {
		if !strings.Contains(a, "=") {
			return fmt.Errorf("wrong format of argument %s (should be key=value)", a)
		}
		parts := strings.Split(a, "=")
		viper.Set(parts[0], parts[1])
	}

	return nil
}

func getServices(rawServices interface{}) []Service {
	var services []Service

	for name := range rawServices.(map[string]interface{}) {
		s := Service{}
		viper.UnmarshalKey(fmt.Sprintf("service.%s", name), &s)
		s.Name = name
		services = append(services, s)
	}

	return services
}

func (o *DeployAllOptions) Complete(cmd *cobra.Command, args []string) error {
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	argsLenAtDash := cmd.ArgsLenAtDash()

	if argsLenAtDash != -1 {
		err = setServiceFromArgs(args[argsLenAtDash:])
		if err != nil {
			return err
		}
	}

	o.Services = getServices(viper.AllSettings()["service"])
	viper.UnmarshalKey("infra.terraform", &o.Infra)

	viper.BindPFlags(cmd.Flags())
	o.Env = viper.GetString("env")
	o.Namespace = viper.GetString("namespace")
	o.Region = viper.GetString("aws-region")
	o.Profile = viper.GetString("aws-profile")
	o.Tag = viper.GetString("tag")
	o.Version = viper.GetString("terraform-version")

	//Set global options if infra options are empty
	if len(o.Infra.Version) == 0 {
		o.Infra.Version = o.Version
	}
	if len(o.Infra.Region) == 0 {
		o.Infra.Version = o.Region
	}
	if len(o.Infra.Profile) == 0 {
		o.Infra.Profile = o.Profile
	}

	if len(o.Profile) == 0 {
		return fmt.Errorf("profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("region must be specified")
	}

	return nil
}

func (o *DeployAllOptions) Validate() error {
	if len(o.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	if len(o.Profile) == 0 {
		return fmt.Errorf("can't validate options: profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("can't validate options: region must be specified")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified")
	}

	for i, _ := range o.Services {
		if len(o.Services[i].Type) == 0 {
			return fmt.Errorf("can't validate options: type for service %s must be specified", o.Services[i].Name)
		}

		if len(o.Services[i].Image) == 0 {
			if len(o.Services[i].Path) == 0 {
				return fmt.Errorf("can't validate options: image or path for service %s must be specified", o.Services[i].Name)
			}
		}

		if len(o.Services[i].EcsCluster) == 0 {
			o.Services[i].EcsCluster = fmt.Sprintf("%s-%s", o.Env, o.Namespace)
		}
	}

	return nil
}

func (o *DeployAllOptions) Run() error {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Infra.Region,
		Profile: o.Infra.Profile,
	})
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	logrus.Infof("infra: %s", o.Infra)
	spinner := &pterm.SpinnerPrinter{}

	//terraform init run options
	opts := terraform.Options{
		ContainerName: "terraform",
		Cmd:           []string{"init", "-input=true"},
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Infra.Profile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
		TerraformVersion: o.Infra.Version,
	}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform init")
	}

	err = terraform.Run(opts)
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

	name := fmt.Sprintf("/%s/terraform-output", o.Env)

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

	_, err = ssm.New(sess).PutParameter(&ssm.PutParameterInput{
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

	for _, svc := range o.Services {
		if logrus.GetLevel() < 4 {
			spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("deploy service %s", svc.Name))
		}

		err = deployService(
			svc,
			o.Infra.Profile,
			o.Namespace,
			o.Env,
			o.Tag,
			sess,
		)
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("can't deploy all: %w", err)
		}

		if logrus.GetLevel() < 4 {
			spinner.Success(fmt.Sprintf("deploy service %s completed", svc.Name))
		} else {
			pterm.Success.Printfln("deploy service %s completed", svc.Name)
		}
	}

	return nil
}
