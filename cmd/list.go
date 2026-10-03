package cmd

import (
	"fmt"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/software"
	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "softwares", "apps"},
	Short:   "List all available PiltiSmart software components and their live status",
	Long: `List all supported software applications in PiltiSmart including:
  - tb-app        (ThingsBoard Core Application)
  - tb-db         (TimescaleDB / PostgreSQL database)
  - tb-edge       (ThingsBoard Edge Gateway)
  - jenkins       (Jenkins CI/CD Engine)
  - piltiservices (PiltiSmart Microservices)
  - kafka         (Apache Kafka Broker)
  - pulseX        (PulseX Cloud Gateway / PMX - formerly PiltiCloud)`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("PiltiSmart Software Catalog")

		fmt.Println("Available PiltiSmart software components:")
		fmt.Println(strings.Repeat("-", 95))
		fmt.Printf("%-14s %-28s %-16s %-18s %s\n", "SOFTWARE ID", "SOFTWARE NAME", "CATEGORY", "PORTS", "STATUS")
		fmt.Println(strings.Repeat("-", 95))

		for _, s := range software.Registry {
			status := software.CheckStatus(&s)
			statusFormatted := fmt.Sprintf("%s%s%s", ui.ColorGreen, status, ui.ColorReset)
			if strings.Contains(status, "STOPPED") {
				statusFormatted = fmt.Sprintf("%s%s%s", ui.ColorYellow, status, ui.ColorReset)
			} else if strings.Contains(status, "NOT") {
				statusFormatted = fmt.Sprintf("%s%s%s", ui.ColorRed, status, ui.ColorReset)
			}

			portsStr := "-"
			if len(s.DefaultPorts) > 0 {
				portsStr = strings.Split(s.DefaultPorts[0], " ")[0]
				if len(s.DefaultPorts) > 1 {
					portsStr += ", " + strings.Split(s.DefaultPorts[1], " ")[0]
				}
			}

			fmt.Printf("%-14s %-28s %-16s %-18s %s\n",
				s.ID,
				s.Name,
				s.Category,
				portsStr,
				statusFormatted,
			)
		}
		fmt.Println(strings.Repeat("-", 95))

		fmt.Println("\nTo manage any software, use subcommands or run interactively:")
		fmt.Printf("  • %spilti <software-id> install%s  (e.g., %spilti tb-app install%s)\n", ui.ColorBold, ui.ColorReset, ui.ColorCyan, ui.ColorReset)
		fmt.Printf("  • %spilti <software-id> status%s   (e.g., %spilti jenkins status%s)\n", ui.ColorBold, ui.ColorReset, ui.ColorCyan, ui.ColorReset)
		fmt.Printf("  • %spilti <software-id> restart%s  (e.g., %spilti tb-db restart%s)\n", ui.ColorBold, ui.ColorReset, ui.ColorCyan, ui.ColorReset)
		fmt.Printf("  • %spilti install%s                (Launches interactive software picker)\n\n", ui.ColorBold, ui.ColorReset)
	},
}
