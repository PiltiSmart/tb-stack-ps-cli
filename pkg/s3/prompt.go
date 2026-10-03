package s3

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

// PromptServerConfig interactively prompts for server IP, port, username, and password.
func PromptServerConfig(reader *bufio.Reader, autoYes bool, defaultCfg *Config) *Config {
	if defaultCfg == nil {
		defaultCfg = GetDefaultConfig()
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

	fmt.Println()
	fmt.Println("==================================================================")
	fmt.Printf("%s%s[S3 / MinIO Server Connection Details]%s\n", ui.ColorBold, ui.ColorCyan, ui.ColorReset)
	fmt.Println("==================================================================")

	// 1. Server IP / Host
	defaultHost := "localhost"
	fmt.Printf("  Enter MinIO Server IP or Hostname [default: %s]: ", defaultHost)
	inputHost, _ := reader.ReadString('\n')
	inputHost = strings.TrimSpace(inputHost)
	if inputHost == "" {
		inputHost = defaultHost
	}

	// 2. Port (Display default port 1st, then ask)
	defaultPort := 9000
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

	// Format Endpoint (handle http / https)
	scheme := "http"
	cleanHost := inputHost
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
		cfg.Endpoint = fmt.Sprintf("%s://%s", scheme, cleanHost)
	} else {
		cfg.Endpoint = fmt.Sprintf("%s://%s:%d", scheme, cleanHost, finalPort)
	}

	// 3. Username / Access Key
	defaultUser := "minioadmin"
	fmt.Printf("\n  Enter Username / Access Key [default: %s]: ", defaultUser)
	inputUser, _ := reader.ReadString('\n')
	inputUser = strings.TrimSpace(inputUser)
	if inputUser == "" {
		inputUser = defaultUser
	}
	cfg.AccessKey = inputUser

	// 4. Password / Secret Key
	defaultPass := "minioadmin123"
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

// ResolveServerConfig handles interactive prompts or flag-based overrides.
func ResolveServerConfig(host string, port int, user, password string, autoYes bool) *Config {
	defaultCfg := GetDefaultConfig()

	// If host was passed via flag
	if host != "" {
		cfg := &Config{
			Alias:     defaultCfg.Alias,
			AccessKey: user,
			SecretKey: password,
		}
		if cfg.AccessKey == "" {
			cfg.AccessKey = "minioadmin"
		}
		if cfg.SecretKey == "" {
			cfg.SecretKey = "minioadmin123"
		}
		if port == 0 {
			port = 9000
		}
		scheme := "http"
		cleanHost := host
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
			cfg.Endpoint = fmt.Sprintf("%s://%s", scheme, cleanHost)
		} else {
			cfg.Endpoint = fmt.Sprintf("%s://%s:%d", scheme, cleanHost, port)
		}
		return cfg
	}

	// If autoYes is true and no host was specified, use default
	if autoYes {
		return defaultCfg
	}

	// Interactive prompt
	return PromptServerConfig(nil, false, defaultCfg)
}
