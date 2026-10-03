package software

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

const DefaultBaseDir = "/opt/piltismart"

type InstallOptions struct {
	BaseDir    string
	Version    string
	Port       int
	AutoYes    bool
	DBUser     string
	DBPassword string
}

// Install provisions configuration and starts containers for a given software using interactive prompts.
func Install(s *Software, baseDir string) error {
	return InstallWithOptions(s, InstallOptions{BaseDir: baseDir})
}

// InstallWithOptions provisions configuration with full interactive or parameterized options.
func InstallWithOptions(s *Software, opts InstallOptions) error {
	baseDir := opts.BaseDir
	if baseDir == "" {
		baseDir = DefaultBaseDir
	}

	targetDir := filepath.Join(baseDir, s.Subdir)
	ui.PrintBanner(fmt.Sprintf("Installing PiltiSmart Software: %s", s.Name))
	ui.Info("Identifier : %s%s%s", ui.ColorCyan, s.ID, ui.ColorReset)
	ui.Info("Directory  : %s", targetDir)
	ui.Info("Default Ports : %s", strings.Join(s.DefaultPorts, ", "))

	// Check dependencies if any
	for _, depID := range s.Dependencies {
		if dep, exists := GetSoftware(depID); exists {
			depStatus := CheckStatus(dep)
			if depStatus != "RUNNING" {
				fmt.Println()
				ui.PrintBanner("Dependency Requirement Error")
				ui.Error("Cannot install '%s' (%s): Missing required prerequisite!", s.Name, s.ID)
				ui.Error("Prerequisite dependency '%s' (%s) is NOT running (Current status: %s).", dep.Name, dep.ID, depStatus)
				fmt.Println()
				ui.Info("👉 Step 1: You must install and start '%s' first:", dep.ID)
				ui.Info("   pilti %s install   (or: pilti install %s)", dep.ID, dep.ID)
				fmt.Println()
				ui.Info("👉 Step 2: Once '%s' is RUNNING, install '%s':", dep.ID, s.ID)
				ui.Info("   pilti %s install   (or: pilti install %s)", s.ID, s.ID)
				fmt.Println()
				return fmt.Errorf("dependency '%s' is not running (status: %s); please install '%s' first", dep.ID, depStatus, dep.ID)
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	normID := strings.ToLower(s.ID)

	// 1. Version Selection for all software
	chosenVersion := opts.Version
	if chosenVersion != "" {
		ui.Info("Using specified %s version: %s", s.Name, chosenVersion)
	} else if opts.AutoYes {
		chosenVersion = s.Version
		if chosenVersion == "" {
			chosenVersion = "latest"
		}
		ui.Info("Using default recommended %s version: %s", s.Name, chosenVersion)
	} else {
		var vErr error
		chosenVersion, vErr = PromptSoftwareVersion(s, reader, opts.AutoYes)
		if vErr != nil {
			return fmt.Errorf("version selection failed: %w", vErr)
		}
	}

	// 2. Port Configuration & Live Availability Checking
	var configuredPorts map[string]int
	if opts.Port > 0 {
		configuredPorts = make(map[string]int)
		if len(s.PortConfigs) > 0 {
			configuredPorts[s.PortConfigs[0].Name] = opts.Port
			ui.Info("Using specified port %d for %s", opts.Port, s.PortConfigs[0].Name)
			if !CheckPortAvailable(opts.Port) {
				ui.Warning("Port %d is occupied by another process on target system!", opts.Port)
			} else {
				ui.Success("Port %d is available on target system!", opts.Port)
			}
		}
	} else {
		var pErr error
		configuredPorts, pErr = PromptAndCheckPorts(s, reader, opts.AutoYes)
		if pErr != nil {
			return fmt.Errorf("port check failed: %w", pErr)
		}
	}

	// 2.5 PostgreSQL Database Credentials (for tb-db)
	var dbUser, dbPass string
	if normID == "tb-db" || normID == "db" {
		dbUser = opts.DBUser
		dbPass = opts.DBPassword
		if opts.AutoYes {
			if dbUser == "" {
				dbUser = "postgres"
			}
			if dbPass == "" {
				dbPass = "postgres"
			}
		} else if dbUser == "" || dbPass == "" {
			var err error
			dbUser, dbPass, err = PromptPostgresCredentials(reader, dbUser, dbPass)
			if err != nil {
				return fmt.Errorf("database credentials prompt failed: %w", err)
			}
		}
	}

	// 3. Write configured templates
	ui.Info("Writing configuration and compose manifests...")
	if err := WriteConfiguredTemplatesWithCreds(s, targetDir, configuredPorts, chosenVersion, dbUser, dbPass); err != nil {
		return fmt.Errorf("failed to write templates for %s: %w", s.ID, err)
	}

	// Detect if .env / env_file is mentioned in docker compose
	envFiles := DetectEnvFiles(targetDir, s)
	if len(envFiles) > 0 {
		DisplayEnvNoticeWhileInstalling(envFiles, targetDir)
	}

	// 4. Run docker compose up -d
	ui.Info("Pulling and launching Docker container(s) for %s...", s.ID)
	if err := runDockerCompose(targetDir, "up", "-d"); err != nil {
		return fmt.Errorf("failed to start container for %s: %w", s.ID, err)
	}

	time.Sleep(2 * time.Second)
	currStatus := CheckStatus(s)
	ui.Success("Software '%s' deployed successfully! Container Status: %s", s.Name, currStatus)

	if len(envFiles) > 0 {
		DisplayEnvNoticeInstallationDone(envFiles, targetDir, s.ID)
	}

	PrintSoftwareSummaryDetails(s, targetDir, configuredPorts, chosenVersion)
	return nil
}

// Start brings up the software container.
func Start(s *Software, baseDir string) error {
	if baseDir == "" {
		baseDir = DefaultBaseDir
	}
	targetDir := filepath.Join(baseDir, s.Subdir)

	ui.Info("Starting %s (%s)...", s.Name, s.ID)
	if err := runDockerCompose(targetDir, "start"); err != nil {
		// Fallback to up -d
		return runDockerCompose(targetDir, "up", "-d")
	}
	ui.Success("%s is now running.", s.Name)
	return nil
}

// Stop halts the software container.
func Stop(s *Software, baseDir string) error {
	if baseDir == "" {
		baseDir = DefaultBaseDir
	}
	targetDir := filepath.Join(baseDir, s.Subdir)

	ui.Info("Stopping %s (%s)...", s.Name, s.ID)
	if err := runDockerCompose(targetDir, "stop"); err != nil {
		return err
	}
	ui.Success("%s has been stopped.", s.Name)
	return nil
}

// Restart restarts the software container.
func Restart(s *Software, baseDir string) error {
	if baseDir == "" {
		baseDir = DefaultBaseDir
	}
	targetDir := filepath.Join(baseDir, s.Subdir)

	ui.Info("Restarting %s (%s)...", s.Name, s.ID)
	if err := runDockerCompose(targetDir, "restart"); err != nil {
		return err
	}
	ui.Success("%s restarted successfully.", s.Name)
	return nil
}

// Status prints the detailed runtime container status for a software.
func Status(s *Software) error {
	ui.PrintBanner(fmt.Sprintf("Software Status: %s", s.Name))
	fmt.Printf("  - Software ID   : %s%s%s\n", ui.ColorCyan, s.ID, ui.ColorReset)
	fmt.Printf("  - Container     : %s\n", s.ContainerName)
	fmt.Printf("  - Category      : %s\n", s.Category)
	fmt.Printf("  - State         : ")

	st := CheckStatus(s)
	if strings.Contains(st, "RUNNING") {
		fmt.Printf("%s%s%s\n", ui.ColorGreen, st, ui.ColorReset)
	} else if strings.Contains(st, "STOPPED") {
		fmt.Printf("%s%s%s\n", ui.ColorYellow, st, ui.ColorReset)
	} else {
		fmt.Printf("%s%s%s\n", ui.ColorRed, st, ui.ColorReset)
	}

	fmt.Printf("  - Ports         : %s\n", strings.Join(s.DefaultPorts, ", "))
	fmt.Printf("  - Description   : %s\n\n", s.Description)

	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", s.ContainerName), "--format", "table {{.Names}}\t{{.Status}}\t{{.Ports}}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	return nil
}

// Logs streams or prints logs for the container.
func Logs(s *Software, follow bool) error {
	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, "--tail", "100", s.ContainerName)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Remove tears down the container and associated resources.
func Remove(s *Software, baseDir string) error {
	if baseDir == "" {
		baseDir = DefaultBaseDir
	}
	targetDir := filepath.Join(baseDir, s.Subdir)

	ui.Info("Removing container and resources for %s...", s.Name)
	return runDockerCompose(targetDir, "down")
}

func runDockerCompose(dir string, args ...string) error {
	var cmd *exec.Cmd

	// Check if docker compose or docker-compose is available
	testCmd := exec.Command("docker", "compose", "version")
	if err := testCmd.Run(); err == nil {
		allArgs := append([]string{"compose"}, args...)
		cmd = exec.Command("docker", allArgs...)
	} else {
		cmd = exec.Command("docker-compose", args...)
	}

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}

func PrintSoftwareSummary(s *Software, targetDir string) {
	PrintSoftwareSummaryDetails(s, targetDir, nil, "")
}

func PrintSoftwareSummaryDetails(s *Software, targetDir string, configuredPorts map[string]int, version string) {
	fmt.Println()
	fmt.Println("==================================================================")
	ui.Success("Software '%s' ready!", s.Name)
	fmt.Println("==================================================================")
	fmt.Printf("  %-25s : %s%s%s\n", "Software ID", ui.ColorCyan, s.ID, ui.ColorReset)
	fmt.Printf("  %-25s : %s\n", "Container Name", s.ContainerName)
	fmt.Printf("  %-25s : %s\n", "Config Directory", targetDir)
	if version != "" {
		fmt.Printf("  %-25s : %s%s%s\n", "Software Version", ui.ColorBold, version, ui.ColorReset)
	}

	if len(configuredPorts) > 0 {
		for name, port := range configuredPorts {
			fmt.Printf("  %-25s : %sport %d%s\n", name, ui.ColorGreen, port, ui.ColorReset)
		}
	} else {
		for i, p := range s.DefaultPorts {
			label := "Default Endpoint"
			if i > 0 {
				label = fmt.Sprintf("Endpoint (%d)", i+1)
			}
			fmt.Printf("  %-25s : %s%s%s\n", label, ui.ColorGreen, p, ui.ColorReset)
		}
	}

	if s.ID == "tb-db" {
		dbUser := "postgres"
		composePath := filepath.Join(targetDir, "docker-compose.yml")
		if data, err := os.ReadFile(composePath); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "POSTGRES_USER:") {
					parts := strings.SplitN(trimmed, ":", 2)
					if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
						dbUser = strings.TrimSpace(parts[1])
					}
				}
			}
		}
		fmt.Printf("  %-25s : %s%s%s\n", "Database User", ui.ColorBold, dbUser, ui.ColorReset)
		fmt.Printf("  %-25s : %s\n", "Database Name", "thingsboard")
	}

	envFiles := DetectEnvFiles(targetDir, s)
	if len(envFiles) > 0 {
		fmt.Println("------------------------------------------------------------------")
		for _, env := range envFiles {
			fullPath := filepath.Join(targetDir, env)
			fmt.Printf("  %-25s : %s%s%s\n", "Environment Config", ui.ColorYellow, ".env manually paste #installation is done plz update .env (infisical file)", ui.ColorReset)
			fmt.Printf("  %-25s : %s%s%s\n", "Env File Path", ui.ColorBold, fullPath, ui.ColorReset)
		}
	}

	fmt.Println("==================================================================")
	fmt.Printf("Quick commands:\n")
	fmt.Printf("  pilti %s status  -> Check status\n", s.ID)
	fmt.Printf("  pilti %s logs    -> View logs\n", s.ID)
	fmt.Printf("  pilti %s restart -> Restart service\n", s.ID)
	fmt.Printf("  pilti %s stop    -> Stop service\n", s.ID)
	fmt.Println("==================================================================")
}

// DetectEnvFiles inspects docker-compose.yml, software definitions, and the target directory
// to detect any environment files referenced in the compose configuration.
func DetectEnvFiles(targetDir string, s *Software) []string {
	var envFiles []string
	seen := make(map[string]bool)

	addFile := func(name string) {
		name = strings.TrimSpace(name)
		if name == "" {
			return
		}
		base := filepath.Base(name)
		if !seen[base] {
			seen[base] = true
			envFiles = append(envFiles, base)
		}
	}

	// 1. Explicitly configured software env file
	if s != nil && s.EnvFile != "" {
		addFile(s.EnvFile)
	}

	// 2. Parse docker-compose.yml if present
	composePath := filepath.Join(targetDir, "docker-compose.yml")
	if data, err := os.ReadFile(composePath); err == nil {
		content := string(data)
		lines := strings.Split(content, "\n")
		inEnvFileBlock := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "env_file:") {
				inEnvFileBlock = true
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					val := strings.TrimSpace(parts[1])
					val = strings.Trim(val, `"'[]`)
					if val != "" && !strings.HasPrefix(val, "#") {
						addFile(val)
					}
				}
				continue
			}
			if inEnvFileBlock {
				if strings.HasPrefix(trimmed, "- ") {
					val := strings.TrimPrefix(trimmed, "- ")
					val = strings.TrimSpace(val)
					val = strings.Trim(val, `"'`)
					if val != "" {
						addFile(val)
					}
				} else if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
					inEnvFileBlock = false
				}
			}
		}

		// Also check if .env is mentioned anywhere in the compose content
		if len(envFiles) == 0 && (strings.Contains(content, "env_file") || strings.Contains(content, ".env")) {
			addFile(".env")
		}
	}

	// 3. Scan directory on disk for any .env* files
	if entries, err := os.ReadDir(targetDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				if strings.HasPrefix(name, ".env") || strings.HasSuffix(name, ".env") || strings.Contains(name, ".env.") {
					addFile(name)
				}
			}
		}
	}

	return envFiles
}

// DisplayEnvNoticeWhileInstalling displays notice while installation is in progress.
func DisplayEnvNoticeWhileInstalling(envFiles []string, targetDir string) {
	fmt.Println()
	fmt.Printf("%s%s▲ NOTICE: Environment configuration file detected in docker compose%s\n", ui.ColorBold, ui.ColorYellow, ui.ColorReset)
	for _, env := range envFiles {
		fullPath := filepath.Join(targetDir, env)
		fmt.Printf("  %s%s.env manually paste%s %s(path: %s)%s\n", ui.ColorBold, ui.ColorCyan, ui.ColorReset, ui.ColorDim, fullPath, ui.ColorReset)
	}
	fmt.Println()
}

// DisplayEnvNoticeInstallationDone displays notice when installation is complete.
func DisplayEnvNoticeInstallationDone(envFiles []string, targetDir string, softwareID string) {
	fmt.Println()
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("%s%s🔔 [ENVIRONMENT CONFIGURATION REQUIRED]%s\n", ui.ColorBold, ui.ColorYellow, ui.ColorReset)
	for _, env := range envFiles {
		fullPath := filepath.Join(targetDir, env)
		fmt.Printf("  %s%s.env manually paste #installation is done plz update .env (infisical file)%s\n", ui.ColorBold, ui.ColorYellow, ui.ColorReset)
		fmt.Printf("  %s--> Target File :%s %s%s%s\n", ui.ColorCyan, ui.ColorReset, ui.ColorBold, fullPath, ui.ColorReset)
	}
	fmt.Printf("  %sAfter updating credentials, restart service: %spilti %s restart%s\n", ui.ColorDim, ui.ColorCyan, softwareID, ui.ColorReset)
	fmt.Println("------------------------------------------------------------------")
}

func IsPortOccupied(port int) bool {
	return !CheckPortAvailable(port)
}

// PromptPostgresCredentials interactively asks for PostgreSQL username and password.
func PromptPostgresCredentials(reader *bufio.Reader, defaultUser, defaultPass string) (string, string, error) {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	if defaultUser == "" {
		defaultUser = "postgres"
	}
	if defaultPass == "" {
		defaultPass = "postgres"
	}

	fmt.Println()
	fmt.Println("==================================================================")
	fmt.Printf("%s%s[PostgreSQL / TimescaleDB Database Credentials]%s\n", ui.ColorBold, ui.ColorCyan, ui.ColorReset)
	fmt.Println("==================================================================")

	fmt.Printf("  Enter PostgreSQL Username [default: %s]: ", defaultUser)
	inputUser, err := reader.ReadString('\n')
	if err != nil {
		return defaultUser, defaultPass, err
	}
	inputUser = strings.TrimSpace(inputUser)
	if inputUser == "" {
		inputUser = defaultUser
	}

	fmt.Printf("  Enter PostgreSQL Password [default: %s]: ", defaultPass)
	inputPass, err := reader.ReadString('\n')
	if err != nil {
		return inputUser, defaultPass, err
	}
	inputPass = strings.TrimSpace(inputPass)
	if inputPass == "" {
		inputPass = defaultPass
	}

	fmt.Println("==================================================================")
	ui.Success("Database user configured: %s%s%s", ui.ColorBold, inputUser, ui.ColorReset)
	fmt.Println()

	return inputUser, inputPass, nil
}
