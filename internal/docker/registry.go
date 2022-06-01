package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
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

func (r *Registry) Push(ctx context.Context, w io.Writer, image string, tags []string) error {
	if len(tags) == 0 {
		tags = []string{"latest"}
	}

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

	var termFd uintptr
	if f, ok := w.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(
		resp,
		w,
		termFd,
		true,
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to stream push logs to the terminal: %s", err)
	}

	return nil
}
