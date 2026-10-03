package cmd

import (
	"os"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pilti",
	Short: "PiltiSmart Enterprise CLI - Software Catalog & Service Manager",
	Long: `==================================================================
  PiltiSmart Enterprise CLI Utility ('pilti')
  Unified management tool to inspect, deploy, and operate 
  PiltiSmart software components:
  - tb-app        (ThingsBoard Core Application)
  - tb-db         (TimescaleDB / PostgreSQL Database)
  - tb-edge       (ThingsBoard Edge Gateway)
  - jenkins       (Jenkins CI/CD Automation Engine)
  - piltiservices (PiltiSmart Microservices)
  - kafka         (Apache Kafka Broker)
  - pulseX        (PulseX Cloud Gateway / PMX - formerly PiltiCloud)
  - minio         (MinIO S3-Compatible Object Storage)
==================================================================`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("PiltiSmart Software Management Engine")
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	// Software catalog listing
	rootCmd.AddCommand(listCmd)

	// Global action commands
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(catalogCmd)
	rootCmd.AddCommand(stackCmd)
	rootCmd.AddCommand(s3Cmd)

	// Individual software commands (pilti tb-app, pilti jenkins, etc.)
	RegisterSoftwareCommands(rootCmd)
}
