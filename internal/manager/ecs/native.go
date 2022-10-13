package ecs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pterm/pterm"
)

func (e *Manager) deployLocal(w io.Writer) error {
	pterm.SetDefaultOutput(w)

	svc := e.Project.AWSClient.ECSClient

	name := fmt.Sprintf("%s-%s", e.Project.Env, e.App.Name)

	dso, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  &e.App.Cluster,
		Services: []*string{&name},
	})
	if err != nil {
		return err
	}

	if len(dso.Services) == 0 {
		return fmt.Errorf("app %s not found not found in %s cluster", name, e.App.Cluster)
	}

	dtdo, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: dso.Services[0].TaskDefinition,
	})
	if err != nil {
		return err
	}

	definitions, err := svc.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: &name,
		Sort:         aws.String(ecs.SortOrderDesc),
	})
	if err != nil {
		return err
	}

	var oldTaskDef ecs.TaskDefinition
	var newTaskDef ecs.TaskDefinition

	if *dtdo.TaskDefinition.TaskDefinitionArn == *definitions.TaskDefinitionArns[0] {
		oldTaskDef = *dtdo.TaskDefinition
	} else {
		definition, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
			TaskDefinition: definitions.TaskDefinitionArns[0],
		})
		if err != nil {
			return err
		}

		oldTaskDef = *definition.TaskDefinition
	}

	pterm.Printfln("Deploying based on task definition: %s:%d", *oldTaskDef.Family, *oldTaskDef.Revision)

	var image string

	for i := 0; i < len(oldTaskDef.ContainerDefinitions); i++ {
		container := oldTaskDef.ContainerDefinitions[i]

		// We are changing the image/tag only for the app-specific container (not sidecars)
		if *container.Name == e.App.Name {
			if len(e.Project.Tag) != 0 && len(e.App.Image) == 0 {
				name := strings.Split(*container.Image, ":")[0]
				image = fmt.Sprintf("%s:%s", name, e.Project.Tag)
			} else {
				image = e.App.Image
			}

			pterm.Printfln(`Changed image of container "%s" to : "%s" (was: "%s")`, *container.Name, image, *container.Image)
			container.Image = &image
		}
	}

	pterm.Println("Creating new task definition revision")

	rtdo, err := svc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    oldTaskDef.ContainerDefinitions,
		Family:                  oldTaskDef.Family,
		Volumes:                 oldTaskDef.Volumes,
		TaskRoleArn:             oldTaskDef.TaskRoleArn,
		ExecutionRoleArn:        oldTaskDef.ExecutionRoleArn,
		RuntimePlatform:         oldTaskDef.RuntimePlatform,
		RequiresCompatibilities: oldTaskDef.RequiresCompatibilities,
		NetworkMode:             oldTaskDef.NetworkMode,
		Cpu:                     oldTaskDef.Cpu,
		Memory:                  oldTaskDef.Memory,
	})
	if err != nil {
		return err
	}

	newTaskDef = *rtdo.TaskDefinition

	pterm.Printfln("Successfully created revision: %s:%d", *rtdo.TaskDefinition.Family, *rtdo.TaskDefinition.Revision)

	if err = e.updateTaskDefinition(&newTaskDef, &oldTaskDef, name, "Deploying new task definition"); err != nil {
		err := e.getLastContainerLogs(fmt.Sprintf("%s-%s", e.Project.Env, e.App.Name))
		if err != nil {
			pterm.Println("Failed to get logs:", err)
		}

		sr, err := getStoppedReason(e.App.Cluster, name, svc)
		if err != nil {
			return err
		}

		pterm.Printfln("Container %s couldn't start: %s", name, sr)

		pterm.Printfln("Rolling back to old task definition: %s:%d", *oldTaskDef.Family, *oldTaskDef.Revision)
		e.App.Timeout = 600
		if err = e.updateTaskDefinition(&oldTaskDef, &newTaskDef, name, "Deploying previous task definition"); err != nil {
			return fmt.Errorf("unable to rollback to old task definition: %w", err)
		}

		pterm.Println("Rollback successful")

		return fmt.Errorf("deployment failed, but service has been rolled back to previous task definition: %s", *oldTaskDef.Family)
	}

	return nil
}

func (e *Manager) redeployLocal(w io.Writer) error {
	pterm.SetDefaultOutput(w)

	svc := e.Project.AWSClient.ECSClient

	name := fmt.Sprintf("%s-%s", e.Project.Env, e.App.Name)

	dso, err := getService(name, e.App.Cluster, svc)
	if err != nil {
		return err
	}

	var td *ecs.TaskDefinition

	switch e.App.TaskDefinitionRevision {
	case "latest":
		tds, err := svc.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
			FamilyPrefix: aws.String(name),
			Sort:         aws.String("DESC"),
		})
		if err != nil {
			return fmt.Errorf("unable to list task definitions: %w", err)
		}

		dtdo, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
			TaskDefinition: tds.TaskDefinitionArns[0],
		})
		if err != nil {
			return fmt.Errorf("unable to describe task definition: %w", err)
		}

		td = dtdo.TaskDefinition
	case "current":
		dtdo, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
			TaskDefinition: dso.Services[0].TaskDefinition,
		})
		if err != nil {
			return fmt.Errorf("unable to describe task definition: %w", err)
		}

		td = dtdo.TaskDefinition
	default:
		r, err := strconv.Atoi(e.App.TaskDefinitionRevision)
		if err == nil && r > 0 {
			arn := fmt.Sprintf("%s:%s", name, e.App.TaskDefinitionRevision)

			dtdo, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
				TaskDefinition: &arn,
			})
			if err != nil {
				return fmt.Errorf("unable to describe task definition: %w", err)
			}

			td = dtdo.TaskDefinition
		} else {
			return fmt.Errorf("invalid task definition revision: %s", e.App.TaskDefinitionRevision)
		}
	}

	if err = e.updateTaskDefinition(td, nil, name, "Redeploying new task definition"); err != nil {
		pterm.Println(err)
		err := e.getLastContainerLogs(fmt.Sprintf("%s-%s", e.Project.Env, e.App.Name))
		if err != nil {
			pterm.Println("Failed to get logs:", err)
		}
		return fmt.Errorf("redeployment failed")
	}

	return nil
}

func getService(name string, cluster string, svc ecsiface.ECSAPI) (*ecs.DescribeServicesOutput, error) {
	dso, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []*string{&name},
	})
	if err != nil {
		return nil, err
	}

	if len(dso.Services) == 0 {
		return nil, fmt.Errorf("app %s not found", name)
	}
	return dso, nil
}

func (e *Manager) updateTaskDefinition(newTD *ecs.TaskDefinition, oldTD *ecs.TaskDefinition, serviceName string, title string) error {
	pterm.Println("Updating service")

	svc := e.Project.AWSClient.ECSClient

	uso, err := svc.UpdateService(&ecs.UpdateServiceInput{
		Service:            aws.String(serviceName),
		Cluster:            aws.String(e.App.Cluster),
		TaskDefinition:     aws.String(*newTD.TaskDefinitionArn),
		ForceNewDeployment: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("unable to update service: %w", err)
	}

	var dtgo *elbv2.DescribeTargetGroupsOutput
	if e.App.Unsafe {
		elb := e.Project.AWSClient.ELBV2Client
		dtgo, err = elb.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
			TargetGroupArns: aws.StringSlice([]string{*uso.Service.LoadBalancers[0].TargetGroupArn}),
		})
		if err != nil {
			return fmt.Errorf("can't describe target groups: %w", err)
		}

		_, err = e.Project.AWSClient.ELBV2Client.ModifyTargetGroup(&elbv2.ModifyTargetGroupInput{
			HealthyThresholdCount:      aws.Int64(2),
			HealthCheckIntervalSeconds: aws.Int64(5),
			HealthCheckTimeoutSeconds:  aws.Int64(2),
			UnhealthyThresholdCount:    aws.Int64(2),
			TargetGroupArn:             uso.Service.LoadBalancers[0].TargetGroupArn,
		})
		if err != nil {
			return fmt.Errorf("unable to modify target group: %w", err)
		}
	}

	pterm.Printfln("Successfully changed task definition to: %s:%d", *newTD.Family, *newTD.Revision)
	pterm.Println(title)

	waitingTimeout := time.Now().Add(time.Duration(e.App.Timeout) * time.Second)
	waiting := true

	for waiting && time.Now().Before(waitingTimeout) {
		d, err := isDeployed(svc, serviceName, e.App.Cluster)
		if err != nil {
			return err
		}

		waiting = !d

		if waiting {
			time.Sleep(time.Second * 5)
		}
	}

	if waiting && time.Now().After(waitingTimeout) {
		pterm.Println("Deployment failed due to timeout")
		return fmt.Errorf("deployment failed due to timeout")
	}

	if e.App.Unsafe {
		_, err = e.Project.AWSClient.ELBV2Client.ModifyTargetGroup(&elbv2.ModifyTargetGroupInput{
			HealthyThresholdCount:      dtgo.TargetGroups[0].HealthyThresholdCount,
			HealthCheckIntervalSeconds: dtgo.TargetGroups[0].HealthCheckIntervalSeconds,
			HealthCheckTimeoutSeconds:  dtgo.TargetGroups[0].HealthCheckTimeoutSeconds,
			UnhealthyThresholdCount:    dtgo.TargetGroups[0].UnhealthyThresholdCount,
			TargetGroupArn:             uso.Service.LoadBalancers[0].TargetGroupArn,
		})
		if err != nil {
			return fmt.Errorf("unable to modify target group: %w", err)
		}
	}

	if oldTD != nil {
		if err = deregisterTaskDefinition(svc, oldTD); err != nil {
			return err
		}
	}

	return nil
}

func isDeployed(svc ecsiface.ECSAPI, name string, cluster string) (bool, error) {
	dso, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []*string{&name},
	})
	if err != nil {
		return false, err
	}

	if len(dso.Services) == 0 {
		return false, nil
	}

	if len(dso.Services[0].Deployments) != 1 {
		return false, nil
	}

	runningTasks, err := svc.ListTasks(&ecs.ListTasksInput{
		Cluster:     &cluster,
		ServiceName: &name,
	})
	if err != nil {
		return false, err
	}

	if len(runningTasks.TaskArns) == 0 {
		return *dso.Services[0].DesiredCount == 0, nil
	}

	runningCount, err := getRunningTaskCount(cluster, runningTasks.TaskArns, *dso.Services[0].TaskDefinition, svc)
	if err != nil {
		return false, err
	}

	return runningCount == *dso.Services[0].DesiredCount, nil
}

func getRunningTaskCount(cluster string, tasks []*string, serviceArn string, svc ecsiface.ECSAPI) (int64, error) {
	count := 0

	dto, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: &cluster,
		Tasks:   tasks,
	})
	if err != nil {
		return 0, err
	}

	for _, t := range dto.Tasks {
		if *t.TaskDefinitionArn == serviceArn && *t.LastStatus == "RUNNING" {
			count++
		}
	}

	return int64(count), nil
}

func getStoppedReason(cluster string, name string, svc ecsiface.ECSAPI) (string, error) {
	stopped := ecs.DesiredStatusStopped

	runningTasks, err := svc.ListTasks(&ecs.ListTasksInput{
		Cluster:       &cluster,
		ServiceName:   &name,
		DesiredStatus: &stopped,
	})
	if err != nil {
		return "", err
	}

	dto, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: &cluster,
		Tasks:   runningTasks.TaskArns,
	})
	if err != nil {
		return "", err
	}

	if dto.Tasks[0].StoppedReason == nil {
		return "", nil
	}

	return *dto.Tasks[0].StoppedReason, nil
}

func deregisterTaskDefinition(svc ecsiface.ECSAPI, td *ecs.TaskDefinition) error {
	pterm.Println("Deregister task definition revision")

	_, err := svc.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: td.TaskDefinitionArn,
	})
	if err != nil {
		return err
	}

	pterm.Printfln("Successfully deregistered revision: %s:%d", *td.Family, *td.Revision)

	return nil
}

func (e *Manager) getLastContainerLogs(logGroup string) error {
	cwl := e.Project.AWSClient.CloudWatchLogsClient
	out, err := cwl.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroup,
		Limit:        aws.Int64(1),
		Descending:   aws.Bool(true),
		OrderBy:      aws.String("LastEventTime"),
	})
	if err != nil {
		return err
	}

	if len(out.LogStreams) == 0 {
		return nil
	}

	pterm.Println("Container logs:")

	for _, stream := range out.LogStreams {
		out, err := cwl.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &logGroup,
			LogStreamName: stream.LogStreamName,
		})
		if err != nil {
			return err
		}

		for _, event := range out.Events {
			pterm.Println("| " + *event.Message)
		}
	}

	pterm.Println()

	return nil
}
