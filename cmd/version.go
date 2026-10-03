package cmd

import (
	"fmt"
	"runtime"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show pilti CLI version and system runtime info",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("PiltiSmart Enterprise CLI (pilti)")
		fmt.Printf("  - CLI Binary    : %spilti%s\n", ui.ColorCyan, ui.ColorReset)
		fmt.Printf("  - CLI Version   : %s2.1.2 (Enterprise Software & Cloud Suite)%s\n", ui.ColorGreen, ui.ColorReset)
		fmt.Printf("  - Framework     : Cobra v1.8.1\n")
		fmt.Printf("  - Go Runtime    : %s\n", runtime.Version())
		fmt.Printf("  - OS / Arch     : %s / %s\n\n", runtime.GOOS, runtime.GOARCH)
	},
}
