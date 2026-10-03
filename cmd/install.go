package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/software"
	"github.com/piltismart/tb-stack-ps-cli/pkg/stack"
	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	deployDir      string
	repoURL        string
	edgeWebPort    int
	edgeMqttPort   int
	installPort    int
	installVersion string
	installAutoYes bool
)

var installCmd = &cobra.Command{
	Use:   "install [software-id]",
	Short: "Install any PiltiSmart software component or launch interactive wizard",
	Long: `Install any PiltiSmart software component.
If no argument is given, launches the interactive terminal software selector.

Available software IDs:
  - tb-app        (ThingsBoard Core)
  - tb-db         (TimescaleDB / Postgres)
  - tb-edge       (ThingsBoard Edge)
  - jenkins       (Jenkins CI/CD Automation)
  - piltiservices (PiltiSmart Microservices)
  - kafka         (Apache Kafka Broker)
  - pulseX        (PulseX Cloud Gateway / PMX - formerly PiltiCloud)
  - pilticloud    (Alias for pulseX)
  - all           (Install all components)
  - tb-stack      (Legacy 3-component ThingsBoard stack)

Examples:
  pilti install tb-app
  pilti install jenkins
  pilti install kafka
  pilti install piltiservices
  pilti install pulseX
  pilti install pilticloud
  pilti install             # Interactive selection wizard`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var target string

		if len(args) == 0 {
			// Interactive mode
			selected, err := software.PromptSelectSoftware()
			if err != nil {
				ui.Warning("%v", err)
				return
			}
			target = selected.ID
		} else {
			target = strings.TrimSpace(args[0])
		}

		if target == "tb-stack" {
			err := stack.InstallTBStack(deployDir, repoURL, edgeWebPort, edgeMqttPort)
			if err != nil {
				ui.Error("Installation failed: %v", err)
				os.Exit(1)
			}
			return
		}

		if target == "all" {
			ui.PrintBanner("Installing All PiltiSmart Software Packages")
			for _, s := range software.Registry {
				if err := software.Install(&s, deployDir); err != nil {
					ui.Error("Failed to install %s: %v", s.ID, err)
				}
			}
			return
		}

		sw, found := software.GetSoftware(target)
		if !found {
			ui.Error("Unknown software '%s'.", target)
			fmt.Println("\nAvailable software packages (run 'pilti list' to see all):")
			for _, s := range software.Registry {
				fmt.Printf("  - %-14s : %s\n", s.ID, s.Name)
			}
			os.Exit(1)
		}

		opts := software.InstallOptions{
			BaseDir: deployDir,
			Version: installVersion,
			Port:    installPort,
			AutoYes: installAutoYes,
		}

		err := software.InstallWithOptions(sw, opts)
		if err != nil {
			ui.Error("Installation failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	installCmd.Flags().StringVarP(&deployDir, "dir", "d", software.DefaultBaseDir, "Base installation directory")
	installCmd.Flags().StringVarP(&repoURL, "repo", "r", stack.DefaultCatalogRepo, "Catalog raw repository URL (for remote templates)")
	installCmd.Flags().IntVar(&edgeWebPort, "edge-web-port", 0, "Host port for ThingsBoard Edge Web UI (default: 8082)")
	installCmd.Flags().IntVar(&edgeMqttPort, "edge-mqtt-port", 0, "Host port for ThingsBoard Edge MQTT broker (default: 1884)")
	installCmd.Flags().IntVarP(&installPort, "port", "p", 0, "Custom host port number")
	installCmd.Flags().StringVarP(&installVersion, "version", "v", "", "Custom software version tag (e.g. for pulseX)")
	installCmd.Flags().BoolVarP(&installAutoYes, "yes", "y", false, "Automatically accept defaults without interactive prompts")
}
