package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute runs the command based on program arguments.
func Execute() error { return rootCmd().Execute() }

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "khub",
		Short: "Manage the khub modular monolith",
	}
	version := os.Getenv("khub_VERSION")
	if version == "" {
		version = "0.0.1"
	}
	cmd.AddCommand(
		versionCmd(version),
		serverCmd(version),
		dataSinkCmd(version),
		mySQLReplTopoCmd(version),
	)
	return cmd
}

func versionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "display the khub client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(rootCmd().Use + " v" + version)
		},
	}
	return cmd
}
