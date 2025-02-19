package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/sullivtr/k8s_platform/internal/config"
	"github.com/sullivtr/k8s_platform/internal/server"
)

func dataSinkCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-datasink",
		Short: "start the khub application data sink server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting khub data-sink server")
			c := config.Load(version, "")
			// Start the redis data sink server
			server.NewDataSink(c)
			waitOnSignal()
		},
	}

	return cmd
}

// waitOnSignal waits on signal
func waitOnSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info().Msg("Received signal to shutdown")
}
