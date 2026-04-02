package ecscron

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	ecsmanager "github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	Project *config.Project
	App     *config.EcsCron
}

func (m *Manager) prepare() {
	if m.App.Path == "" {
		appsPath := m.Project.AppsPath
		if !filepath.IsAbs(appsPath) {
			appsPath = filepath.Join(os.Getenv("PWD"), appsPath)
		}
		m.App.Path = filepath.Join(appsPath, m.App.Name)
	} else if !filepath.IsAbs(m.App.Path) {
		m.App.Path = filepath.Join(m.Project.RootDir, m.App.Path)
	}

	if m.App.Cluster == "" {
		m.App.Cluster = fmt.Sprintf("%s-%s", m.Project.Env, m.Project.Namespace)
	}

	if m.App.DockerRegistry == "" {
		m.App.DockerRegistry = m.Project.DockerRegistry
	}

	if m.App.Timeout == 0 {
		m.App.Timeout = 300
	}
}

func (m *Manager) ruleName() string {
	return fmt.Sprintf("%s-%s-cron", m.Project.Env, m.App.Name)
}

func (m *Manager) taskFamily() string {
	if m.App.TaskDefinition != "" {
		return m.App.TaskDefinition
	}
	return fmt.Sprintf("%s-%s", m.Project.Env, m.App.Name)
}

func (m *Manager) ensureSession() error {
	if m.App.AwsRegion != "" && m.App.AwsProfile != "" {
		sess, err := utils.GetSession(&utils.SessionConfig{
			Region:  m.App.AwsRegion,
			Profile: m.App.AwsProfile,
		})
		if err != nil {
			return fmt.Errorf("can't get session: %w", err)
		}
		m.Project.SettingAWSClient(sess)
	}
	return nil
}

func (m *Manager) toEcsApp() *config.Ecs {
	return &config.Ecs{
		Name:           m.App.Name,
		Path:           m.App.Path,
		Image:          m.App.Image,
		Cluster:        m.App.Cluster,
		DockerRegistry: m.App.DockerRegistry,
		Timeout:        m.App.Timeout,
		SkipDeploy:     m.App.SkipDeploy,
		Icon:           m.App.Icon,
		AwsProfile:     m.App.AwsProfile,
		AwsRegion:      m.App.AwsRegion,
		DependsOn:      m.App.DependsOn,
	}
}

// Deploy creates or updates an EventBridge rule and its ECS task target.
func (m *Manager) Deploy(ui terminal.UI) error {
	m.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	if err := m.ensureSession(); err != nil {
		return err
	}

	if m.App.SkipDeploy {
		s := sg.Add("%s: deploy will be skipped", m.App.Name)
		defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()
		s.Done()
		return nil
	}

	s := sg.Add("%s: deploying cron task...", m.App.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	// Resolve image
	image := m.App.Image
	if image == "" {
		image = fmt.Sprintf("%s/%s:%s",
			m.App.DockerRegistry,
			fmt.Sprintf("%s-%s", m.Project.Namespace, m.App.Name),
			fmt.Sprintf("%s-%s", m.Project.Env, "latest"))
	}

	ecsSvc := m.Project.AWSClient.ECSClient
	ebSvc := m.Project.AWSClient.EventBridgeClient
	family := m.taskFamily()

	// 1. Get current task definition and register new revision with updated image
	latestTd, err := getLatestTaskDefinition(ecsSvc, family)
	if err != nil {
		return fmt.Errorf("task definition %s not found — it must be created via Terraform first: %w", family, err)
	}

	for _, cd := range latestTd.ContainerDefinitions {
		if *cd.Name == m.App.Name {
			pterm.Printfln(`Changed image of container "%s" to: "%s" (was: "%s")`, *cd.Name, image, *cd.Image)
			cd.Image = &image
		}
	}

	rtdo, err := ecsSvc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    latestTd.ContainerDefinitions,
		Family:                  latestTd.Family,
		Volumes:                 latestTd.Volumes,
		TaskRoleArn:             latestTd.TaskRoleArn,
		ExecutionRoleArn:        latestTd.ExecutionRoleArn,
		RuntimePlatform:         latestTd.RuntimePlatform,
		RequiresCompatibilities: latestTd.RequiresCompatibilities,
		NetworkMode:             latestTd.NetworkMode,
		Cpu:                     latestTd.Cpu,
		Memory:                  latestTd.Memory,
	})
	if err != nil {
		return fmt.Errorf("can't register task definition: %w", err)
	}

	newTdArn := *rtdo.TaskDefinition.TaskDefinitionArn
	pterm.Printfln("Registered task definition: %s:%d", *rtdo.TaskDefinition.Family, *rtdo.TaskDefinition.Revision)

	// 2. Create or update EventBridge rule
	// Three scenarios:
	//   a) schedule in config -> create/update rule with that schedule
	//   b) no schedule, rule exists -> keep existing rule (only update target)
	//   c) no schedule, rule doesn't exist -> error
	ruleName := m.ruleName()

	if m.App.Schedule != "" {
		// Scenario (a): schedule provided — create or update the rule
		state := "ENABLED"
		if m.App.Enabled != nil && !*m.App.Enabled {
			state = "DISABLED"
		}

		_, err = ebSvc.PutRule(&eventbridge.PutRuleInput{
			Name:               &ruleName,
			ScheduleExpression: &m.App.Schedule,
			State:              &state,
			Description:        aws.String(fmt.Sprintf("IZE cron: %s", m.App.Name)),
		})
		if err != nil {
			return fmt.Errorf("can't create/update EventBridge rule: %w", err)
		}

		pterm.Printfln("EventBridge rule %s: schedule=%s state=%s", ruleName, m.App.Schedule, state)
	} else {
		// Scenario (b)/(c): no schedule — check if rule already exists
		dro, err := ebSvc.DescribeRule(&eventbridge.DescribeRuleInput{
			Name: &ruleName,
		})
		if err != nil {
			return fmt.Errorf("schedule not specified and EventBridge rule %s not found — schedule is required to create a new rule: %w", ruleName, err)
		}

		pterm.Printfln("EventBridge rule %s exists (schedule=%s), updating target only", ruleName, *dro.ScheduleExpression)
	}

	// 3. Get role ARN from task definition execution role
	roleArn := latestTd.ExecutionRoleArn
	if roleArn == nil {
		return fmt.Errorf("task definition %s has no execution role — required for EventBridge target", family)
	}

	// 4. Get network configuration from SSM
	nc, err := getNetworkConfigFromSSM(m.Project.AWSClient.SSMClient, m.Project.Env)
	if err != nil {
		return fmt.Errorf("can't get network configuration: %w", err)
	}

	if len(nc.VpcPrivateSubnets.Value) == 0 {
		return fmt.Errorf("vpc_private_subnets is empty in terraform output")
	}

	// 5. Get account ID for cluster ARN
	accountId, err := m.getAccountId()
	if err != nil {
		return fmt.Errorf("can't get AWS account ID: %w", err)
	}

	// 6. Put target
	_, err = ebSvc.PutTargets(&eventbridge.PutTargetsInput{
		Rule: &ruleName,
		Targets: []*eventbridge.Target{
			{
				Id:      aws.String(m.App.Name),
				Arn:     aws.String(fmt.Sprintf("arn:aws:ecs:%s:%s:cluster/%s", m.Project.AwsRegion, accountId, m.App.Cluster)),
				RoleArn: roleArn,
				EcsParameters: &eventbridge.EcsParameters{
					TaskDefinitionArn: &newTdArn,
					TaskCount:         aws.Int64(1),
					LaunchType:        aws.String("FARGATE"),
					NetworkConfiguration: &eventbridge.NetworkConfiguration{
						AwsvpcConfiguration: &eventbridge.AwsVpcConfiguration{
							Subnets: aws.StringSlice(nc.VpcPrivateSubnets.Value),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("can't put EventBridge target: %w", err)
	}

	s.Done()
	s = sg.Add("%s: cron task deployed! schedule=%s", m.App.Name, m.App.Schedule)
	s.Done()
	return nil
}

// Destroy removes the EventBridge rule, its targets, and deregisters task definitions.
func (m *Manager) Destroy(ui terminal.UI, autoApprove bool) error {
	m.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	if err := m.ensureSession(); err != nil {
		return err
	}

	s := sg.Add("%s: destroying cron task...", m.App.Name)
	defer func() { s.Abort(); time.Sleep(200 * time.Millisecond) }()

	ebSvc := m.Project.AWSClient.EventBridgeClient
	ecsSvc := m.Project.AWSClient.ECSClient
	ruleName := m.ruleName()

	// 1. Remove all targets from the rule
	lto, err := ebSvc.ListTargetsByRule(&eventbridge.ListTargetsByRuleInput{
		Rule: &ruleName,
	})
	if err != nil {
		// Rule may not exist — treat as already destroyed
		logrus.Debugf("rule %s may not exist: %v", ruleName, err)
	} else if len(lto.Targets) > 0 {
		var ids []*string
		for _, t := range lto.Targets {
			ids = append(ids, t.Id)
		}
		_, err = ebSvc.RemoveTargets(&eventbridge.RemoveTargetsInput{
			Rule: &ruleName,
			Ids:  ids,
		})
		if err != nil {
			return fmt.Errorf("can't remove targets from rule %s: %w", ruleName, err)
		}
	}

	// 2. Delete the rule (ignore ResourceNotFoundException)
	_, err = ebSvc.DeleteRule(&eventbridge.DeleteRuleInput{
		Name: &ruleName,
	})
	if err != nil {
		logrus.Debugf("can't delete rule %s (may not exist): %v", ruleName, err)
	}

	// 3. Deregister task definitions
	family := m.taskFamily()
	definitions, err := ecsSvc.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: &family,
		Sort:         aws.String(ecs.SortOrderDesc),
	})
	if err == nil {
		for _, tda := range definitions.TaskDefinitionArns {
			_, _ = ecsSvc.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{
				TaskDefinition: tda,
			})
		}
	}

	s.Done()
	s = sg.Add("%s: cron task destroyed!", m.App.Name)
	s.Done()
	return nil
}

// Build delegates to the existing ECS manager (same Docker build pipeline).
func (m *Manager) Build(ui terminal.UI) error {
	m.prepare()
	em := &ecsmanager.Manager{Project: m.Project, App: m.toEcsApp()}
	return em.Build(ui)
}

// Push delegates to the existing ECS manager (same ECR push pipeline).
func (m *Manager) Push(ui terminal.UI) error {
	m.prepare()
	em := &ecsmanager.Manager{Project: m.Project, App: m.toEcsApp()}
	return em.Push(ui)
}

// Redeploy re-deploys with the current schedule (same as Deploy for cron tasks).
func (m *Manager) Redeploy(ui terminal.UI) error {
	return m.Deploy(ui)
}

func getLatestTaskDefinition(svc ecsiface.ECSAPI, family string) (*ecs.TaskDefinition, error) {
	tds, err := svc.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: &family,
		Sort:         aws.String("DESC"),
	})
	if err != nil {
		return nil, err
	}
	if len(tds.TaskDefinitionArns) == 0 {
		return nil, fmt.Errorf("no task definitions found for family %s", family)
	}

	dtdo, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: tds.TaskDefinitionArns[0],
	})
	if err != nil {
		return nil, err
	}
	return dtdo.TaskDefinition, nil
}

func (m *Manager) getAccountId() (string, error) {
	out, err := m.Project.AWSClient.STSClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("can't get caller identity: %w", err)
	}
	return *out.Account, nil
}
