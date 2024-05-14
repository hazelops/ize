package commands

import (
	"errors"
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
	Config           *config.Project
	AppName          string
	EcsCluster       string
	EcsServiceName   string
	Task             string
	CustomPrompt     bool
	EcsContainerName string
	Explain          bool
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
	cmd.Flags().StringVar(&o.EcsContainerName, "container-name", "", "set container name")
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

	return nil
}

func (o *ConsoleOptions) Validate() error {
	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified")
	}

	return nil
}

func (o *ConsoleOptions) Run() error {
	var err error

	if len(o.EcsServiceName) == 0 {
		o.EcsServiceName, err = getEcsServiceName(o)
		if err != nil {
			return err
		}
	}

	if o.Explain {
		err := o.Config.Generate(explainConsoleTmpl, template.FuncMap{
			"svc": func() string {
				return o.EcsServiceName
			},
		})
		if err != nil {
			return err
		}

		return nil
	}

	logrus.Infof("app name: %s, cluster name: %s", o.EcsServiceName, o.EcsCluster)
	logrus.Infof("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	s, _ := pterm.DefaultSpinner.WithRemoveWhenDone().Start("Getting access to container...")

	if o.Task == "" {
		// Infer task name from the app name
		lto, err := o.Config.AWSClient.ECSClient.ListTasks(&ecs.ListTasksInput{
			Cluster:       &o.EcsCluster,
			DesiredStatus: aws.String(ecs.DesiredStatusRunning),
			ServiceName:   &o.EcsServiceName,
		})

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ClusterNotFoundException":
				return fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
			default:
				{
					return err
				}
			}
		}

		logrus.Debugf("list task output: %s", lto)

		if len(lto.TaskArns) == 0 {
			return fmt.Errorf("running task not found")
		}

		o.Task = *lto.TaskArns[0]
	}

	if len(o.EcsContainerName) == 0 {
		o.EcsContainerName, err = getEcsContainerName(o)
		if err != nil {
			return err
		}
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
		Container:   &o.EcsContainerName,
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

func getEcsServiceName(o *ConsoleOptions) (string, error) {
	// TODO: Move core logic to a shared function (since it's used in deploy too)
	ecsServiceCandidates := []string{
		o.AppName,
		fmt.Sprintf("%s-%s", o.Config.Env, o.AppName),
		fmt.Sprintf("%s-%s-%s", o.Config.Env, o.Config.Namespace, o.AppName),
	}

	for _, v := range ecsServiceCandidates {
		logrus.Debugf("Checking if ECS service %s exists in cluster %s.", v, o.EcsCluster)
		_, err := o.Config.AWSClient.ECSClient.ListTasks(&ecs.ListTasksInput{
			Cluster:       &o.EcsCluster,
			DesiredStatus: aws.String(ecs.DesiredStatusRunning),
			ServiceName:   &v,
		})

		var aerr awserr.Error
		if errors.As(err, &aerr) {
			switch aerr.Code() {
			case "ClusterNotFoundException":
				return "", fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
			case "ServiceNotFoundException":
				{
					logrus.Infof("ECS Service not found: %s in cluster %s.", v, o.EcsCluster)
					continue
				}
			default:
				{
					return "", err
				}
			}

		}
		return v, err
	}
	err := errors.New("ECS Service not found")
	return "", err
}

func getEcsContainerName(o *ConsoleOptions) (string, error) {
	ecsContainerNameCandidates := []string{
		o.AppName,
		fmt.Sprintf("%s-%s", o.Config.Namespace, o.AppName),
		fmt.Sprintf("%s-%s", o.Config.Env, o.AppName),
		fmt.Sprintf("%s-%s-%s", o.Config.Env, o.Config.Namespace, o.AppName),
	}

	for _, v := range ecsContainerNameCandidates {
		logrus.Debugf("Checking if ECS container %s exists in task %s.", v, o.Task)
		t, err := o.Config.AWSClient.ECSClient.ListTasks(&ecs.ListTasksInput{
			Cluster:       &o.EcsCluster,
			DesiredStatus: aws.String(ecs.DesiredStatusRunning),
			ServiceName:   &o.EcsServiceName,
		})

		if err != nil {
			return "", err
		}

		if len(t.TaskArns) > 0 {
			tasks, err := o.Config.AWSClient.ECSClient.DescribeTasks(&ecs.DescribeTasksInput{
				Cluster: &o.EcsCluster,
				Tasks:   t.TaskArns,
			})

			if err != nil {
				return "", err
			}

			if len(tasks.Tasks) > 0 {
				for _, task := range tasks.Tasks {
					logrus.Debugf("Task arn is %s", *task.TaskArn)

					for _, container := range task.Containers {
						for _, ecsContainerNameCandidate := range ecsContainerNameCandidates {
							logrus.Debugf("Checking if %s==%s", ecsContainerNameCandidate, *container.Name)
							if ecsContainerNameCandidate == *container.Name {
								return *container.Name, nil
							} else {
								continue
							}
						}
					}

					return "", errors.New(fmt.Sprintf("Can't find a container for %s in %s", o.AppName, task.TaskDefinitionArn))
				}
			} else {
				fmt.Println("No tasks found.")
			}

		}

		//_, err := o.Config.AWSClient.ECSClient.ListContainerInstances(&ecs.ListContainerInstancesInput{
		//	Cluster: &o.EcsCluster,
		//	//Filter:  "",
		//})

		var aerr awserr.Error
		if errors.As(err, &aerr) {
			switch aerr.Code() {
			case "ClusterNotFoundException":
				return "", fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
			case "ServiceNotFoundException":
				{
					logrus.Infof("ECS Service not found: %s in cluster %s.", v, o.EcsCluster)
					continue
				}
			default:
				{
					return "", err
				}
			}

		}
		return v, err
	}
	err := errors.New(fmt.Sprintf("ECS Container for %s not found", o.AppName))
	return "", err
}
