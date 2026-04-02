package ecscron

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
)

// RunNow manually triggers the cron task outside its schedule via ECS RunTask.
func (m *Manager) RunNow(ui terminal.UI) error {
	m.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	if err := m.ensureSession(); err != nil {
		return err
	}

	s := sg.Add("%s: starting cron task manually...", m.App.Name)
	defer func() { s.Abort() }()

	ecsSvc := m.Project.AWSClient.ECSClient
	family := m.taskFamily()

	nc, err := getNetworkConfigFromSSM(m.Project.AWSClient.SSMClient, m.Project.Env)
	if err != nil {
		return fmt.Errorf("can't get network configuration: %w", err)
	}

	if len(nc.VpcPrivateSubnets.Value) == 0 {
		return fmt.Errorf("vpc_private_subnets is empty in terraform output")
	}

	out, err := ecsSvc.RunTaskWithContext(context.Background(), &ecs.RunTaskInput{
		TaskDefinition: &family,
		StartedBy:      aws.String("IZE-cron-manual"),
		Cluster:        &m.App.Cluster,
		LaunchType:     aws.String(ecs.LaunchTypeFargate),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets: aws.StringSlice(nc.VpcPrivateSubnets.Value),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("can't run task: %w", err)
	}

	if len(out.Tasks) == 0 {
		return fmt.Errorf("no tasks started")
	}

	taskArn := *out.Tasks[0].TaskArn
	logrus.Debugf("started task: %s", taskArn)

	s.Done()
	pterm.Success.Printfln("Task %s started: %s", m.App.Name, taskArn)
	return nil
}
