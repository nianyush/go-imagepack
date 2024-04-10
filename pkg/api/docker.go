package api

import (
	"bytes"
	_ "embed"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nianyush/go-imagepack/pkg/util"
	"github.com/twpayne/go-vfs/v4"
)

//go:embed Dockerfile
var dockerfile []byte

func DockerBuild(ctx util.Contenxt, image string, path ...string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	buildContext, err := getBuildContext(ctx.FS, true, path...)
	if err != nil {
		return err
	}

	buildOpts := types.ImageBuildOptions{
		Context:    buildContext,
		Dockerfile: "Dockerfile",
		Remove:     true,
		Tags:       []string{image},
	}

	// build docker image
	resp, err := cli.ImageBuild(ctx.Ctx, buildContext, buildOpts)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err 
	}
	return nil
}

func getBuildContext(fs vfs.FS, includeDockerfile bool, path ...string) (io.Reader, error) {
	ta := util.NewTar(fs)
	if includeDockerfile {
		ta.AddBytes("Dockerfile", dockerfile)
	}
	for _, p := range path {
		ta.AddPath(p)
	}

	b, err := ta.Build()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
