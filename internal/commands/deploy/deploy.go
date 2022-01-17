package deploy

import (
	"fmt"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/ecsdeploy"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeployOptions struct {
	ServiceName       string
	EcsCluster        string
	TaskDefinitionArn string
	Image             string
	Env               string
	Namespace         string
	Path              string
	Profile           string
	Region            string
	Tag               string
	SkipBuildAndPush  bool
}

var deployLongDesc = templates.LongDesc(`
	Deploy infraftructure or sevice.
	For deploy service the service name must be specimfied. 
	The infrastructure for the service must be prepared in advance.
`)

var deployExample = templates.Examples(`
	# Deploy service form image
	ize deploy --image foo/bar:latest <service name>

	# Deploy service form path
	ize deploy --path /path/to/service <service name>

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
		Args:    cobra.MinimumNArgs(1),
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

	cmd.Flags().StringVar(&o.Image, "image", "", "set image name")
	cmd.Flags().StringVar(&o.EcsCluster, "ecs-cluster", "", "set ECS cluster name")
	cmd.Flags().StringVar(&o.TaskDefinitionArn, "task-definition-arn", "", "set task definition arn")

	cmd.AddCommand(
		NewCmdDeployAll(),
		NewCmdDeployInfra(),
	)

	return cmd
}

func (o *DeployOptions) Complete(cmd *cobra.Command, args []string) error {
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	viper.BindPFlags(cmd.Flags())
	o.Env = viper.GetString("env")
	o.Namespace = viper.GetString("namespace")
	o.TaskDefinitionArn = viper.GetString("task-definition-arn")
	o.EcsCluster = viper.GetString("ecs-cluster")
	o.Region = viper.GetString("aws-region")
	o.Profile = viper.GetString("aws-profile")
	o.Tag = viper.GetString("tag")
	o.Tag = viper.GetString("image")

	o.ServiceName = cmd.Flags().Args()[0]

	var cfg ecsServiceConfig

	err = mapstructure.Decode(viper.GetStringMap(fmt.Sprintf("service.%s", o.ServiceName)), &cfg)
	if err != nil {
		return err
	}

	o.Path = cfg.Path

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Env, o.Namespace)
	}

	return nil
}

func (o *DeployOptions) Validate() error {
	if o.Image == "" {
		if o.Path == "" {
			return fmt.Errorf("image or path must be specified")
		}
	} else {
		o.SkipBuildAndPush = true
	}

	if len(o.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	if len(o.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified")
	}

	if len(o.EcsCluster) == 0 {
		return fmt.Errorf("ECS cluster must be specified")
	}

	if len(o.Profile) == 0 {
		return fmt.Errorf("profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("region must be specified")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("tag must be specified")
	}

	if len(o.ServiceName) == 0 {
		return fmt.Errorf("service name be specified")
	}
	return nil
}

func (o *DeployOptions) Run() error {
	logrus.Debugf("profile: %s, region: %s", o.Profile, o.Region)

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Region,
		Profile: o.Profile,
	})
	if err != nil {
		return fmt.Errorf("can't deploy : %w", err)
	}

	err = ecsdeploy.DeployService(
		ecsdeploy.Service{Name: o.ServiceName, Path: o.Path, Image: o.Image, EcsCluster: o.EcsCluster, TaskDefinitionArn: o.TaskDefinitionArn},
		o.Profile,
		o.Namespace,
		o.Env,
		o.Tag,
		sess,
	)
	if err != nil {
		return fmt.Errorf("can't deploy: %w", err)
	}

	return nil
}

type ecsServiceConfig struct {
	Path string
}
