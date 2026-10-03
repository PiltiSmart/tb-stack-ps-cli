package doctor

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

type CheckResult struct {
	Name    string
	Status  string // PASS, WARN, FAIL
	Details string
}

func RunDiagnostics() bool {
	ui.PrintBanner("PiltiSmart Pre-flight Health & Dependency Diagnostics")

	var results []CheckResult
	allPass := true

	// 1. Host OS & Architecture
	results = append(results, CheckResult{
		Name:    "Host OS / Architecture",
		Status:  "PASS",
		Details: fmt.Sprintf("%s / %s (CPUs: %d)", runtime.GOOS, runtime.GOARCH, runtime.NumCPU()),
	})

	// 2. Docker Engine check
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		results = append(results, CheckResult{
			Name:    "Docker Engine Installed",
			Status:  "FAIL",
			Details: "docker binary not found in PATH",
		})
		allPass = false
	} else {
		// Test daemon reachability
		cmd := exec.Command("docker", "info", "--format", "{{.ServerVersion}}")
		out, err := cmd.Output()
		if err != nil {
			results = append(results, CheckResult{
				Name:    "Docker Engine Running",
				Status:  "FAIL",
				Details: fmt.Sprintf("Daemon not responding or no permission (%s)", dockerPath),
			})
			allPass = false
		} else {
			serverVer := strings.TrimSpace(string(out))
			results = append(results, CheckResult{
				Name:    "Docker Engine Daemon",
				Status:  "PASS",
				Details: fmt.Sprintf("Running v%s (%s)", serverVer, dockerPath),
			})
		}
	}

	// 3. Docker Compose check
	composeCmd := exec.Command("docker", "compose", "version", "--short")
	out, err := composeCmd.Output()
	if err == nil {
		results = append(results, CheckResult{
			Name:    "Docker Compose Plugin",
			Status:  "PASS",
			Details: fmt.Sprintf("v%s (docker compose)", strings.TrimSpace(string(out))),
		})
	} else {
		standalonePath, err2 := exec.LookPath("docker-compose")
		if err2 == nil {
			cmd2 := exec.Command(standalonePath, "version", "--short")
			out2, _ := cmd2.Output()
			results = append(results, CheckResult{
				Name:    "Docker Compose Standalone",
				Status:  "PASS",
				Details: fmt.Sprintf("v%s (%s)", strings.TrimSpace(string(out2)), standalonePath),
			})
		} else {
			results = append(results, CheckResult{
				Name:    "Docker Compose",
				Status:  "FAIL",
				Details: "docker compose plugin or docker-compose executable not found",
			})
			allPass = false
		}
	}

	// 4. Memory check
	memTotal, memAvail := getMemoryInfo()
	if memTotal > 0 {
		memStatus := "PASS"
		if memAvail < 1024*1024*1024 { // Less than 1GB available
			memStatus = "WARN"
		}
		results = append(results, CheckResult{
			Name:    "Host System Memory (RAM)",
			Status:  memStatus,
			Details: fmt.Sprintf("Total: %.1f GB | Available: %.1f GB", float64(memTotal)/(1024*1024*1024), float64(memAvail)/(1024*1024*1024)),
		})
	}

	// 5. Disk space check
	diskFree, diskTotal := getDiskSpace("/")
	if diskTotal > 0 {
		diskStatus := "PASS"
		if diskFree < 2*1024*1024*1024 { // Less than 2GB free
			diskStatus = "WARN"
		}
		results = append(results, CheckResult{
			Name:    "Root Storage Disk Space",
			Status:  diskStatus,
			Details: fmt.Sprintf("Free: %.1f GB / Total: %.1f GB", float64(diskFree)/(1024*1024*1024), float64(diskTotal)/(1024*1024*1024)),
		})
	}

	// 6. Network Port Availability checks for PiltiSmart software components
	portsToCheck := []struct {
		Port int
		Desc string
	}{
		{5432, "tb-db (TimescaleDB / PostgreSQL)"},
		{8080, "tb-app (ThingsBoard Core Web UI)"},
		{1883, "tb-app (ThingsBoard Core MQTT Broker)"},
		{7070, "tb-app (ThingsBoard Core RPC)"},
		{8082, "tb-edge (ThingsBoard Edge Web UI)"},
		{1884, "tb-edge (ThingsBoard Edge MQTT)"},
		{8085, "jenkins (Jenkins Web UI)"},
		{50000, "jenkins (Jenkins Agent Listener)"},
		{9000, "piltiservices (API Gateway) / minio (S3 API)"},
		{9001, "minio (MinIO Web Console)"},
		{9092, "kafka (Apache Kafka PLAINTEXT Broker)"},
		{8088, "pulseX (PulseX Cloud Gateway / PMX)"},
	}

	fmt.Printf("\n[Dependency & Component Checks]\n")
	for _, r := range results {
		badge := fmt.Sprintf("%s✔ PASS%s", ui.ColorGreen, ui.ColorReset)
		if r.Status == "FAIL" {
			badge = fmt.Sprintf("%s✖ FAIL%s", ui.ColorRed, ui.ColorReset)
		} else if r.Status == "WARN" {
			badge = fmt.Sprintf("%s▲ WARN%s", ui.ColorYellow, ui.ColorReset)
		}
		fmt.Printf("  [%s] %-28s : %s\n", badge, r.Name, r.Details)
	}

	fmt.Printf("\n[PiltiSmart Software Port Availability]\n")
	for _, p := range portsToCheck {
		isAvailable := checkPortAvailable(p.Port)
		if isAvailable {
			fmt.Printf("  [%s✔ FREE%s] Port %-5d : Available (%s)\n", ui.ColorGreen, ui.ColorReset, p.Port, p.Desc)
		} else {
			fmt.Printf("  [%s▲ USED%s] Port %-5d : In Use/Bound (%s)\n", ui.ColorYellow, ui.ColorReset, p.Port, p.Desc)
		}
	}

	fmt.Println()
	if allPass {
		ui.Success("Core dependencies verified. System is ready to deploy PiltiSmart software!")
	} else {
		ui.Error("Some core dependencies are missing or require attention. See above for details.")
	}
	return allPass
}

func checkPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func getMemoryInfo() (uint64, uint64) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0
	}
	var total, avail uint64
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			if fields[0] == "MemTotal:" {
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				total = val * 1024
			} else if fields[0] == "MemAvailable:" {
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				avail = val * 1024
			}
		}
	}
	return total, avail
}
