package s3

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

const S3ConfigFileName = "s3_config.json"

// GetConfigFilePath returns the persistent config file path in ~/.pilti/s3_config.json.
func GetConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".pilti", S3ConfigFileName)
}

// LoadSavedConfig loads connection settings from ~/.pilti/s3_config.json.
// If missing, it checks for existing aliases in ~/.mc/config.json.
func LoadSavedConfig() (*Config, error) {
	p := GetConfigFilePath()
	if p != "" {
		data, err := os.ReadFile(p)
		if err == nil {
			var cfg Config
			if jErr := json.Unmarshal(data, &cfg); jErr == nil && cfg.Endpoint != "" {
				if cfg.Alias == "" {
					cfg.Alias = DefaultAlias
				}
				return &cfg, nil
			}
		}
	}

	// Fallback: check ~/.mc/config.json for existing alias
	if mcCfg := loadFromMCConfig(DefaultAlias); mcCfg != nil {
		_ = SaveConfig(mcCfg)
		return mcCfg, nil
	}

	return nil, fmt.Errorf("no saved S3 configuration found")
}

// loadFromMCConfig reads existing alias credentials from ~/.mc/config.json.
func loadFromMCConfig(alias string) *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	mcConfigFile := filepath.Join(home, ".mc", "config.json")
	data, err := os.ReadFile(mcConfigFile)
	if err != nil {
		return nil
	}

	type mcConfigStruct struct {
		Aliases map[string]struct {
			URL       string `json:"url"`
			AccessKey string `json:"accessKey"`
			SecretKey string `json:"secretKey"`
		} `json:"aliases"`
	}

	var m mcConfigStruct
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}

	if entry, ok := m.Aliases[alias]; ok && entry.URL != "" {
		return &Config{
			Alias:     alias,
			Endpoint:  entry.URL,
			AccessKey: entry.AccessKey,
			SecretKey: entry.SecretKey,
		}
	}

	return nil
}

// SaveConfig persists MinIO/S3 connection parameters to ~/.pilti/s3_config.json.
func SaveConfig(cfg *Config) error {
	p := GetConfigFilePath()
	if p == "" {
		return fmt.Errorf("unable to determine user home directory")
	}
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

// HasSavedConfig checks whether a valid persistent configuration exists.
func HasSavedConfig() bool {
	cfg, err := LoadSavedConfig()
	return err == nil && cfg != nil && cfg.Endpoint != ""
}

// GetActiveConfig returns the active configuration in order of priority:
// 1. Environment variables (MINIO_ENDPOINT, etc.)
// 2. Saved configuration in ~/.pilti/s3_config.json
// 3. Default configuration (http://localhost:9000)
func GetActiveConfig() *Config {
	if envEp := os.Getenv("MINIO_ENDPOINT"); envEp != "" {
		return GetDefaultConfig()
	}

	if saved, err := LoadSavedConfig(); err == nil && saved != nil && saved.Endpoint != "" {
		return saved
	}

	return GetDefaultConfig()
}

// FormatEndpoint constructs a valid HTTP/HTTPS endpoint from host and port.
func FormatEndpoint(host string, port int) string {
	scheme := "http"
	cleanHost := strings.TrimSpace(host)
	if strings.HasPrefix(strings.ToLower(cleanHost), "https://") {
		scheme = "https"
		cleanHost = strings.TrimPrefix(cleanHost, "https://")
		cleanHost = strings.TrimPrefix(cleanHost, "HTTPS://")
	} else if strings.HasPrefix(strings.ToLower(cleanHost), "http://") {
		scheme = "http"
		cleanHost = strings.TrimPrefix(cleanHost, "http://")
		cleanHost = strings.TrimPrefix(cleanHost, "HTTP://")
	}
	cleanHost = strings.TrimRight(cleanHost, "/")

	if strings.Contains(cleanHost, ":") {
		return fmt.Sprintf("%s://%s", scheme, cleanHost)
	}

	if port <= 0 {
		port = 9000
	}
	return fmt.Sprintf("%s://%s:%d", scheme, cleanHost, port)
}

func extractHostAndPort(endpoint string) (string, int) {
	clean := endpoint
	clean = strings.TrimPrefix(clean, "http://")
	clean = strings.TrimPrefix(clean, "https://")
	clean = strings.TrimPrefix(clean, "HTTP://")
	clean = strings.TrimPrefix(clean, "HTTPS://")
	clean = strings.TrimRight(clean, "/")

	if strings.Contains(clean, ":") {
		parts := strings.Split(clean, ":")
		p, err := strconv.Atoi(parts[1])
		if err == nil {
			return parts[0], p
		}
		return parts[0], 9000
	}
	return clean, 9000
}

// PromptServerConfig interactively prompts for server IP, port, username, and password.
func PromptServerConfig(reader *bufio.Reader, autoYes bool, defaultCfg *Config) *Config {
	if defaultCfg == nil {
		defaultCfg = GetActiveConfig()
	}

	cfg := &Config{
		Alias:     defaultCfg.Alias,
		Endpoint:  defaultCfg.Endpoint,
		AccessKey: defaultCfg.AccessKey,
		SecretKey: defaultCfg.SecretKey,
	}

	if autoYes {
		return cfg
	}

	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	currHost, currPort := extractHostAndPort(defaultCfg.Endpoint)
	if currHost == "" {
		currHost = "localhost"
	}

	fmt.Println()
	fmt.Println("==================================================================")
	fmt.Printf("%s%s[S3 / MinIO Server Connection Details]%s\n", ui.ColorBold, ui.ColorCyan, ui.ColorReset)
	fmt.Println("==================================================================")

	// 1. Server IP / Host
	fmt.Printf("  Enter MinIO Server IP or Hostname [default: %s]: ", currHost)
	inputHost, _ := reader.ReadString('\n')
	inputHost = strings.TrimSpace(inputHost)
	if inputHost == "" {
		inputHost = currHost
	}

	// 2. Port (Display default port 1st, then ask)
	defaultPort := currPort
	if defaultPort <= 0 {
		defaultPort = 9000
	}
	fmt.Printf("\n  Default Port: %s%d%s\n", ui.ColorGreen, defaultPort, ui.ColorReset)
	fmt.Printf("  Use default port %d? [Y/n] (or enter custom port): ", defaultPort)
	inputPortStr, _ := reader.ReadString('\n')
	inputPortStr = strings.TrimSpace(inputPortStr)

	finalPort := defaultPort
	if strings.ToLower(inputPortStr) == "n" || strings.ToLower(inputPortStr) == "no" {
		for {
			fmt.Printf("  Enter custom port number [1-65535]: ")
			pStr, _ := reader.ReadString('\n')
			pStr = strings.TrimSpace(pStr)
			if p, err := strconv.Atoi(pStr); err == nil && p > 0 && p <= 65535 {
				finalPort = p
				break
			}
			fmt.Printf("  %sInvalid port number. Please enter a value between 1 and 65535.%s\n", ui.ColorRed, ui.ColorReset)
		}
	} else if inputPortStr != "" && strings.ToLower(inputPortStr) != "y" && strings.ToLower(inputPortStr) != "yes" {
		if p, err := strconv.Atoi(inputPortStr); err == nil && p > 0 && p <= 65535 {
			finalPort = p
		} else {
			fmt.Printf("  %sInput '%s' not recognized as port; using default %d%s\n", ui.ColorYellow, inputPortStr, defaultPort, ui.ColorReset)
			finalPort = defaultPort
		}
	}

	cfg.Endpoint = FormatEndpoint(inputHost, finalPort)

	// 3. Username / Access Key
	defaultUser := defaultCfg.AccessKey
	if defaultUser == "" || defaultUser == DefaultAccessKey {
		defaultUser = "minioadmin"
	}
	fmt.Printf("\n  Enter Username / Access Key [default: %s]: ", defaultUser)
	inputUser, _ := reader.ReadString('\n')
	inputUser = strings.TrimSpace(inputUser)
	if inputUser == "" {
		inputUser = defaultUser
	}
	cfg.AccessKey = inputUser

	// 4. Password / Secret Key
	defaultPass := defaultCfg.SecretKey
	if defaultPass == "" || defaultPass == DefaultSecretKey {
		defaultPass = "minioadmin123"
	}
	fmt.Printf("  Enter Password / Secret Key [default: %s]: ", defaultPass)
	inputPass, _ := reader.ReadString('\n')
	inputPass = strings.TrimSpace(inputPass)
	if inputPass == "" {
		inputPass = defaultPass
	}
	cfg.SecretKey = inputPass

	fmt.Println("==================================================================")
	ui.Info("Target: %s%s%s | User: %s%s%s", ui.ColorCyan, cfg.Endpoint, ui.ColorReset, ui.ColorBold, cfg.AccessKey, ui.ColorReset)
	fmt.Println()

	return cfg
}

// ResolveServerConfig handles saved config, flag updates, and interactive setup:
// - If flags are supplied (--host, --port, --user, --password): updates saved config & runs immediately.
// - If reconfigure is true (--reconfigure): interactively re-prompts and updates saved config.
// - If saved config already exists: uses it immediately WITHOUT PROMPTING.
// - If first time run: interactively prompts once, saves to ~/.pilti/s3_config.json, and runs.
func ResolveServerConfig(host string, port int, user, password string, autoYes bool, reconfigure bool) *Config {
	activeCfg := GetActiveConfig()

	// 1. If explicit flags were passed, update configuration immediately without prompting
	if host != "" || port != 0 || user != "" || password != "" {
		cfg := &Config{
			Alias:     activeCfg.Alias,
			Endpoint:  activeCfg.Endpoint,
			AccessKey: activeCfg.AccessKey,
			SecretKey: activeCfg.SecretKey,
		}

		targetHost := host
		if targetHost == "" {
			targetHost, _ = extractHostAndPort(activeCfg.Endpoint)
		}
		targetPort := port
		if targetPort == 0 {
			_, targetPort = extractHostAndPort(activeCfg.Endpoint)
		}
		if targetPort == 0 {
			targetPort = 9000
		}

		cfg.Endpoint = FormatEndpoint(targetHost, targetPort)
		if user != "" {
			cfg.AccessKey = user
		}
		if password != "" {
			cfg.SecretKey = password
		}

		_ = SaveConfig(cfg)
		// Register with mc alias
		mcPath := FindMCExecutable()
		if mcPath != "" {
			_ = SetAlias(mcPath, cfg)
		}

		ui.Success("Updated S3 configuration: Target=%s | User=%s", cfg.Endpoint, cfg.AccessKey)
		return cfg
	}

	// 2. If reconfigure requested explicitly, run interactive prompt
	if reconfigure {
		cfg := PromptServerConfig(nil, false, activeCfg)
		_ = SaveConfig(cfg)
		mcPath := FindMCExecutable()
		if mcPath != "" {
			_ = SetAlias(mcPath, cfg)
		}
		ui.Success("Saved new S3 configuration successfully!")
		return cfg
	}

	// 3. If a saved configuration already exists, USE IT DIRECTLY (no prompt!)
	if HasSavedConfig() {
		return activeCfg
	}

	// 4. First-time setup (no config saved yet)
	if autoYes {
		_ = SaveConfig(activeCfg)
		return activeCfg
	}

	// Interactive prompt for first run
	cfg := PromptServerConfig(nil, false, activeCfg)
	_ = SaveConfig(cfg)
	mcPath := FindMCExecutable()
	if mcPath != "" {
		_ = SetAlias(mcPath, cfg)
	}

	ui.Success("S3 configuration saved! Subsequent commands will use this connection automatically.")
	fmt.Printf("ℹ TIP: To update settings in the future, pass flags (--host, -u, -p, --port) or run 'pilti s3 setup'\n\n")

	return cfg
}
