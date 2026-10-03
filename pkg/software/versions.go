package software

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

// FallbackPulseXVersions provides a curated, verified list of PulseX versions
// used if the GitHub or Docker APIs are temporarily unreachable.
var FallbackPulseXVersions = []string{
	"v8.4.41",
	"v8.4.38",
	"v8.4.37",
	"v8.4.36",
	"v8.4.35",
	"v8.4.34",
	"v8.4.4",
	"v8.4.3",
	"v8.4.2",
	"latest",
}

type gitHubTag struct {
	Name string `json:"name"`
}

type dockerHubTagItem struct {
	Name string `json:"name"`
}

type dockerHubTagsResponse struct {
	Results []dockerHubTagItem `json:"results"`
}

// FetchPulseXVersions fetches all available PulseX (PiltiCloud) release versions
// from GitHub repository tags and Docker registry.
func FetchPulseXVersions() []string {
	var collected []string
	seen := make(map[string]bool)

	client := &http.Client{Timeout: 5 * time.Second}

	// 1. Fetch tags from GitHub (PiltiSmart/vmxService - source repo for PulseX / pilticloud)
	ghURL := "https://api.github.com/repos/PiltiSmart/vmxService/tags"
	req, err := http.NewRequest("GET", ghURL, nil)
	if err == nil {
		req.Header.Set("User-Agent", "pilti-cli")
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			req.Header.Set("Authorization", "token "+token)
		}
		resp, gErr := client.Do(req)
		if gErr == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			var tags []gitHubTag
			if jsonErr := json.Unmarshal(body, &tags); jsonErr == nil {
				for _, t := range tags {
					tag := strings.TrimSpace(t.Name)
					if tag != "" && !seen[tag] {
						seen[tag] = true
						collected = append(collected, tag)
					}
				}
			}
		}
	}

	// 2. Fetch tags from Docker Hub (piltismartsolutions/pilticloud)
	dhURL := "https://hub.docker.com/v2/repositories/piltismartsolutions/pilticloud/tags?page_size=50"
	dhResp, dhErr := client.Get(dhURL)
	if dhErr == nil && dhResp.StatusCode == http.StatusOK {
		defer dhResp.Body.Close()
		body, _ := io.ReadAll(dhResp.Body)
		var dhTags dockerHubTagsResponse
		if jsonErr := json.Unmarshal(body, &dhTags); jsonErr == nil {
			for _, item := range dhTags.Results {
				tag := strings.TrimSpace(item.Name)
				if tag != "" && !seen[tag] {
					seen[tag] = true
					collected = append(collected, tag)
				}
			}
		}
	}

	// 3. Fallback or merge with known versions
	if len(collected) == 0 {
		return FallbackPulseXVersions
	}

	// Ensure top default v8.4.41 is at index 0 if present
	for i, v := range collected {
		if v == "v8.4.41" && i > 0 {
			collected[0], collected[i] = collected[i], collected[0]
			break
		}
	}

	return collected
}

// PromptPulseXVersion interactively prompts the user to select an available PulseX version.
func PromptPulseXVersion(reader *bufio.Reader) (string, error) {
	fmt.Println()
	ui.PrintBanner("PulseX Version Selection")
	ui.Info("Connecting to GitHub to fetch available PulseX release versions...")

	versions := FetchPulseXVersions()
	defaultVersion := "v8.4.41"
	if len(versions) > 0 {
		defaultVersion = versions[0]
	}

	fmt.Println("\nAvailable PulseX versions:")
	fmt.Println(strings.Repeat("-", 65))
	for i, v := range versions {
		label := v
		if v == defaultVersion {
			label = fmt.Sprintf("%-12s [Latest / Recommended]", v)
		}
		fmt.Printf(" [%d] %s\n", i+1, label)
	}
	fmt.Println(" [c] Custom version (enter your own version/tag)")
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("Select PulseX version [1-%d] (Press Enter for default: %s): ", len(versions), defaultVersion)

	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultVersion, nil
	}
	input = strings.TrimSpace(input)

	if input == "" {
		ui.Success("Selected default version: %s", defaultVersion)
		return defaultVersion, nil
	}

	if strings.ToLower(input) == "c" {
		fmt.Print("Enter custom PulseX version tag (e.g., v8.4.41 or latest): ")
		customTag, err := reader.ReadString('\n')
		if err != nil {
			return defaultVersion, nil
		}
		customTag = strings.TrimSpace(customTag)
		if customTag == "" {
			return defaultVersion, nil
		}
		ui.Success("Selected custom version: %s", customTag)
		return customTag, nil
	}

	idx, err := strconv.Atoi(input)
	if err == nil && idx >= 1 && idx <= len(versions) {
		selected := versions[idx-1]
		ui.Success("Selected version: %s", selected)
		return selected, nil
	}

	ui.Success("Selected version: %s", input)
	return input, nil
}

// GetAvailableSoftwareVersions returns known or available versions/tags for any given software.
func GetAvailableSoftwareVersions(s *Software) []string {
	norm := strings.ToLower(s.ID)
	switch norm {
	case "pulsex", "pilticloud":
		return FetchPulseXVersions()
	case "jenkins":
		return []string{"lts", "latest", "lts-jdk17", "lts-jdk21", "2.479.1"}
	case "tb-app":
		return []string{"v-4.1.2", "v-4.1.1", "3.8.1", "latest"}
	case "tb-db":
		return []string{"pg17", "pg16", "latest"}
	case "tb-edge":
		return []string{"3.9.1EDGE", "3.8.0EDGE", "latest"}
	case "piltiservices":
		return []string{"v7.10.7", "v7.10.6", "latest"}
	case "kafka":
		return []string{"4.1.1", "3.9.0", "3.8.0", "latest"}
	case "minio", "minio-server", "s3-server":
		return []string{"latest", "RELEASE.2025-07-16T15-35-03Z", "RELEASE.2024-11-07T00-52-28Z"}
	default:
		if s.Version != "" {
			return []string{s.Version, "latest"}
		}
		return []string{"latest"}
	}
}

// PromptSoftwareVersion interactively prompts for the version for ANY software component.
func PromptSoftwareVersion(s *Software, reader *bufio.Reader, autoYes bool) (string, error) {
	defaultVersion := s.Version
	if defaultVersion == "" {
		defaultVersion = "latest"
	}

	if autoYes {
		return defaultVersion, nil
	}

	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	norm := strings.ToLower(s.ID)
	// If PulseX, use the multi-choice catalog with dynamic GitHub versions
	if norm == "pulsex" || norm == "pilticloud" {
		return PromptPulseXVersion(reader)
	}

	available := GetAvailableSoftwareVersions(s)

	fmt.Println()
	fmt.Println("==================================================================")
	ui.PrintBanner(fmt.Sprintf("Version Selection: %s", s.Name))
	fmt.Printf("  • Default / Recommended Version: %s%s%s\n", ui.ColorGreen, defaultVersion, ui.ColorReset)
	if len(available) > 1 {
		fmt.Printf("  • Known Versions: %s\n", strings.Join(available, ", "))
	}
	fmt.Printf("  • Do you want to use the default version (%s)? [Y/n] (or enter custom version/tag): ", defaultVersion)

	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultVersion, nil
	}
	input = strings.TrimSpace(input)

	if input == "" || strings.ToLower(input) == "y" || strings.ToLower(input) == "yes" {
		ui.Success("Selected default version: %s", defaultVersion)
		return defaultVersion, nil
	}

	if strings.ToLower(input) == "n" || strings.ToLower(input) == "no" {
		for {
			fmt.Printf("  Enter custom version/tag for %s: ", s.Name)
			custom, err := reader.ReadString('\n')
			if err != nil {
				return defaultVersion, nil
			}
			custom = strings.TrimSpace(custom)
			if custom != "" {
				ui.Success("Selected custom version: %s", custom)
				return custom, nil
			}
			fmt.Printf("  %sVersion cannot be empty. Please enter a valid tag or version.%s\n", ui.ColorRed, ui.ColorReset)
		}
	}

	// User directly typed custom version tag
	ui.Success("Selected version: %s", input)
	return input, nil
}
