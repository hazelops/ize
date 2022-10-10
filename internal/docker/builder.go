package docker

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/oklog/ulid"
)

type Builder struct {
	BuildArgs  map[string]*string
	Tags       []string
	Dockerfile string
	CacheFrom  []string
	Platform   string
}

func NewBuilder(buildArgs map[string]*string, tags []string, dockerfile string, cacheFrom []string, platform string) Builder {
	return Builder{
		BuildArgs:  buildArgs,
		Tags:       tags,
		Dockerfile: dockerfile,
		CacheFrom:  cacheFrom,
		Platform:   platform,
	}
}

func (b *Builder) Build(w io.Writer, contextDir string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("unable to create Docker client: %s", err)
	}

	dockerfile := b.Dockerfile
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	if !filepath.IsAbs(dockerfile) {
		dockerfile = filepath.Join(contextDir, dockerfile)
	}

	// If the dockerfile is outside of our build context, then we copy it
	// into our build context.
	relDockerfile, err := filepath.Rel(contextDir, dockerfile)
	if err != nil || strings.HasPrefix(relDockerfile, "..") {
		id, err := ulid.New(ulid.Now(), rand.Reader)
		if err != nil {
			return err
		}

		newPath := filepath.Join(contextDir, fmt.Sprintf("Dockerfile-%s", id.String()))
		if err := copyFile(dockerfile, newPath); err != nil {
			return err
		}
		defer os.Remove(newPath)

		dockerfile = newPath
	}

	contextDir, relDockerfile, err = build.GetContextFromLocalDir(contextDir, dockerfile)
	if err != nil {
		return fmt.Errorf("unable to create Docker context: %s", err)
	}

	if err := b.buildWithDocker(w, cli, contextDir, relDockerfile, b.Tags, b.BuildArgs); err != nil {
		return err
	}

	return nil
}

func (b *Builder) buildWithDocker(
	w io.Writer,
	cli *client.Client,
	contextDir string,
	relDockerfile string,
	tags []string,
	buildArgs map[string]*string,
) error {
	excludes, err := build.ReadDockerignore(contextDir)
	if err != nil {
		return fmt.Errorf("unable to read .dockerignore: %s", err)
	}

	if err := build.ValidateContextDirectory(contextDir, excludes); err != nil {
		return fmt.Errorf("error checking context: %s", err)
	}

	// And canonicalize dockerfile name to a platform-independent one
	relDockerfile = archive.CanonicalTarNameForPath(relDockerfile)

	excludes = build.TrimBuildFilesFromExcludes(excludes, relDockerfile, false)
	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return fmt.Errorf("unable to compress context: %s", err)
	}

	buildOpts := types.ImageBuildOptions{
		Dockerfile: relDockerfile,
		Tags:       tags,
		BuildArgs:  buildArgs,
		Platform:   b.Platform,
	}

	resp, err := cli.ImageBuild(context.Background(), buildCtx, buildOpts)
	if err != nil {
		return fmt.Errorf("error building image: %s", err)
	}
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, w, os.Stdout.Fd(), true, nil)
	if err != nil {
		return fmt.Errorf("unable to stream build logs to the terminal: %s", err)
	}

	return nil
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}
