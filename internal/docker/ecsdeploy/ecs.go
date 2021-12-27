package ecsdeploy

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Option struct {
	Tags       []string
	BuildArgs  map[string]*string
	Dockerfile string
	CacheFrom  []string
	ContextDir string
}

func Build(log logrus.Logger, opts Option) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Error("docker client initialization")
		return err
	}

	log.Debugf("build options: %s", opts)

	contextDir := opts.ContextDir
	dockerfile := opts.Dockerfile

	dockerfile, err = filepath.Rel(contextDir, dockerfile)
	if err != nil {
		return err
	}

	excludes, err := build.ReadDockerignore(contextDir)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to read .dockerignore: %s", err)
	}

	if err := build.ValidateContextDirectory(contextDir, excludes); err != nil {
		return status.Errorf(codes.Internal, "error checking context: %s", err)
	}

	excludes = build.TrimBuildFilesFromExcludes(excludes, opts.Dockerfile, false)

	fmt.Println(contextDir)

	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return status.Errorf(codes.Internal, "unable to compress context: %s", err)
	}

	resp, err := cli.ImageBuild(context.Background(), buildCtx, types.ImageBuildOptions{
		CacheFrom:  opts.CacheFrom,
		Tags:       opts.Tags,
		BuildArgs:  opts.BuildArgs,
		Dockerfile: dockerfile,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	wr := ioutil.Discard
	if log.GetLevel() >= 4 {
		wr = os.Stdout
	}

	var termFd uintptr

	err = jsonmessage.DisplayJSONMessagesStream(
		resp.Body,
		wr,
		termFd,
		true,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
