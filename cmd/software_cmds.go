package cmd

import (
	"fmt"
	"os"

	"github.com/piltismart/tb-stack-ps-cli/pkg/software"
	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	softwareBaseDir string
	softwarePort    int
	softwareVersion string
	softwareAutoYes bool
	softwareDBUser  string
	softwareDBPass  string
	logFollow       bool
)

func createSoftwareCommand(s software.Software) *cobra.Command {
	swCmd := &cobra.Command{
		Use:     s.ID,
		Aliases: s.Aliases,
		Short:   fmt.Sprintf("Manage %s (%s)", s.Name, s.Category),
		Long: fmt.Sprintf(`Manage %s.
Category: %s
Default Ports: %v
Description: %s`, s.Name, s.Category, s.DefaultPorts, s.Description),
		Run: func(cmd *cobra.Command, args []string) {
			_ = software.Status(&s)
			fmt.Printf("\nAvailable subcommands for '%s':\n", s.ID)
			fmt.Printf("  pilti %s install  - Deploy and run container\n", s.ID)
			fmt.Printf("  pilti %s status   - Check live runtime container status\n", s.ID)
			fmt.Printf("  pilti %s logs     - View container logs (-f to follow)\n", s.ID)
			fmt.Printf("  pilti %s restart  - Restart container\n", s.ID)
			fmt.Printf("  pilti %s stop     - Stop container\n", s.ID)
			fmt.Printf("  pilti %s remove   - Remove container\n", s.ID)
		},
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: fmt.Sprintf("Install and run %s", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
			opts := software.InstallOptions{
				BaseDir:    softwareBaseDir,
				Version:    softwareVersion,
				Port:       softwarePort,
				AutoYes:    softwareAutoYes,
				DBUser:     softwareDBUser,
				DBPassword: softwareDBPass,
			}
			if err := software.InstallWithOptions(&s, opts); err != nil {
				ui.Error("Failed to install %s: %v", s.ID, err)
				os.Exit(1)
			}
		},
	}
	installCmd.Flags().StringVarP(&softwareBaseDir, "dir", "d", software.DefaultBaseDir, "Base installation directory")
	installCmd.Flags().IntVarP(&softwarePort, "port", "p", 0, "Custom host port number")
	installCmd.Flags().StringVarP(&softwareVersion, "version", "v", "", "Custom software version tag")
	installCmd.Flags().BoolVarP(&softwareAutoYes, "yes", "y", false, "Automatically accept defaults without interactive prompts")
	installCmd.Flags().StringVar(&softwareDBUser, "db-user", "", "PostgreSQL database username (for tb-db)")
	installCmd.Flags().StringVar(&softwareDBPass, "db-password", "", "PostgreSQL database password (for tb-db)")

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: fmt.Sprintf("Check status of %s", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
			_ = software.Status(&s)
		},
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf("Start %s container", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
			if err := software.Start(&s, softwareBaseDir); err != nil {
				ui.Error("Failed to start %s: %v", s.ID, err)
				os.Exit(1)
			}
		},
	}
	startCmd.Flags().StringVarP(&softwareBaseDir, "dir", "d", software.DefaultBaseDir, "Base installation directory")

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: fmt.Sprintf("Stop %s container", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
			if err := software.Stop(&s, softwareBaseDir); err != nil {
				ui.Error("Failed to stop %s: %v", s.ID, err)
				os.Exit(1)
			}
		},
	}
	stopCmd.Flags().StringVarP(&softwareBaseDir, "dir", "d", software.DefaultBaseDir, "Base installation directory")

	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: fmt.Sprintf("Restart %s container", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
			if err := software.Restart(&s, softwareBaseDir); err != nil {
				ui.Error("Failed to restart %s: %v", s.ID, err)
				os.Exit(1)
			}
		},
	}
	restartCmd.Flags().StringVarP(&softwareBaseDir, "dir", "d", software.DefaultBaseDir, "Base installation directory")

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: fmt.Sprintf("View logs of %s", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
			if err := software.Logs(&s, logFollow); err != nil {
				ui.Error("Failed to retrieve logs for %s: %v", s.ID, err)
			}
		},
	}
	logsCmd.Flags().BoolVarP(&logFollow, "follow", "f", false, "Follow log output")

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: fmt.Sprintf("Remove container and compose service for %s", s.Name),
		Run: func(cmd *cobra.Command, args []string) {
			if err := software.Remove(&s, softwareBaseDir); err != nil {
				ui.Error("Failed to remove %s: %v", s.ID, err)
				os.Exit(1)
			}
			ui.Success("Removed %s successfully.", s.Name)
		},
	}
	removeCmd.Flags().StringVarP(&softwareBaseDir, "dir", "d", software.DefaultBaseDir, "Base installation directory")

	swCmd.AddCommand(installCmd)
	swCmd.AddCommand(statusCmd)
	swCmd.AddCommand(startCmd)
	swCmd.AddCommand(stopCmd)
	swCmd.AddCommand(restartCmd)
	swCmd.AddCommand(logsCmd)
	swCmd.AddCommand(removeCmd)

	return swCmd
}

// RegisterSoftwareCommands registers individual subcommands for all software in Registry into the root command.
func RegisterSoftwareCommands(root *cobra.Command) {
	for _, s := range software.Registry {
		root.AddCommand(createSoftwareCommand(s))
	}
}
