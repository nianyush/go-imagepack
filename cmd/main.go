package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/kairos-io/kairos-sdk/bundles"
	"github.com/nianyush/go-imagepack/pkg/api"
	"github.com/nianyush/go-imagepack/pkg/util"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-vfs/v4"
)

var rootcmd = &cobra.Command{
	Use: "imagectl",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			opts.Level = slog.LevelDebug
		}

		s := slog.New(slog.NewTextHandler(os.Stdout, opts))
		slog.SetDefault(s)
	},
}

var packcmd = &cobra.Command{
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

		slog.Info(fmt.Sprintf("packing image %s from %s", image, absPath))

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

var unpackcmd = &cobra.Command{
	Use:  "unpack <image> <path>",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		path := args[1]
		local, _ := cmd.Flags().GetBool("local")

		if !strings.HasPrefix(image, "docker://") && !strings.HasPrefix(image, "container://") {
			image = fmt.Sprintf("container://%s", image)
		}

		slog.Info(fmt.Sprintf("unpacking image %s to %s, local: %v", image, path, local))

		opts := []bundles.BundleOption{
			bundles.WithTarget(image),
			bundles.WithRootFS(path),
			bundles.WithLocalFile(local),
		}

		if err := bundles.RunBundles(opts); err != nil {
			slog.Error("failed to unpack image", err)
			os.Exit(1)
		}
	},
}

func main() {
	rootcmd.Flags().BoolP("debug", "d", false, "enable debug mode")

	unpackcmd.Flags().BoolP("local", "l", false, "using local image")

	rootcmd.AddCommand(packcmd)
	rootcmd.AddCommand(unpackcmd)
	if err := rootcmd.Execute(); err != nil {
		slog.Error("failed to execute command", err)
		os.Exit(1)
	}
}
