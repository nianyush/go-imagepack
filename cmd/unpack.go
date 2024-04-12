package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kairos-io/kairos-sdk/bundles"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

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

func init() {
	unpackcmd.Flags().BoolP("local", "l", false, "using local image")
	rootcmd.AddCommand(unpackcmd)
}
