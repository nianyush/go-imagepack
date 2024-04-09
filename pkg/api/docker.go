package api

import (
	"bytes"
	"io"

	_ "embed"

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

	buildContext, err := getBuildContext(ctx.FS, path...)
	if err != nil {
		return err
	}

	buildOpts := types.ImageBuildOptions{
		Context:    buildContext,
		Dockerfile: "Dockerfile",
		Tags:       []string{image},
	}

	// build docker image
	resp, err := cli.ImageBuild(ctx.Ctx, buildContext, buildOpts)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func getBuildContext(fs vfs.FS, path ...string) (io.Reader, error) {
	ta := util.NewTar(fs).AddBytes("Dockerfile", dockerfile)
	for _, p := range path {
		ta.AddPath(p)
	}

	b, err := ta.Build()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
