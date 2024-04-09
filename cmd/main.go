package main

import (
	"context"
	"os"

	"log/slog"

	"github.com/nianyush/go-imagepack/pkg/api"
	"github.com/nianyush/go-imagepack/pkg/util"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-vfs/v4"
)

func main() {
	cmd := &cobra.Command{
		Use:  "pack <image> <path>",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			image := args[0]
			path := args[1]

			absPath, err := util.AbsPath(path)
			if err != nil {
				slog.Error("failed to get absolute path", err)
				os.Exit(1)
			}

			ctx := util.Contenxt{
				Ctx: context.Background(),
				FS:  vfs.OSFS,
			}

			if err := api.DockerBuild(ctx, image, absPath); err != nil {
				panic(err)
			}
		},
	}

	if err := cmd.Execute(); err != nil {
		slog.Error("failed to execute command", err)
		os.Exit(1)
	}
}
