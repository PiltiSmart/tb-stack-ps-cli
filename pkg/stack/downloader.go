package stack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
)

const DefaultCatalogRepo = "https://raw.githubusercontent.com/PiltiSmart/stack-catalog/main"

type CatalogStack struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

type CatalogSchema struct {
	Version      string         `json:"version"`
	Organization string         `json:"organization"`
	UpdatedAt    string         `json:"updated_at"`
	Stacks       []CatalogStack `json:"stacks"`
}

func normalizeRepoURL(repo string) string {
	repo = strings.TrimSpace(repo)
	if repo == "" {
		return DefaultCatalogRepo
	}
	repo = strings.TrimSuffix(repo, "/")
	repo = strings.TrimSuffix(repo, ".git")

	if strings.Contains(repo, "github.com/") && !strings.Contains(repo, "raw.githubusercontent.com") {
		parts := strings.Split(repo, "github.com/")
		if len(parts) == 2 {
			repoPath := parts[1]
			if strings.Contains(repoPath, "/tree/") {
				repoPath = strings.Replace(repoPath, "/tree/", "/", 1)
				return fmt.Sprintf("https://raw.githubusercontent.com/%s", repoPath)
			}
			return fmt.Sprintf("https://raw.githubusercontent.com/%s/main", repoPath)
		}
	}
	return repo
}

func FetchRemoteCatalog(repoBase string) (*CatalogSchema, error) {
	repoBase = normalizeRepoURL(repoBase)
	url := fmt.Sprintf("%s/catalog.json", repoBase)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to contact catalog repo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("repository returned HTTP %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cat CatalogSchema
	if err := json.Unmarshal(body, &cat); err != nil {
		return nil, fmt.Errorf("failed to parse catalog.json: %w", err)
	}

	return &cat, nil
}

func DownloadStackManifests(repoBase, stackID, deployDir string) error {
	repoBase = normalizeRepoURL(repoBase)

	ui.Info("Fetching '%s' dynamic compose manifests from GitHub repository...", stackID)
	ui.Info("Repo URL: %s", repoBase)

	client := &http.Client{Timeout: 15 * time.Second}

	manifestFiles := []struct {
		RelRemotePath string
		LocalSubdir   string
		LocalFilename string
	}{
		{"stacks/tb-stack/tb-db/docker-compose.yml", "tb-db", "docker-compose.yml"},
		{"stacks/tb-stack/tb/docker-compose.yml", "tb", "docker-compose.yml"},
		{"stacks/tb-stack/tb/.tb.env.template", "tb", ".tb.env"},
		{"stacks/tb-stack/edge-tb/docker-compose.yml", "edge-tb", "docker-compose.yml"},
	}

	for _, m := range manifestFiles {
		remoteURL := fmt.Sprintf("%s/%s", repoBase, m.RelRemotePath)
		resp, err := client.Get(remoteURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				resp.Body.Close()
			}
			return fmt.Errorf("failed to download %s (HTTP %v): %w", remoteURL, resp.Status, err)
		}

		content, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read response for %s: %w", remoteURL, err)
		}

		localPath := filepath.Join(deployDir, m.LocalSubdir, m.LocalFilename)
		if err := os.WriteFile(localPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", localPath, err)
		}
		ui.Success("Downloaded: %s -> %s", m.RelRemotePath, localPath)
	}

	ui.Success("All manifests successfully fetched from GitHub repository!")
	return nil
}
