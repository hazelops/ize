package commands

import (
	"fmt"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ConsoleOptions struct {
	Config        *config.Project
	AppName       string
	EcsCluster    string
	Task          string
	CustomPrompt  bool
	ContainerName string
	Explain       bool
}

var explainConsoleTmpl = `
TASK_ID=$(aws ecs list-tasks --cluster {{.Env}}-{{.Namespace}} --service-name {{.Env}}-{{svc}} --desired-status "RUNNING" | jq -r '.taskArns[]' | cut -d'/' -f3 | head -n 1)

aws ecs execute-command  \
    --interactive \
    --region {{.AwsRegion}} \
    --cluster {{.Env}}-{{.Namespace}} \
    --task $TASK_ID \
    --container {{svc}} \
    --command "/bin/sh"
`

func NewConsoleFlags(project *config.Project) *ConsoleOptions {
	return &ConsoleOptions{
		Config: project,
	}
}

func NewCmdConsole(project *config.Project) *cobra.Command {
	o := NewConsoleFlags(project)

	cmd := &cobra.Command{
		Use:               "console [app-name]",
		Short:             "Connect to a container in the ECS",
		Long:              "Connect to a container of the app via AWS SSM.\nTakes app name that is running on ECS as an argument",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete(cmd)
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

	cmd.Flags().StringVar(&o.EcsCluster, "ecs-cluster", "", "set ECS cluster name")
	cmd.Flags().StringVar(&o.ContainerName, "container-name", "", "set container name")
	cmd.Flags().StringVar(&o.Task, "task", "", "set task id")
	cmd.Flags().BoolVar(&o.Explain, "explain", false, "bash alternative shown")
	cmd.Flags().BoolVar(&o.CustomPrompt, "custom-prompt", false, "enable custom prompt in the console")

	return cmd
}

func (o *ConsoleOptions) Complete(cmd *cobra.Command) error {
	if err := requirements.CheckRequirements(requirements.WithSSMPlugin()); err != nil {
		return err
	}

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	if !o.CustomPrompt {
		o.CustomPrompt = o.Config.CustomPrompt
	}

	o.AppName = cmd.Flags().Args()[0]

	if len(o.ContainerName) == 0 {
		o.ContainerName = o.AppName
	}

	return nil
}

func (o *ConsoleOptions) Validate() error {
	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified")
	}

	return nil
}

func (o *ConsoleOptions) Run() error {
	appName := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)

	if o.Explain {
		err := o.Config.Generate(explainConsoleTmpl, template.FuncMap{
			"svc": func() string {
				return o.AppName
			},
		})
		if err != nil {
			return err
		}

		return nil
	}

	logrus.Infof("app name: %s, cluster name: %s", appName, o.EcsCluster)
	logrus.Infof("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	s, _ := pterm.DefaultSpinner.WithRemoveWhenDone().Start("Getting access to container...")

	if o.Task == "" {
		lto, err := o.Config.AWSClient.ECSClient.ListTasks(&ecs.ListTasksInput{
			Cluster:       &o.EcsCluster,
			DesiredStatus: aws.String(ecs.DesiredStatusRunning),
			ServiceName:   &appName,
		})
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ClusterNotFoundException":
				return fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
			default:
				return err
			}
		}

		logrus.Debugf("list task output: %s", lto)

		if len(lto.TaskArns) == 0 {
			return fmt.Errorf("running task not found")
		}

		o.Task = *lto.TaskArns[0]
	}

	s.UpdateText("Executing command...")
	consoleCommand := `/bin/sh`

	if o.CustomPrompt {
		// This is ASCII Prompt string with colors. See https://dev.to/ifenna__/adding-colors-to-bash-scripts-48g4 for reference
		// TODO: Make this customizable via a config
		promptString := fmt.Sprintf(`\e[1;35m★\e[0m $ENV-$APP_NAME\n\e[1;33m\e[0m \w \e[1;34m❯\e[0m `)
		consoleCommand = fmt.Sprintf(`/bin/sh -c '$(echo "export PS1=\"%s\"" > /etc/profile.d/ize.sh) /bin/bash --rcfile /etc/profile'`, promptString)
	}

	out, err := o.Config.AWSClient.ECSClient.ExecuteCommand(&ecs.ExecuteCommandInput{
		Container:   &o.ContainerName,
		Interactive: aws.Bool(true),
		Cluster:     &o.EcsCluster,
		Task:        &o.Task,
		Command:     aws.String(consoleCommand),
	})
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "ClusterNotFoundException":
			return fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
		default:
			return err
		}
	}

	s.Success()

	ssmCmd := ssmsession.NewSSMPluginCommand(o.Config.AwsRegion)
	err = ssmCmd.StartInteractive(out.Session)
	if err != nil {
		return err
	}

	return nil
}
