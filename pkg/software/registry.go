package software

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PortConfig struct {
	Name        string `json:"name"`
	DefaultPort int    `json:"default_port"`
	Container   int    `json:"container"`
	Protocol    string `json:"protocol"`
}

type Software struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Category      string       `json:"category"`
	Description   string       `json:"description"`
	Version       string       `json:"version,omitempty"`
	DefaultPorts  []string     `json:"ports"`
	PortConfigs   []PortConfig `json:"port_configs"`
	ContainerName string       `json:"container_name"`
	Subdir        string       `json:"subdir"`
	Dependencies  []string     `json:"dependencies"`
	DefaultImage  string       `json:"default_image"`
	Aliases       []string     `json:"aliases"`
	EnvFile       string       `json:"env_file,omitempty"`
}

var Registry = []Software{
	{
		ID:            "tb-app",
		Name:          "ThingsBoard Core Application",
		Category:      "Core IoT",
		Description:   "ThingsBoard enterprise IoT server & device orchestrator",
		Version:       "v-4.1.2",
		DefaultPorts:  []string{"8080 (Web UI)", "1883 (MQTT Broker)"},
		PortConfigs: []PortConfig{
			{Name: "Web UI", DefaultPort: 8080, Container: 8080, Protocol: "tcp"},
			{Name: "MQTT Broker", DefaultPort: 1883, Container: 1883, Protocol: "tcp"},
		},
		ContainerName: "Thingsboard-test",
		Subdir:        "tb-app",
		Dependencies:  []string{"tb-db"},
		DefaultImage:  "piltismartsolutions/thingsboard-3.8.1:v-4.1.2",
		EnvFile:       ".tb.env",
	},
	{
		ID:            "tb-db",
		Name:          "TimescaleDB / PostgreSQL",
		Category:      "Database",
		Description:   "High-performance telemetry & relational time-series database",
		Version:       "pg17",
		DefaultPorts:  []string{"5432 (Postgres/Timescale)"},
		PortConfigs: []PortConfig{
			{Name: "TimescaleDB Storage", DefaultPort: 5432, Container: 5432, Protocol: "tcp"},
		},
		ContainerName: "tb-timescaledb",
		Subdir:        "tb-db",
		Dependencies:  nil,
		DefaultImage:  "timescale/timescaledb-ha:pg17",
	},
	{
		ID:            "tb-edge",
		Name:          "ThingsBoard Edge Gateway",
		Category:      "Edge Computing",
		Description:   "Local autonomous ThingsBoard Edge instance for remote sites",
		Version:       "3.9.1EDGE",
		DefaultPorts:  []string{"8082 (Web UI)", "1884 (MQTT Broker)"},
		PortConfigs: []PortConfig{
			{Name: "Edge Web UI", DefaultPort: 8082, Container: 8080, Protocol: "tcp"},
			{Name: "Edge MQTT Broker", DefaultPort: 1884, Container: 1883, Protocol: "tcp"},
		},
		ContainerName: "mytbedge",
		Subdir:        "tb-edge",
		Dependencies:  []string{"tb-app"},
		DefaultImage:  "thingsboard/tb-edge:3.9.1EDGE",
	},
	{
		ID:            "jenkins",
		Name:          "Jenkins CI/CD Automation",
		Category:      "DevOps & CI/CD",
		Description:   "Automated build, test, and release controller engine",
		Version:       "lts",
		DefaultPorts:  []string{"8085 (Web UI)", "50000 (Agent Listener)"},
		PortConfigs: []PortConfig{
			{Name: "Jenkins Web UI", DefaultPort: 8085, Container: 8080, Protocol: "tcp"},
			{Name: "Jenkins Agent Listener", DefaultPort: 50000, Container: 50000, Protocol: "tcp"},
		},
		ContainerName: "jenkins",
		Subdir:        "jenkins",
		Dependencies:  nil,
		DefaultImage:  "jenkins/jenkins:lts",
	},
	{
		ID:            "piltiservices",
		Name:          "PiltiSmart Microservices",
		Category:      "Backend Services",
		Description:   "Modular PiltiSmart specialized API backend services",
		Version:       "v7.10.7",
		DefaultPorts:  []string{"9000 (Web/API Gateway)"},
		PortConfigs: []PortConfig{
			{Name: "API Gateway", DefaultPort: 9000, Container: 80, Protocol: "tcp"},
		},
		ContainerName: "piltiservices-test",
		Subdir:        "piltiservices",
		Dependencies:  nil,
		DefaultImage:  "piltismartsolutions/piltiservices:v7.10.7",
		EnvFile:       ".piltiservices.env",
	},
	{
		ID:            "kafka",
		Name:          "Apache Kafka Broker",
		Category:      "Message Streaming",
		Description:   "KRaft-based distributed event streaming & message broker",
		Version:       "4.1.1",
		DefaultPorts:  []string{"9092 (PLAINTEXT Broker)"},
		PortConfigs: []PortConfig{
			{Name: "Kafka PLAINTEXT Broker", DefaultPort: 9092, Container: 9092, Protocol: "tcp"},
		},
		ContainerName: "kafka",
		Subdir:        "kafka",
		Dependencies:  nil,
		DefaultImage:  "apache/kafka:4.1.1",
	},
	{
		ID:            "pulseX",
		Name:          "PulseX Cloud Gateway",
		Category:      "Cloud Platform",
		Description:   "Hybrid cloud synchronization connector and remote tunnel (PulseX / PMX)",
		Version:       "v8.4.41",
		DefaultPorts:  []string{"8088 (PulseX Cloud Gateway)"},
		PortConfigs: []PortConfig{
			{Name: "PulseX Gateway", DefaultPort: 8088, Container: 80, Protocol: "tcp"},
		},
		ContainerName: "pulseX",
		Subdir:        "pulsex",
		Dependencies:  nil,
		DefaultImage:  "piltismartsolutions/pilticloud:v8.4.41",
		Aliases:       []string{"pulsex", "pilticloud", "pmx", "cloud", "pilti-cloud"},
		EnvFile:       ".pmx.env",
	},
	{
		ID:            "minio",
		Name:          "MinIO Object Storage",
		Category:      "Cloud Storage",
		Description:   "High-performance S3-compatible distributed object storage server",
		Version:       "latest",
		DefaultPorts:  []string{"9000 (S3 API)", "9001 (Web Console)"},
		PortConfigs: []PortConfig{
			{Name: "MinIO S3 API", DefaultPort: 9000, Container: 9000, Protocol: "tcp"},
			{Name: "MinIO Web Console", DefaultPort: 9001, Container: 9001, Protocol: "tcp"},
		},
		ContainerName: "minio",
		Subdir:        "minio",
		Dependencies:  nil,
		DefaultImage:  "quay.io/minio/minio:latest",
		Aliases:       []string{"minio-server", "s3-server", "minio-storage"},
	},
}

// GetSoftware looks up a software by ID (case-insensitive and alias-tolerant).
func GetSoftware(id string) (*Software, bool) {
	norm := strings.ToLower(strings.TrimSpace(id))
	// Support common aliases
	if norm == "tb" {
		norm = "tb-app"
	} else if norm == "edge" || norm == "edge-tb" {
		norm = "tb-edge"
	} else if norm == "db" || norm == "timescale" {
		norm = "tb-db"
	} else if norm == "piltiservice" || norm == "services" {
		norm = "piltiservices"
	} else if norm == "pulsex" || norm == "pilticloud" || norm == "pmx" || norm == "cloud" || norm == "pilti-cloud" {
		norm = "pulsex"
	} else if norm == "kafka-broker" || norm == "broker" || norm == "apache-kafka" {
		norm = "kafka"
	} else if norm == "ci" || norm == "pilti-jenkins" {
		norm = "jenkins"
	} else if norm == "minio" || norm == "minio-server" || norm == "s3-server" || norm == "minio-storage" {
		norm = "minio"
	}

	for _, s := range Registry {
		if strings.ToLower(s.ID) == norm {
			return &s, true
		}
		for _, a := range s.Aliases {
			if strings.ToLower(a) == norm {
				return &s, true
			}
		}
	}
	return nil, false
}

// CheckStatus returns RUNNING, STOPPED, or NOT INSTALLED for a given software.
func CheckStatus(s *Software) string {
	// Query docker inspect for container status
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", s.ContainerName)
	out, err := cmd.Output()
	if err != nil {
		// For pulseX, also check legacy container name piltiCloud
		if s.ID == "pulseX" {
			cmdOld := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", "piltiCloud")
			if outOld, errOld := cmdOld.Output(); errOld == nil {
				out = outOld
			} else {
				return "NOT INSTALLED"
			}
		} else {
			return "NOT INSTALLED"
		}
	}
	status := strings.TrimSpace(string(out))
	switch status {
	case "running":
		return "RUNNING"
	case "exited", "paused", "dead":
		return fmt.Sprintf("STOPPED (%s)", status)
	case "":
		return "NOT INSTALLED"
	default:
		return strings.ToUpper(status)
	}
}

// GetVersion returns the runtime container version or configured default version for a software.
func GetVersion(s *Software) string {
	// 1. Inspect live container image tag
	cmd := exec.Command("docker", "inspect", "--format", "{{.Config.Image}}", s.ContainerName)
	if out, err := cmd.Output(); err == nil {
		img := strings.TrimSpace(string(out))
		if img != "" {
			parts := strings.Split(img, ":")
			if len(parts) > 1 {
				return parts[len(parts)-1]
			}
			return img
		}
	}

	// For pulseX, also check legacy container name piltiCloud
	if s.ID == "pulseX" {
		cmdOld := exec.Command("docker", "inspect", "--format", "{{.Config.Image}}", "piltiCloud")
		if outOld, errOld := cmdOld.Output(); errOld == nil {
			img := strings.TrimSpace(string(outOld))
			if img != "" {
				parts := strings.Split(img, ":")
				if len(parts) > 1 {
					return parts[len(parts)-1]
				}
				return img
			}
		}
	}

	// 2. Check docker-compose.yml in default directory if present
	targetDir := filepath.Join("/opt/piltismart", s.Subdir)
	composePath := filepath.Join(targetDir, "docker-compose.yml")
	if data, err := os.ReadFile(composePath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "image:") {
				parts := strings.Split(trimmed, ":")
				if len(parts) >= 3 {
					return strings.TrimSpace(parts[2])
				} else if len(parts) == 2 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// 3. Fallback to registry Version or DefaultImage tag
	if s.Version != "" {
		return s.Version
	}
	if s.DefaultImage != "" {
		parts := strings.Split(s.DefaultImage, ":")
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
	}
	return "latest"
}
