package cmd

import (
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/software"
	"github.com/piltismart/tb-stack-ps-cli/pkg/stack"
	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	stackDir       string
	stackComponent string
)

var stackCmd = &cobra.Command{
	Use:   "stack [status|down|restart]",
	Short: "Legacy ThingsBoard 3-component stack lifecycle manager",
	Run: func(cmd *cobra.Command, args []string) {
		action := "status"
		if len(args) > 0 {
			action = args[0]
		}

		switch action {
		case "down":
			_ = stack.Down(stackDir)
		case "restart":
			_ = stack.Restart(stackDir, stackComponent)
		default:
			_ = stack.Status(stackDir)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status [software-id]",
	Short: "Check running container status for all or specific software",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			target := strings.TrimSpace(args[0])
			if sw, found := software.GetSoftware(target); found {
				_ = software.Status(sw)
				return
			}
			ui.Warning("Unknown software '%s'. Showing global status instead.", target)
		}

		ui.PrintBanner("PiltiSmart Global Container Status")
		_ = stack.Status(stackDir)
	},
}

func init() {
	stackCmd.Flags().StringVarP(&stackDir, "dir", "d", "/opt/piltismart/tb-stack", "Deployment directory")
	stackCmd.Flags().StringVarP(&stackComponent, "component", "c", "all", "Component to manage (all|tb-db|tb|edge-tb)")

	statusCmd.Flags().StringVarP(&stackDir, "dir", "d", "/opt/piltismart", "Deployment directory")
}
