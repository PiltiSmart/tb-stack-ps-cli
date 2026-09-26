package software

import (
	"fmt"
	"os/exec"
	"strings"
)

type Software struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Category      string   `json:"category"`
	Description   string   `json:"description"`
	DefaultPorts  []string `json:"ports"`
	ContainerName string   `json:"container_name"`
	Subdir        string   `json:"subdir"`
	Dependencies  []string `json:"dependencies"`
}

var Registry = []Software{
	{
		ID:            "tb-app",
		Name:          "ThingsBoard Core Application",
		Category:      "Core IoT",
		Description:   "ThingsBoard enterprise IoT server & device orchestrator",
		DefaultPorts:  []string{"8080 (Web UI)", "1883 (MQTT Broker)"},
		ContainerName: "Thingsboard-test",
		Subdir:        "tb-app",
		Dependencies:  []string{"tb-db"},
	},
	{
		ID:            "tb-db",
		Name:          "TimescaleDB / PostgreSQL",
		Category:      "Database",
		Description:   "High-performance telemetry & relational time-series database",
		DefaultPorts:  []string{"5432 (Postgres/Timescale)"},
		ContainerName: "tb-timescaledb",
		Subdir:        "tb-db",
		Dependencies:  nil,
	},
	{
		ID:            "tb-edge",
		Name:          "ThingsBoard Edge Gateway",
		Category:      "Edge Computing",
		Description:   "Local autonomous ThingsBoard Edge instance for remote sites",
		DefaultPorts:  []string{"8082 (Web UI)", "1884 (MQTT Broker)"},
		ContainerName: "mytbedge",
		Subdir:        "tb-edge",
		Dependencies:  []string{"tb-app"},
	},
	{
		ID:            "jenkins",
		Name:          "Jenkins CI/CD Automation",
		Category:      "DevOps & CI/CD",
		Description:   "Automated build, test, and release controller engine",
		DefaultPorts:  []string{"80:8080 (Web UI)", "50000 (Agent Listener)"},
		ContainerName: "jenkins",
		Subdir:        "jenkins",
		Dependencies:  nil,
	},
	{
		ID:            "piltiservices",
		Name:          "PiltiSmart Microservices",
		Category:      "Backend Services",
		Description:   "Modular PiltiSmart specialized API backend services",
		DefaultPorts:  []string{"80 (Web/API)"},
		ContainerName: "piltiservices-test",
		Subdir:        "piltiservices",
		Dependencies:  nil,
	},
	{
		ID:            "kafka",
		Name:          "Apache Kafka Broker",
		Category:      "Message Streaming",
		Description:   "KRaft-based distributed event streaming & message broker",
		DefaultPorts:  []string{"9092 (PLAINTEXT Broker)"},
		ContainerName: "kafka",
		Subdir:        "kafka",
		Dependencies:  nil,
	},
	{
		ID:            "pilticloud",
		Name:          "PiltiSmart Cloud Gateway",
		Category:      "Cloud Platform",
		Description:   "Hybrid cloud synchronization connector and remote tunnel",
		DefaultPorts:  []string{"80 (Cloud Service / PMX)"},
		ContainerName: "piltiCloud",
		Subdir:        "pilticloud",
		Dependencies:  nil,
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
	} else if norm == "pmx" || norm == "cloud" || norm == "pilti-cloud" {
		norm = "pilticloud"
	} else if norm == "kafka-broker" || norm == "broker" || norm == "apache-kafka" {
		norm = "kafka"
	} else if norm == "ci" || norm == "pilti-jenkins" {
		norm = "jenkins"
	}

	for _, s := range Registry {
		if strings.ToLower(s.ID) == norm {
			return &s, true
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
		return "NOT INSTALLED"
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
