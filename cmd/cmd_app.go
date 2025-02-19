package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sullivtr/k8s_platform/internal/config"
	"github.com/sullivtr/k8s_platform/internal/server"
)

func serverCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-app",
		Short: "start the khub application server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting khub application server")
			c := config.Load(version, "")
			s := server.NewApp(c)

			s.Logger.Fatal(s.Start(fmt.Sprintf(":%d", c.ListenPort)))
		},
	}

	return cmd
}
