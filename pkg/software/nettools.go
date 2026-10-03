package software

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

type NetToolInfo struct {
	Name        string
	Binary      string
	Description string
}

var CoreNetTools = []NetToolInfo{
	{Name: "iproute2 (ip)", Binary: "ip", Description: "IP routing, address & network device management"},
	{Name: "iproute2 (ss)", Binary: "ss", Description: "Modern socket statistics & network port inspection"},
	{Name: "net-tools (netstat)", Binary: "netstat", Description: "Network connections & routing tables"},
	{Name: "net-tools (ifconfig)", Binary: "ifconfig", Description: "Network interface configuration"},
	{Name: "iputils (ping)", Binary: "ping", Description: "ICMP network echo & latency test"},
	{Name: "dnsutils (dig)", Binary: "dig", Description: "DNS lookup & domain name resolution"},
	{Name: "dnsutils (nslookup)", Binary: "nslookup", Description: "DNS query & name server verification"},
	{Name: "curl", Binary: "curl", Description: "HTTP/REST API client & data transfer"},
	{Name: "wget", Binary: "wget", Description: "Non-interactive network file downloader"},
	{Name: "traceroute", Binary: "traceroute", Description: "Network packet route & hop tracing"},
	{Name: "mtr", Binary: "mtr", Description: "Combined traceroute & continuous ping latency monitor"},
	{Name: "netcat (nc)", Binary: "nc", Description: "Arbitrary TCP/UDP connections & port listeners"},
	{Name: "socat", Binary: "socat", Description: "Multipurpose bidirectional relay & socket proxy"},
	{Name: "tcpdump", Binary: "tcpdump", Description: "Command-line packet analyzer & traffic capture"},
	{Name: "openssh (ssh)", Binary: "ssh", Description: "Secure remote shell & cryptographic tunneling"},
	{Name: "openssh (sshd)", Binary: "sshd", Description: "OpenSSH daemon / remote access server"},
	{Name: "openssl", Binary: "openssl", Description: "TLS/SSL certificates & cryptographic toolkit"},
}

type DistroInfo struct {
	Name           string
	PackageManager string
	UpdateCmd      []string
	InstallCmd     []string
	Packages       []string
}

// DetectDistroAndPackageManager inspects the host system to determine the OS, distro, and package manager.
func DetectDistroAndPackageManager() DistroInfo {
	if runtime.GOOS == "windows" {
		return DistroInfo{
			Name:           "Windows",
			PackageManager: "native / winget",
			Packages:       []string{"curl", "openssh", "ping", "tracert", "nslookup", "netstat"},
		}
	}

	if runtime.GOOS == "darwin" {
		return DistroInfo{
			Name:           "macOS (Darwin)",
			PackageManager: "brew",
			UpdateCmd:      []string{"brew", "update"},
			InstallCmd:     []string{"brew", "install"},
			Packages: []string{
				"curl", "wget", "traceroute", "mtr", "netcat", "socat", "tcpdump", "openssl",
			},
		}
	}

	// Linux: Read /etc/os-release
	distroName := "Linux"
	id := ""
	idLike := ""
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				distroName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
			} else if strings.HasPrefix(line, "ID=") {
				id = strings.ToLower(strings.Trim(strings.TrimPrefix(line, "ID="), `"`))
			} else if strings.HasPrefix(line, "ID_LIKE=") {
				idLike = strings.ToLower(strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), `"`))
			}
		}
	}

	sudo := ""
	if os.Geteuid() != 0 {
		if _, err := exec.LookPath("sudo"); err == nil {
			sudo = "sudo"
		}
	}

	// 1. Debian / Ubuntu family
	if id == "ubuntu" || id == "debian" || id == "pop" || id == "kali" || id == "linuxmint" || id == "raspbian" ||
		strings.Contains(idLike, "debian") || strings.Contains(idLike, "ubuntu") || hasCmd("apt-get") {
		update := []string{"apt-get", "update", "-y"}
		install := []string{"apt-get", "install", "-y"}
		if sudo != "" {
			update = append([]string{sudo}, update...)
			install = append([]string{sudo}, install...)
		}
		return DistroInfo{
			Name:           distroName,
			PackageManager: "apt-get",
			UpdateCmd:      update,
			InstallCmd:     install,
			Packages: []string{
				"iproute2",
				"net-tools",
				"iputils-ping",
				"dnsutils",
				"curl",
				"wget",
				"traceroute",
				"mtr-tiny",
				"netcat-openbsd",
				"socat",
				"tcpdump",
				"openssh-server",
				"openssl",
			},
		}
	}

	// 2. RHEL / CentOS / Rocky / AlmaLinux / Fedora family
	if id == "fedora" || id == "rhel" || id == "centos" || id == "rocky" || id == "almalinux" || id == "ol" || id == "amzn" ||
		strings.Contains(idLike, "rhel") || strings.Contains(idLike, "fedora") || hasCmd("dnf") || hasCmd("yum") {
		mgr := "dnf"
		if !hasCmd("dnf") && hasCmd("yum") {
			mgr = "yum"
		}
		update := []string{mgr, "makecache", "-y"}
		install := []string{mgr, "install", "-y"}
		if sudo != "" {
			update = append([]string{sudo}, update...)
			install = append([]string{sudo}, install...)
		}
		return DistroInfo{
			Name:           distroName,
			PackageManager: mgr,
			UpdateCmd:      update,
			InstallCmd:     install,
			Packages: []string{
				"iproute",
				"net-tools",
				"iputils",
				"bind-utils",
				"curl",
				"wget",
				"traceroute",
				"mtr",
				"nc",
				"socat",
				"tcpdump",
				"openssh-server",
				"openssl",
			},
		}
	}

	// 3. Alpine Linux
	if id == "alpine" || hasCmd("apk") {
		install := []string{"apk", "add", "--no-cache"}
		if sudo != "" {
			install = append([]string{sudo}, install...)
		}
		return DistroInfo{
			Name:           distroName,
			PackageManager: "apk",
			UpdateCmd:      nil,
			InstallCmd:     install,
			Packages: []string{
				"iproute2",
				"net-tools",
				"iputils",
				"bind-tools",
				"curl",
				"wget",
				"traceroute",
				"mtr",
				"netcat-openbsd",
				"socat",
				"tcpdump",
				"openssh",
				"openssl",
			},
		}
	}

	// 4. Arch Linux / Manjaro
	if id == "arch" || id == "manjaro" || hasCmd("pacman") {
		install := []string{"pacman", "-S", "--noconfirm", "--needed"}
		if sudo != "" {
			install = append([]string{sudo}, install...)
		}
		return DistroInfo{
			Name:           distroName,
			PackageManager: "pacman",
			UpdateCmd:      nil,
			InstallCmd:     install,
			Packages: []string{
				"iproute2",
				"net-tools",
				"iputils",
				"bind",
				"curl",
				"wget",
				"traceroute",
				"mtr",
				"gnu-netcat",
				"socat",
				"tcpdump",
				"openssh",
				"openssl",
			},
		}
	}

	// 5. openSUSE / SLES
	if id == "opensuse" || id == "sles" || hasCmd("zypper") {
		install := []string{"zypper", "install", "-y"}
		if sudo != "" {
			install = append([]string{sudo}, install...)
		}
		return DistroInfo{
			Name:           distroName,
			PackageManager: "zypper",
			UpdateCmd:      nil,
			InstallCmd:     install,
			Packages: []string{
				"iproute2",
				"net-tools",
				"iputils",
				"bind-utils",
				"curl",
				"wget",
				"traceroute",
				"mtr",
				"netcat-openbsd",
				"socat",
				"tcpdump",
				"openssh",
				"openssl",
			},
		}
	}

	// Generic Linux fallback
	return DistroInfo{
		Name:           distroName,
		PackageManager: "unknown",
		Packages: []string{
			"iproute2", "net-tools", "iputils-ping", "dnsutils", "curl", "wget",
			"traceroute", "mtr-tiny", "netcat-openbsd", "socat", "tcpdump", "openssh-server", "openssl",
		},
	}
}

func hasCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// FindToolBinary searches system paths for a binary.
func FindToolBinary(name string) string {
	if path, err := exec.LookPath(name); err == nil {
		return path
	}
	if runtime.GOOS == "windows" {
		if path, err := exec.LookPath(name + ".exe"); err == nil {
			return path
		}
	}
	// Common Unix system sbin/bin paths
	standardPaths := []string{
		"/usr/sbin/" + name,
		"/sbin/" + name,
		"/usr/bin/" + name,
		"/bin/" + name,
		"/usr/local/bin/" + name,
		"/usr/local/sbin/" + name,
	}
	for _, p := range standardPaths {
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p
		}
	}
	return ""
}

// CheckNetToolsStatus returns the current readiness of the network tools suite.
func CheckNetToolsStatus() string {
	presentCount := 0
	for _, t := range CoreNetTools {
		if FindToolBinary(t.Binary) != "" {
			presentCount++
		}
	}

	if presentCount >= 14 {
		return "INSTALLED (READY)"
	} else if presentCount >= 4 {
		return fmt.Sprintf("PARTIAL (%d/%d)", presentCount, len(CoreNetTools))
	}
	return "NOT INSTALLED"
}

// InstallNetTools installs the entire network diagnostic suite across supported distros.
func InstallNetTools(s *Software, opts InstallOptions) error {
	ui.PrintBanner("PiltiSmart Network & Diagnostic Suite Installation")

	distro := DetectDistroAndPackageManager()
	ui.Info("Target Environment : %s%s%s", ui.ColorBold, distro.Name, ui.ColorReset)
	ui.Info("Package Manager    : %s%s%s", ui.ColorCyan, distro.PackageManager, ui.ColorReset)

	// 1. Version / Profile confirmation
	reader := bufio.NewReader(os.Stdin)
	chosenVersion := opts.Version
	if chosenVersion == "" {
		chosenVersion = "v1.0.0 (Standard Full Suite)"
	}
	if !opts.AutoYes {
		fmt.Printf("\n  • Software Profile / Version : %s%s%s\n", ui.ColorGreen, chosenVersion, ui.ColorReset)
		fmt.Printf("  • Install complete diagnostic toolchain on %s? [Y/n]: ", distro.Name)
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans != "" && strings.ToLower(ans) != "y" && strings.ToLower(ans) != "yes" {
			ui.Info("Installation cancelled by user.")
			return nil
		}
	}

	if runtime.GOOS == "windows" {
		ui.Info("Windows environment detected. Verifying native PowerShell & Windows networking tools...")
		return NetToolsStatus(s)
	}

	if distro.PackageManager == "unknown" || len(distro.InstallCmd) == 0 {
		return fmt.Errorf("unsupported Linux distribution or no recognized package manager found (apt-get, dnf, yum, apk, pacman, zypper)")
	}

	// 2. Run repository update if defined
	if len(distro.UpdateCmd) > 0 {
		ui.Info("Updating package repository index (%s)...", strings.Join(distro.UpdateCmd, " "))
		cmd := exec.Command(distro.UpdateCmd[0], distro.UpdateCmd[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	// 3. Run installation command
	fullInstallCmd := append(distro.InstallCmd, distro.Packages...)
	ui.Info("Executing package installer: %s%s%s", ui.ColorCyan, strings.Join(fullInstallCmd, " "), ui.ColorReset)

	cmd := exec.Command(fullInstallCmd[0], fullInstallCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		ui.Warning("Batch package installation encountered a warning or non-zero exit code: %v", err)
		ui.Info("Retrying individual package installations to ensure maximum coverage...")
		for _, pkg := range distro.Packages {
			singleCmd := append(distro.InstallCmd, pkg)
			sc := exec.Command(singleCmd[0], singleCmd[1:]...)
			_ = sc.Run()
		}
	}

	// 4. Start / Enable OpenSSH service if sshd binary was installed
	if FindToolBinary("sshd") != "" {
		ui.Info("Configuring OpenSSH Server daemon...")
		sudo := ""
		if os.Geteuid() != 0 && hasCmd("sudo") {
			sudo = "sudo"
		}
		if hasCmd("systemctl") {
			if sudo != "" {
				_ = exec.Command(sudo, "systemctl", "enable", "--now", "ssh").Run()
				_ = exec.Command(sudo, "systemctl", "enable", "--now", "sshd").Run()
			} else {
				_ = exec.Command("systemctl", "enable", "--now", "ssh").Run()
				_ = exec.Command("systemctl", "enable", "--now", "sshd").Run()
			}
		} else if hasCmd("service") {
			if sudo != "" {
				_ = exec.Command(sudo, "service", "ssh", "start").Run()
			} else {
				_ = exec.Command("service", "ssh", "start").Run()
			}
		}
	}

	fmt.Println()
	ui.Success("PiltiSmart Network & Diagnostic Suite installed successfully!")
	return NetToolsStatus(s)
}

// NetToolsStatus prints the detailed checklist of all installed network tools.
func NetToolsStatus(s *Software) error {
	distro := DetectDistroAndPackageManager()

	ui.PrintBanner("PiltiSmart Network & Diagnostic Suite Status")
	fmt.Printf("  - Suite Identifier : %s%s%s\n", ui.ColorCyan, s.ID, ui.ColorReset)
	fmt.Printf("  - Host Distro      : %s\n", distro.Name)
	fmt.Printf("  - Overall Status   : %s\n\n", CheckNetToolsStatus())

	fmt.Printf("%-24s %-12s %-30s %s\n", "TOOL / PACKAGE", "BINARY", "SYSTEM PATH", "STATUS")
	fmt.Println(strings.Repeat("-", 80))

	installedCount := 0
	for _, t := range CoreNetTools {
		path := FindToolBinary(t.Binary)
		statusStr := fmt.Sprintf("%s[✔ INSTALLED]%s", ui.ColorGreen, ui.ColorReset)
		displayPath := path
		if path == "" {
			statusStr = fmt.Sprintf("%s[✖ MISSING]%s", ui.ColorRed, ui.ColorReset)
			displayPath = "-"
		} else {
			installedCount++
		}
		fmt.Printf("%-24s %-12s %-30s %s\n", t.Name, t.Binary, displayPath, statusStr)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total Available: %d / %d tools ready on host system.\n\n", installedCount, len(CoreNetTools))

	fmt.Println("Quick Diagnostic Commands:")
	fmt.Println("  • Check IP addresses      : ip a  (or: ifconfig)")
	fmt.Println("  • Check listening ports   : ss -tulpn  (or: netstat -tulpn)")
	fmt.Println("  • Test connectivity/ping  : ping -c 4 8.8.8.8")
	fmt.Println("  • Query DNS resolution    : dig google.com  (or: nslookup google.com)")
	fmt.Println("  • Trace network hops      : traceroute 8.8.8.8  (or: mtr 8.8.8.8)")
	fmt.Println("  • Test TCP port listener  : nc -zv <host> <port>")
	fmt.Println("  • Capture network traffic : tcpdump -i any -n\n")

	return nil
}

// RestartNetTools restarts sshd service if present.
func RestartNetTools(s *Software) error {
	ui.Info("Restarting OpenSSH service...")
	sudo := ""
	if os.Geteuid() != 0 && hasCmd("sudo") {
		sudo = "sudo"
	}
	if hasCmd("systemctl") {
		args := []string{"systemctl", "restart", "ssh"}
		if sudo != "" {
			args = append([]string{sudo}, args...)
		}
		_ = exec.Command(args[0], args[1:]...).Run()
		args = []string{"systemctl", "restart", "sshd"}
		if sudo != "" {
			args = append([]string{sudo}, args...)
		}
		_ = exec.Command(args[0], args[1:]...).Run()
	}
	ui.Success("Network tools and SSH service refreshed.")
	return nil
}

// StopNetTools stops ssh service if running.
func StopNetTools(s *Software) error {
	ui.Info("Stopping OpenSSH service...")
	sudo := ""
	if os.Geteuid() != 0 && hasCmd("sudo") {
		sudo = "sudo"
	}
	if hasCmd("systemctl") {
		args := []string{"systemctl", "stop", "ssh"}
		if sudo != "" {
			args = append([]string{sudo}, args...)
		}
		_ = exec.Command(args[0], args[1:]...).Run()
	}
	ui.Success("OpenSSH service stopped.")
	return nil
}

// RemoveNetTools displays instructions or executes package purge.
func RemoveNetTools(s *Software) error {
	distro := DetectDistroAndPackageManager()
	ui.PrintBanner("PiltiSmart Network Tools Removal")
	ui.Warning("Network tools (curl, ip, openssh, etc.) are core system utilities.")
	fmt.Printf("To remove these packages on %s, run:\n", distro.Name)
	if distro.PackageManager == "apt-get" {
		fmt.Printf("  sudo apt-get remove --purge -y %s\n\n", strings.Join(distro.Packages, " "))
	} else if distro.PackageManager == "dnf" || distro.PackageManager == "yum" {
		fmt.Printf("  sudo %s remove -y %s\n\n", distro.PackageManager, strings.Join(distro.Packages, " "))
	} else {
		fmt.Printf("  Use your package manager (%s) to remove packages.\n\n", distro.PackageManager)
	}
	return nil
}
