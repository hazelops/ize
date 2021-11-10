package commands

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/pkg/ssmsession.go"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type sshCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newSSHCmd() *sshCmd {
	cc := &sshCmd{}

	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "",
		Long:  "",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			serviceName := fmt.Sprintf("%s-%s", cc.config.Env, os.Args[2])
			clusterName := fmt.Sprintf("%s-%s", cc.config.Env, cc.config.Namespace)

			sess, err := utils.GetSession(&utils.SessionConfig{
				Region:  cc.config.AwsRegion,
				Profile: cc.config.AwsProfile,
			})
			if err != nil {
				pterm.Error.Printfln("Getting AWS session")
				return err
			}

			pterm.Success.Printfln("Getting AWS session")

			ecsSvc := ecs.New(sess)

			lto, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
				Cluster:       &clusterName,
				DesiredStatus: aws.String(ecs.DesiredStatusRunning),
				ServiceName:   &serviceName,
			})
			if err != nil {
				pterm.Error.Printfln("Getting running task")
				return err
			}

			if len(lto.TaskArns) == 0 {
				return fmt.Errorf("running task not found")
			}

			pterm.Success.Printfln("Getting running task")

			out, err := ecsSvc.ExecuteCommand(&ecs.ExecuteCommandInput{
				Container:   &os.Args[2],
				Interactive: aws.Bool(true),
				Cluster:     &clusterName,
				Task:        lto.TaskArns[0],
				Command:     aws.String("/bin/sh"),
			})
			if err != nil {
				pterm.Error.Printfln("Executing command")
				return err
			}

			pterm.Success.Printfln("Executing command")

			ssmCmd := ssmsession.NewSSMPluginCommand(cc.config.AwsRegion)
			ssmCmd.Start((out.Session))

			if err != nil {
				return err
			}

			return nil
		},
	}

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}
