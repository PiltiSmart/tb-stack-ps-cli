package cmd

import (
	"os"

	"github.com/piltismart/tb-stack-ps-cli/pkg/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run pre-flight health checks and dependency diagnostics",
	Long:  "Inspect host OS, Docker engine, Docker compose plugin, memory, disk, and network ports availability for ThingsBoard.",
	Run: func(cmd *cobra.Command, args []string) {
		success := doctor.RunDiagnostics()
		if !success {
			os.Exit(1)
		}
	},
}
