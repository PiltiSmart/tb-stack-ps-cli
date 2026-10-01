package s3

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

const (
	DefaultAlias     = "myminio"
	DefaultEndpoint  = "http://localhost:9000"
	DefaultAccessKey = "your_access_key"
	DefaultSecretKey = "your_secret_key"
	MCReleaseVersion = "RELEASE.2025-07-16T15-35-03Z"
)

// Config represents MinIO/S3 connection parameters.
type Config struct {
	Alias     string
	Endpoint  string
	AccessKey string
	SecretKey string
}

// GetDefaultConfig returns the active or environment-overridden configuration.
func GetDefaultConfig() *Config {
	cfg := &Config{
		Alias:     DefaultAlias,
		Endpoint:  DefaultEndpoint,
		AccessKey: DefaultAccessKey,
		SecretKey: DefaultSecretKey,
	}

	if envEp := os.Getenv("MINIO_ENDPOINT"); envEp != "" {
		cfg.Endpoint = envEp
	}
	if envAK := os.Getenv("MINIO_ACCESS_KEY"); envAK != "" {
		cfg.AccessKey = envAK
	}
	if envSK := os.Getenv("MINIO_SECRET_KEY"); envSK != "" {
		cfg.SecretKey = envSK
	}
	if envAlias := os.Getenv("MINIO_ALIAS"); envAlias != "" {
		cfg.Alias = envAlias
	}

	return cfg
}

// GetMCURLs returns verified download URLs based on OS and Architecture.
func GetMCURLs() ([]string, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	var filename string
	switch osName {
	case "linux":
		switch arch {
		case "amd64":
			filename = fmt.Sprintf("mc.linux-amd64.%s", MCReleaseVersion)
		case "arm64":
			filename = fmt.Sprintf("mc.linux-arm64.%s", MCReleaseVersion)
		default:
			return nil, fmt.Errorf("unsupported linux architecture: %s", arch)
		}
	case "darwin":
		switch arch {
		case "arm64":
			filename = fmt.Sprintf("mc.darwin-arm64.%s", MCReleaseVersion)
		case "amd64":
			filename = fmt.Sprintf("mc.darwin-amd64.%s", MCReleaseVersion)
		default:
			return nil, fmt.Errorf("unsupported darwin architecture: %s", arch)
		}
	case "windows":
		if arch == "amd64" {
			filename = fmt.Sprintf("mc.windows-amd64.%s.exe", MCReleaseVersion)
		} else {
			return nil, fmt.Errorf("unsupported windows architecture: %s", arch)
		}
	default:
		return nil, fmt.Errorf("unsupported OS: %s", osName)
	}

	urls := []string{
		fmt.Sprintf("https://github.com/minio/mc/releases/download/%s/%s", MCReleaseVersion, filename),
	}
	return urls, nil
}

// FindMCExecutable checks PATH and standard installation paths for mc.
func FindMCExecutable() string {
	// 1. Check system PATH
	if path, err := exec.LookPath("mc"); err == nil {
		return path
	}
	if runtime.GOOS == "windows" {
		if path, err := exec.LookPath("mc.exe"); err == nil {
			return path
		}
	}

	// 2. Check standard Unix paths
	standardPaths := []string{
		"/usr/local/bin/mc",
		"/usr/bin/mc",
		filepath.Join(os.Getenv("HOME"), ".local/bin/mc"),
	}
	for _, p := range standardPaths {
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p
		}
	}

	// 3. Check Windows local app data
	if runtime.GOOS == "windows" {
		localApp := os.Getenv("LOCALAPPDATA")
		if localApp != "" {
			winPath := filepath.Join(localApp, "pilti", "bin", "mc.exe")
			if fi, err := os.Stat(winPath); err == nil && !fi.IsDir() {
				return winPath
			}
		}
	}

	return ""
}

// EnsureMC installs the mc binary if missing and configures the default alias.
func EnsureMC(cfg *Config) (string, error) {
	mcPath := FindMCExecutable()
	if mcPath == "" {
		ui.Info("MinIO Client ('mc') utility not found on host. Downloading official binary...")
		var err error
		mcPath, err = InstallMC()
		if err != nil {
			return "", fmt.Errorf("failed to auto-install MinIO client: %w", err)
		}
	}

	// Ensure alias is configured
	if err := EnsureAlias(mcPath, cfg); err != nil {
		return mcPath, fmt.Errorf("failed to configure MinIO alias: %w", err)
	}

	return mcPath, nil
}

// InstallMC downloads and installs the official MinIO client binary for the host OS/Arch.
func InstallMC() (string, error) {
	urls, err := GetMCURLs()
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 180 * time.Second}
	var lastErr error
	var resp *http.Response

	for _, downloadURL := range urls {
		ui.Info("Downloading 'mc' from: %s", downloadURL)
		req, rErr := http.NewRequest("GET", downloadURL, nil)
		if rErr != nil {
			lastErr = rErr
			continue
		}
		req.Header.Set("User-Agent", "curl/7.88.1")

		r, getErr := client.Do(req)
		if getErr != nil {
			lastErr = fmt.Errorf("network error downloading mc: %w", getErr)
			continue
		}

		if r.StatusCode != http.StatusOK {
			r.Body.Close()
			lastErr = fmt.Errorf("server returned HTTP %d while downloading mc from %s", r.StatusCode, downloadURL)
			continue
		}

		resp = r
		break
	}

	if resp == nil {
		return "", fmt.Errorf("failed to download mc binary: %v", lastErr)
	}
	defer resp.Body.Close()

	var targetPath string
	if runtime.GOOS == "windows" {
		targetDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "pilti", "bin")
		_ = os.MkdirAll(targetDir, 0755)
		targetPath = filepath.Join(targetDir, "mc.exe")
	} else {
		// Prefer /usr/local/bin if writable or running as root
		if os.Geteuid() == 0 || isDirWritable("/usr/local/bin") {
			targetPath = "/usr/local/bin/mc"
		} else {
			targetDir := filepath.Join(os.Getenv("HOME"), ".local/bin")
			_ = os.MkdirAll(targetDir, 0755)
			targetPath = filepath.Join(targetDir, "mc")
		}
	}

	// Write binary
	outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		// Try fallback to ~/.local/bin on Unix
		if runtime.GOOS != "windows" && targetPath != filepath.Join(os.Getenv("HOME"), ".local/bin/mc") {
			targetDir := filepath.Join(os.Getenv("HOME"), ".local/bin")
			_ = os.MkdirAll(targetDir, 0755)
			targetPath = filepath.Join(targetDir, "mc")
			outFile, err = os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		}
		if err != nil {
			return "", fmt.Errorf("failed to write mc binary to %s: %w", targetPath, err)
		}
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save mc binary: %w", err)
	}
	_ = outFile.Chmod(0755)

	ui.Success("MinIO Client ('mc') installed successfully to: %s", targetPath)
	return targetPath, nil
}

func isDirWritable(path string) bool {
	testFile := filepath.Join(path, ".pilti_test")
	if err := os.WriteFile(testFile, []byte(""), 0644); err != nil {
		return false
	}
	_ = os.Remove(testFile)
	return true
}

// EnsureAlias checks if the alias exists; if not, registers it with mc alias set.
func EnsureAlias(mcPath string, cfg *Config) error {
	// Test if alias is already configured
	cmd := exec.Command(mcPath, "alias", "list", cfg.Alias)
	if err := cmd.Run(); err == nil {
		// Already exists
		return nil
	}

	// Configure alias
	ui.Info("Configuring MinIO alias '%s' -> %s...", cfg.Alias, cfg.Endpoint)
	setCmd := exec.Command(mcPath, "alias", "set", cfg.Alias, cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	setCmd.Stdout = os.Stdout
	setCmd.Stderr = os.Stderr
	if err := setCmd.Run(); err != nil {
		return fmt.Errorf("failed to set mc alias '%s': %w", cfg.Alias, err)
	}

	ui.Success("MinIO alias '%s' configured successfully!", cfg.Alias)
	return nil
}

// NormalizePath converts s3://bucket/path or bare bucket/path into alias/bucket/path.
func NormalizePath(input string, defaultAlias string) string {
	val := strings.TrimSpace(input)
	if val == "" {
		return defaultAlias
	}

	// s3://mybucket/object -> myminio/mybucket/object
	if strings.HasPrefix(strings.ToLower(val), "s3://") {
		trimmed := strings.TrimPrefix(val, "s3://")
		trimmed = strings.TrimPrefix(trimmed, "S3://")
		return fmt.Sprintf("%s/%s", defaultAlias, trimmed)
	}

	// If it is a flag (e.g. -r, --force, --recursive), return as-is
	if strings.HasPrefix(val, "-") {
		return val
	}

	// If it refers to an existing local file or path starting with . or / or \, keep as local
	if strings.HasPrefix(val, ".") || strings.HasPrefix(val, "/") || strings.HasPrefix(val, "\\") || strings.Contains(val, ":\\") {
		return val
	}
	if fi, err := os.Stat(val); err == nil && (fi.Mode().IsRegular() || fi.IsDir()) {
		return val
	}

	// If already prefixed with alias, keep it
	if strings.HasPrefix(val, defaultAlias+"/") || val == defaultAlias {
		return val
	}

	// Default to alias prefix
	return fmt.Sprintf("%s/%s", defaultAlias, val)
}

// RunMC runs mc with the given arguments, connecting std streams.
func RunMC(mcPath string, args ...string) error {
	cmd := exec.Command(mcPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}
