package stack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

func Status(deployDir string) error {
	ui.PrintBanner("ThingsBoard Stack Container Status")

	cmd := exec.Command("docker", "ps", "--format", "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Down(deployDir string) error {
	ui.PrintBanner("Tearing Down ThingsBoard 3-Component Stack")

	if deployDir == "" {
		deployDir = "/opt/piltismart/tb-stack"
	}

	components := []string{"edge-tb", "tb", "tb-db"}
	for _, c := range components {
		cDir := filepath.Join(deployDir, c)
		if _, err := os.Stat(cDir); err == nil {
			ui.Info("Stopping and removing component: %s...", c)
			cmd := exec.Command("docker", "compose", "down")
			cmd.Dir = cDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
		}
	}

	ui.Success("Stack teardown complete!")
	return nil
}

func Restart(deployDir string, component string) error {
	ui.PrintBanner(fmt.Sprintf("Restarting Stack: %s", component))

	if deployDir == "" {
		deployDir = "/opt/piltismart/tb-stack"
	}

	var targets []string
	if component == "all" || component == "" {
		targets = []string{"tb-db", "tb", "edge-tb"}
	} else {
		targets = []string{component}
	}

	for _, c := range targets {
		cDir := filepath.Join(deployDir, c)
		ui.Info("Restarting component: %s...", c)
		cmd := exec.Command("docker", "compose", "restart")
		cmd.Dir = cDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	ui.Success("Restart complete.")
	return nil
}
