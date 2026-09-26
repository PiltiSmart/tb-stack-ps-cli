package cmd

import (
	"fmt"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/stack"
	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var catalogRepoURL string

var catalogCmd = &cobra.Command{
	Use:   "catalog [list]",
	Short: "Discover stacks from the centralized GitHub catalog repository",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("PiltiSmart Centralized Stack Catalog")
		ui.Info("Connecting to GitHub: %s...", catalogRepoURL)

		cat, err := stack.FetchRemoteCatalog(catalogRepoURL)
		if err != nil {
			ui.Warning("Could not fetch remote catalog: %v", err)
			fmt.Println("\nLocally available stacks in binary:")
			fmt.Println(strings.Repeat("-", 75))
			fmt.Printf("%-12s %-32s %s\n", "STACK ID", "NAME", "COMPONENTS")
			fmt.Println(strings.Repeat("-", 75))
			fmt.Printf("%-12s %-32s %s\n", "tb-stack", "ThingsBoard Microservice Stack", "tb-db, tb, edge-tb")
			fmt.Println(strings.Repeat("-", 75))
			return
		}

		fmt.Printf("\nOrganization: %s | Catalog Version: %s | Updated: %s\n\n", cat.Organization, cat.Version, cat.UpdatedAt)
		fmt.Printf("%-12s %-10s %-35s %s\n", "STACK ID", "VERSION", "NAME", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 85))
		for _, s := range cat.Stacks {
			fmt.Printf("%-12s %-10s %-35s %s\n", s.ID, s.Version, s.Name, s.Description)
		}
		fmt.Println(strings.Repeat("-", 85))
		fmt.Printf("\nRun 'ps install <stack-id>' to deploy any stack.\n\n")
	},
}

func init() {
	catalogCmd.Flags().StringVarP(&catalogRepoURL, "repo", "r", stack.DefaultCatalogRepo, "Catalog raw repository URL")
}
