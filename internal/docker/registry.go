package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/pkg/terminal"
)

type Registry struct {
	Registry string
	Token    string
}

func NewRegistry(registry, token string) Registry {
	return Registry{
		Registry: registry,
		Token:    token,
	}
}

func (r *Registry) Push(ctx context.Context, ui terminal.UI, image string, tags []string) error {
	if len(tags) == 0 {
		tags = []string{"latest"}
	}

	sg := ui.StepGroup()
	defer sg.Wait()
	s := sg.Add("%s: pushing image...", image)
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("unable to create Docker client: %s", err)
	}

	if len(tags) > 1 {
		for i := 1; i < len(tags); i++ {
			err = cli.ImageTag(context.Background(), image+":"+tags[0], fmt.Sprintf("%s/%s:%s", r.Registry, image, tags[i]))
			if err != nil {
				return err
			}
		}
	}

	resp, err := cli.ImagePush(ctx, image+":"+tags[0], types.ImagePushOptions{
		RegistryAuth: r.Token,
		All:          true,
	})
	if err != nil {
		return fmt.Errorf("unable to push image: %s", err)
	}

	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return err
	}

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(
		resp,
		s.TermOutput(),
		termFd,
		true,
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to stream push logs to the terminal: %s", err)
	}

	s.Done()

	return nil
}
