package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sullivtr/k8s_platform/internal/config"
	"github.com/sullivtr/k8s_platform/internal/providers"
)

func mySQLReplTopoCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capture-replication-topology",
		Short: "capture mysql replication topology from sources",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting capture of mysql replication topology from sources")
			c := config.Load(version, "")

			prvds := &providers.ModuleProviders{
				Config: c,
			}

			prvds.InitMySQLTopoProvider()
			prvds.InitStorageProvider()
			prvds.InitCacheProvider()

			catalog, err := prvds.StorageProvider.GetMySQLCatalog()
			if err != nil {
				fmt.Println("Error fetching mysql catalog")
			}

			prvds.MySQLTopoProvider.Session.SDK.Databases = catalog

			nodes, edges, err := prvds.MySQLTopoProvider.CaptureReplicationTopology()
			if err != nil {
				fmt.Printf("Error capturing mysql replication topology: %s\n", err.Error())
			}

			if err := prvds.CacheProvider.Put("mysql-repl-topo-nodes", nodes); err != nil {
				fmt.Printf("Error caching mysql replication topology nodes: %s\n", err.Error())
			}
			if err := prvds.CacheProvider.Put("mysql-repl-topo-edges", edges); err != nil {
				fmt.Printf("Error caching mysql replication topology edges: %s\n", err.Error())
			}

			fmt.Println("Capture of mysql replication topology from sources complete")
		},
	}

	return cmd
}
