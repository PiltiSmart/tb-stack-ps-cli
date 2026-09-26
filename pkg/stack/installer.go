package stack

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

func InstallTBStack(deployDir, repoURL string, edgeWebPort, edgeMqttPort int) error {
	ui.PrintBanner("Orchestrating ThingsBoard 3-Component Stack")

	if deployDir == "" {
		deployDir = "/opt/piltismart/tb-stack"
	}

	// Auto-detect if default ports are in use
	if edgeWebPort == 0 {
		if isPortInUse(8080) {
			edgeWebPort = 8082
			ui.Info("Port 8080 is currently in use on host; auto-mapping ThingsBoard Edge Web UI to port 8082")
		} else {
			edgeWebPort = 8080
		}
	}

	if edgeMqttPort == 0 {
		if isPortInUse(1883) {
			edgeMqttPort = 1884
			ui.Info("Port 1883 is allocated to ThingsBoard Core; mapping ThingsBoard Edge MQTT to port 1884")
		} else {
			edgeMqttPort = 1884
		}
	}

	ui.Info("Target deployment directory: %s", deployDir)

	// Step 1: Ensure directory structure & timescale permissions
	if err := prepareDirectories(deployDir); err != nil {
		return fmt.Errorf("failed to prepare directories: %w", err)
	}

	// Step 2: Attempt dynamic pull from GitHub catalog repository
	err := DownloadStackManifests(repoURL, "tb-stack", deployDir)
	if err != nil {
		ui.Warning("Remote GitHub repository pull notice: %v", err)
		ui.Info("Falling back to embedded high-performance offline template manifests...")
		if err := writeEmbeddedStackFiles(deployDir); err != nil {
			return fmt.Errorf("failed to write fallback manifests: %w", err)
		}
		ui.Success("Embedded templates prepared successfully.")
	}

	// Step 3: Phase A - Deploy tb-db (TimescaleDB)
	ui.Info("[Phase 1/3] Deploying Database Stack: TimescaleDB / Postgres (tb-db)...")
	dbDir := filepath.Join(deployDir, "tb-db")
	if err := runDockerComposeUp(dbDir, nil); err != nil {
		return fmt.Errorf("failed to start tb-db: %w", err)
	}
	ui.Success("TimescaleDB container started! Waiting for database readiness on port 5432...")
	if err := waitForPort("127.0.0.1", 5432, 25*time.Second); err != nil {
		ui.Warning("Database port 5432 probe timed out, continuing...")
	} else {
		ui.Success("Database is listening on port 5432!")
	}

	// Step 4: Phase B - Deploy tb (ThingsBoard Core)
	ui.Info("[Phase 2/3] Deploying ThingsBoard Core Application Stack (tb)...")
	tbDir := filepath.Join(deployDir, "tb")
	if err := runDockerComposeUp(tbDir, nil); err != nil {
		return fmt.Errorf("failed to start tb: %w", err)
	}
	ui.Success("ThingsBoard Core container (Thingsboard-test) started!")

	// Step 5: Phase C - Deploy edge-tb (ThingsBoard Edge)
	ui.Info("[Phase 3/3] Deploying ThingsBoard Edge Stack (edge-tb)...")
	edgeDir := filepath.Join(deployDir, "edge-tb")
	envVars := map[string]string{
		"EDGE_WEB_PORT":  fmt.Sprintf("%d", edgeWebPort),
		"EDGE_MQTT_PORT": fmt.Sprintf("%d", edgeMqttPort),
	}
	if err := runDockerComposeUp(edgeDir, envVars); err != nil {
		return fmt.Errorf("failed to start edge-tb: %w", err)
	}
	ui.Success("ThingsBoard Edge container (mytbedge) started!")

	time.Sleep(3 * time.Second)
	PrintSummary(deployDir, edgeWebPort, edgeMqttPort)
	return nil
}

func prepareDirectories(deployDir string) error {
	dirs := []string{
		filepath.Join(deployDir, "tb-db"),
		filepath.Join(deployDir, "tb"),
		filepath.Join(deployDir, "edge-tb"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	tsDataDir := filepath.Join(deployDir, "tb-db", "timescale_data")
	if err := os.MkdirAll(tsDataDir, 0777); err != nil {
		return err
	}
	_ = os.Chmod(tsDataDir, 0777)
	return nil
}

func writeEmbeddedStackFiles(deployDir string) error {
	if err := os.WriteFile(filepath.Join(deployDir, "tb-db", "docker-compose.yml"), []byte(RawTbDbCompose), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(deployDir, "tb", "docker-compose.yml"), []byte(RawTbCompose), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(deployDir, "tb", ".tb.env"), []byte(RawTbEnv), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(deployDir, "edge-tb", "docker-compose.yml"), []byte(RawEdgeTbCompose), 0644); err != nil {
		return err
	}
	return nil
}

func isPortInUse(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

func runDockerComposeUp(dir string, extraEnv map[string]string) error {
	var cmd *exec.Cmd

	testCmd := exec.Command("docker", "compose", "version")
	if err := testCmd.Run(); err == nil {
		cmd = exec.Command("docker", "compose", "up", "-d")
	} else {
		cmd = exec.Command("docker-compose", "up", "-d")
	}

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	for k, v := range extraEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return cmd.Run()
}

func waitForPort(host string, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 1*time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout waiting for %s:%d", host, port)
}

func PrintSummary(deployDir string, edgeWebPort, edgeMqttPort int) {
	fmt.Println()
	fmt.Println("==================================================================")
	ui.Success("All 3 Compose Stacks successfully deployed and orchestrated!")
	fmt.Println("==================================================================")
	fmt.Printf("  %-30s : %shttp://<node-ip>:80%s\n", "ThingsBoard Core Web UI", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("  %-30s : %shttp://<node-ip>:%d%s\n", "ThingsBoard Edge Web UI", ui.ColorCyan, edgeWebPort, ui.ColorReset)
	fmt.Printf("  %-30s : %sport 1883%s\n", "ThingsBoard Core MQTT Broker", ui.ColorGreen, ui.ColorReset)
	fmt.Printf("  %-30s : %sport %d%s\n", "ThingsBoard Edge MQTT Port", ui.ColorGreen, edgeMqttPort, ui.ColorReset)
	fmt.Printf("  %-30s : %sport 5432 (user: postgres, db: thingsboard)%s\n", "TimescaleDB Storage Port", ui.ColorGreen, ui.ColorReset)
	fmt.Printf("  %-30s : %s%s%s\n", "Deployment Directory", ui.ColorBold, deployDir, ui.ColorReset)
	fmt.Println("==================================================================")
}
