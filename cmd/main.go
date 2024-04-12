package main

import (
	"os"

	"log/slog"

	"github.com/spf13/cobra"
)

var rootcmd = &cobra.Command{
	Use: "imagectl",
}

func main() {
	if err := rootcmd.Execute(); err != nil {
		slog.Error("failed to execute command", err)
		os.Exit(1)
	}
}
