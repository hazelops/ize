package commands

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ExecOptions struct {
	Config        *config.Project
	AppName       string
	EcsCluster    string
	Command       []string
	Task          string
	ContainerName string
	Explain       bool
}

var explainExecTmpl = `
TASK_ID=$(aws ecs list-tasks --cluster {{.Env}}-{{.Namespace}} --service-name {{.Env}}-{{svc}} --desired-status "RUNNING" | jq -r '.taskArns[]' | cut -d'/' -f3 | head -n 1)

aws ecs execute-command  \
    --interactive \
    --region {{.AwsRegion}} \
    --cluster {{.Env}}-{{.Namespace}} \
    --task $TASK_ID \
    --container {{svc}} \
    --command {{command}}
`

var execExample = templates.Examples(`
	# Connect to a container in the ECS via AWS SSM and run command.
	ize exec goblin ps aux
`)

func NewExecFlags(project *config.Project) *ExecOptions {
	return &ExecOptions{
		Config: project,
	}
}

func NewCmdExec(project *config.Project) *cobra.Command {
	o := NewExecFlags(project)

	cmd := &cobra.Command{
		Use:               "exec [app-name] -- [commands]",
		Example:           execExample,
		Short:             "Execute command in ECS container",
		Long:              "Connect to a container in the ECS via AWS SSM and run command.\nIt uses app name as an argument.",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			argsLenAtDash := cmd.ArgsLenAtDash()
			err := o.Complete(cmd, args, argsLenAtDash)
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
	cmd.Flags().StringVar(&o.Task, "task", "", "set task id")
	cmd.Flags().StringVar(&o.ContainerName, "container-name", "", "set container name")
	cmd.Flags().BoolVar(&o.Explain, "explain", false, "bash alternative shown")

	return cmd
}

func (o *ExecOptions) Complete(cmd *cobra.Command, args []string, argsLenAtDash int) error {
	if err := requirements.CheckRequirements(requirements.WithSSMPlugin()); err != nil {
		return err
	}

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	o.AppName = cmd.Flags().Args()[0]

	if len(o.ContainerName) == 0 {
		o.ContainerName = o.AppName
	}

	if argsLenAtDash > -1 {
		o.Command = args[argsLenAtDash:]
	}

	return nil
}

func (o *ExecOptions) Validate() error {
	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified")
	}

	if len(o.Command) == 0 {
		return fmt.Errorf("can't validate: you must specify at least one command for the container")
	}

	return nil
}

func (o *ExecOptions) Run() error {
	appName := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)

	if o.Explain {
		err := o.Config.Generate(explainExecTmpl, template.FuncMap{
			"svc": func() string {
				return o.AppName
			},
			"command": func() string {
				return strings.Join(o.Command, " ")
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

	out, err := o.Config.AWSClient.ECSClient.ExecuteCommand(&ecs.ExecuteCommandInput{
		Container:   &o.AppName,
		Interactive: aws.Bool(true),
		Cluster:     &o.EcsCluster,
		Task:        &o.Task,
		Command:     aws.String(strings.Join(o.Command, " ")),
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
	err = ssmCmd.Start(out.Session)
	if err != nil {
		return err
	}

	return nil
}
