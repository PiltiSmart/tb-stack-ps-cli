package software

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const RawTbDbCompose = `version: '3.8'
services:
  timescaledb:
    image: timescale/timescaledb-ha:pg17
    container_name: tb-timescaledb
    restart: always
    environment:
      POSTGRES_PASSWORD: {{POSTGRES_PASSWORD}}
      POSTGRES_USER: {{POSTGRES_USER}}
      POSTGRES_DB: thingsboard
      PGDATA: /var/lib/postgresql/data/pgdata
    command: postgres -c shared_preload_libraries=pg_stat_statements,timescaledb
    ports:
      - "5432:5432"
    volumes:
      - ./timescale_data:/var/lib/postgresql/data
`

const RawTbAppCompose = `services:
  tb-app:
    image: piltismartsolutions/thingsboard-3.8.1:v-4.1.2
    container_name: Thingsboard-test
    ports:
      - "8080:8080"
      - "1883:1883"
      - "7070:7070"
    env_file:
      - ./.tb.env
    restart: always
`

const RawTbEnv = `INFISICAL_CLIENT_ID=your_infisical_client_id
INFISICAL_CLIENT_SECRET=your_infisical_client_secret
INFISICAL_PROJECT_ID=your_infisical_project_id
INFISICAL_URL=https://eu.infisical.com
INFISICAL_ENV=test
`

const RawTbEdgeCompose = `services:
  mytbedge:
    restart: always
    image: "thingsboard/tb-edge:3.9.1EDGE"
    container_name: mytbedge
    ports:
      - "${EDGE_WEB_PORT:-8082}:8080"
      - "${EDGE_MQTT_PORT:-1884}:1883"
      - "5683-5688:5683-5688/udp"
    environment:
      SPRING_DATASOURCE_URL: jdbc:postgresql://postgres:5432/tb-edge
      CLOUD_ROUTING_KEY: your_cloud_routing_key
      CLOUD_ROUTING_SECRET: your_cloud_routing_secret
      CLOUD_RPC_HOST: your_cloud_rpc_host
    volumes:
      - tb-edge-data:/data
      - tb-edge-logs:/var/log/tb-edge
  postgres:
    restart: always
    image: "postgres:15"
    ports:
      - "5433:5432"
    environment:
      POSTGRES_DB: tb-edge
      POSTGRES_PASSWORD: your_postgres_password
    volumes:
      - tb-edge-postgres-data:/var/lib/postgresql/data

volumes:
  tb-edge-data:
    name: tb-edge-data
  tb-edge-logs:
    name: tb-edge-logs
  tb-edge-postgres-data:
    name: tb-edge-postgres-data
`

const RawJenkinsCompose = `services:
  jenkins:
    image: jenkins/jenkins:lts
    container_name: jenkins
    privileged: true
    user: root
    ports:
      - "8085:8080"
      - "50000:50000"
    volumes:
      - ./data:/var/jenkins_home
      - /var/run/docker.sock:/var/run/docker.sock
    restart: always
`

const RawPiltiServicesCompose = `services:
  piltiservice:
    image: piltismartsolutions/piltiservices:v7.10.7
    container_name: piltiservices-test
  #  security_opt:
   #   - apparmor:unconfined
    ports:
      - "9000:80"
    env_file:
      - ./.piltiservices.env
    volumes:
      - ./logs:/app/logs
    restart: always
    healthcheck:
      # Command to run. If it returns 200 OK, the container is healthy.
      test: ["CMD", "curl", "-f", "http://localhost:80/pilti/piltiUrls"]
      interval: 10s       # How often to check
      timeout: 5s         # How long to wait for a response before failing
      retries: 5          # How many consecutive failures mean it's "unhealthy"
      start_period: 45s
`

const RawPiltiServicesEnv = `INFISICAL_CLIENT_ID=your_infisical_client_id
INFISICAL_CLIENT_SECRET=your_infisical_client_secret
INFISICAL_PROJECT_ID=your_infisical_project_id
INFISICAL_SITE_URL=https://eu.infisical.com
INFISICAL_ENV=test
`

const RawKafkaCompose = `services:
  kafka:
    image: apache/kafka:4.1.1
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      KAFKA_NODE_ID: "1"
      KAFKA_PROCESS_ROLES: "broker,controller"
      KAFKA_CONTROLLER_QUORUM_VOTERS: "1@kafka:9093"
      KAFKA_LISTENERS: "PLAINTEXT://:9092,CONTROLLER://:9093"
      KAFKA_ADVERTISED_LISTENERS: "PLAINTEXT://localhost:9092"
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT"
      KAFKA_CONTROLLER_LISTENER_NAMES: "CONTROLLER"
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: "1"
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: "1"
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: "1"
      KAFKA_LOG_DIRS: "/var/lib/kafka/data"
    restart: always
    volumes:
      - kafka_data:/var/lib/kafka/data
volumes:
  kafka_data:
`

const RawPulseXCompose = `services:
  pulseX:
    image: piltismartsolutions/pilticloud:{{VERSION}}
    container_name: pulseX
    ports:
      - "{{PORT}}:80"
    env_file:
      - ./.pmx.env
#    working_dir: /app
#    volumes:
#      - ./logs:/app/logs
    restart: always
`

const RawPulseXEnv = `INFISICAL_CLIENT_ID=your_infisical_client_id
INFISICAL_CLIENT_SECRET=your_infisical_client_secret
INFISICAL_PROJECT_ID=your_infisical_project_id
INFISICAL_SITE_URL=https://eu.infisical.com
INFISICAL_ENV=test
KAFKA_SERVER=localhost:9092
`

const RawMinioCompose = `version: '3.7'

services:
  minio:
    image: quay.io/minio/minio:{{VERSION}}
    container_name: minio
    restart: always
    ports:
      - "{{API_PORT}}:9000"
      - "{{CONSOLE_PORT}}:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin123
    volumes:
      - ./minio/data:/data
    command: server /data --console-address ":9001"
`

// Aliases for backwards compatibility
const RawPiltiCloudCompose = RawPulseXCompose
const RawPiltiCloudEnv = RawPulseXEnv

// WriteTemplates provisions the configuration and compose files for a given software with default settings.
func WriteTemplates(s *Software, targetDir string) error {
	return WriteConfiguredTemplates(s, targetDir, nil, "")
}

// WriteConfiguredTemplates provisions the configuration and compose files with customized ports and versions.
func WriteConfiguredTemplates(s *Software, targetDir string, ports map[string]int, customVersion string) error {
	return WriteConfiguredTemplatesWithCreds(s, targetDir, ports, customVersion, "", "")
}

// WriteConfiguredTemplatesWithCreds provisions configuration and compose files with customized ports, versions, and database credentials.
func WriteConfiguredTemplatesWithCreds(s *Software, targetDir string, ports map[string]int, customVersion string, dbUser, dbPass string) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	normID := strings.ToLower(s.ID)

	switch normID {
	case "tb-app":
		compose := RawTbAppCompose
		if customVersion != "" {
			compose = strings.Replace(compose, ":v-4.1.2", fmt.Sprintf(":%s", customVersion), 1)
		}
		if p, ok := ports["Web UI"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"8080:8080\"", fmt.Sprintf("\"%d:8080\"", p), 1)
		}
		if p, ok := ports["MQTT Broker"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"1883:1883\"", fmt.Sprintf("\"%d:1883\"", p), 1)
		}
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(targetDir, ".tb.env"), []byte(RawTbEnv), 0644); err != nil {
			return err
		}
	case "tb-db":
		user := "postgres"
		if dbUser != "" {
			user = dbUser
		}
		pass := "postgres"
		if dbPass != "" {
			pass = dbPass
		}

		compose := RawTbDbCompose
		if customVersion != "" {
			compose = strings.Replace(compose, ":pg17", fmt.Sprintf(":%s", customVersion), 1)
		}
		if p, ok := ports["TimescaleDB Storage"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"5432:5432\"", fmt.Sprintf("\"%d:5432\"", p), 1)
		}
		compose = strings.Replace(compose, "{{POSTGRES_USER}}", user, 1)
		compose = strings.Replace(compose, "{{POSTGRES_PASSWORD}}", pass, 1)
		compose = strings.Replace(compose, "your_postgres_user", user, 1)
		compose = strings.Replace(compose, "your_postgres_password", pass, 1)

		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
		tsDir := filepath.Join(targetDir, "timescale_data")
		_ = os.MkdirAll(tsDir, 0777)
		_ = os.Chmod(tsDir, 0777)
	case "tb-edge":
		compose := RawTbEdgeCompose
		if customVersion != "" {
			compose = strings.Replace(compose, ":3.9.1EDGE", fmt.Sprintf(":%s", customVersion), 1)
		}
		if p, ok := ports["Edge Web UI"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"${EDGE_WEB_PORT:-8082}:8080\"", fmt.Sprintf("\"%d:8080\"", p), 1)
		}
		if p, ok := ports["Edge MQTT Broker"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"${EDGE_MQTT_PORT:-1884}:1883\"", fmt.Sprintf("\"%d:1883\"", p), 1)
		}
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
	case "jenkins":
		compose := RawJenkinsCompose
		if customVersion != "" {
			compose = strings.Replace(compose, ":lts", fmt.Sprintf(":%s", customVersion), 1)
		}
		if p, ok := ports["Jenkins Web UI"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"8085:8080\"", fmt.Sprintf("\"%d:8080\"", p), 1)
		}
		if p, ok := ports["Jenkins Agent Listener"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"50000:50000\"", fmt.Sprintf("\"%d:50000\"", p), 1)
		}
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
		dataDir := filepath.Join(targetDir, "data")
		_ = os.MkdirAll(dataDir, 0777)
		_ = os.Chmod(dataDir, 0777)
	case "piltiservices":
		compose := RawPiltiServicesCompose
		if customVersion != "" {
			compose = strings.Replace(compose, ":v7.10.7", fmt.Sprintf(":%s", customVersion), 1)
		}
		if p, ok := ports["API Gateway"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"9000:80\"", fmt.Sprintf("\"%d:80\"", p), 1)
		}
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(targetDir, ".piltiservices.env"), []byte(RawPiltiServicesEnv), 0644); err != nil {
			return err
		}
		logsDir := filepath.Join(targetDir, "logs")
		_ = os.MkdirAll(logsDir, 0777)
		_ = os.Chmod(logsDir, 0777)
	case "kafka":
		compose := RawKafkaCompose
		if customVersion != "" {
			compose = strings.Replace(compose, ":4.1.1", fmt.Sprintf(":%s", customVersion), 1)
		}
		if p, ok := ports["Kafka PLAINTEXT Broker"]; ok && p > 0 {
			compose = strings.Replace(compose, "\"9092:9092\"", fmt.Sprintf("\"%d:%d\"", p, p), 1)
			compose = strings.Replace(compose, "PLAINTEXT://:9092", fmt.Sprintf("PLAINTEXT://:%d", p), 1)
			compose = strings.Replace(compose, "PLAINTEXT://localhost:9092", fmt.Sprintf("PLAINTEXT://localhost:%d", p), 1)
		}
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
	case "pulsex", "pilticloud":
		version := "v8.4.41"
		if customVersion != "" {
			version = customVersion
		}
		port := 8088
		if p, ok := ports["PulseX Gateway"]; ok && p > 0 {
			port = p
		}

		compose := strings.Replace(RawPulseXCompose, "{{VERSION}}", version, 1)
		compose = strings.Replace(compose, "{{PORT}}", strconv.Itoa(port), 1)

		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(targetDir, ".pmx.env"), []byte(RawPulseXEnv), 0644); err != nil {
			return err
		}
	case "minio", "minio-server", "s3-server", "minio-storage":
		version := "latest"
		if customVersion != "" {
			version = customVersion
		}
		apiPort := 9000
		if p, ok := ports["MinIO S3 API"]; ok && p > 0 {
			apiPort = p
		}
		consolePort := 9001
		if p, ok := ports["MinIO Web Console"]; ok && p > 0 {
			consolePort = p
		}

		compose := strings.Replace(RawMinioCompose, "{{VERSION}}", version, 1)
		compose = strings.Replace(compose, "{{API_PORT}}", strconv.Itoa(apiPort), 1)
		compose = strings.Replace(compose, "{{CONSOLE_PORT}}", strconv.Itoa(consolePort), 1)

		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(compose), 0644); err != nil {
			return err
		}
		dataDir := filepath.Join(targetDir, "minio", "data")
		_ = os.MkdirAll(dataDir, 0777)
		_ = os.Chmod(dataDir, 0777)
	default:
		return fmt.Errorf("unknown software template ID: %s", s.ID)
	}

	return nil
}
