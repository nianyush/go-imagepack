package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/nianyush/go-imagepack/pkg/api"
	"github.com/nianyush/go-imagepack/pkg/util"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-vfs/v4"
)

var	packcmd = &cobra.Command{
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

			pathFS := vfs.NewPathFS(vfs.OSFS, absPath)

			ctx := util.Contenxt{
				Ctx: context.Background(),
				FS:  pathFS,
			}

			if err := api.DockerBuild(ctx, image, "/"); err != nil {
				panic(err)
			}
		},
	}

func init() {
	rootcmd.AddCommand(packcmd)
}