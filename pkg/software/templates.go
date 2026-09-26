package software

import (
	"fmt"
	"os"
	"path/filepath"
)

const RawTbDbCompose = `version: '3.8'
services:
  timescaledb:
    image: timescale/timescaledb-ha:pg17
    container_name: tb-timescaledb
    restart: always
    environment:
      POSTGRES_PASSWORD: qwer1234
      POSTGRES_USER: postgres
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

const RawTbEnv = `INFISICAL_CLIENT_ID=0c4a1a06-34d9-4500-b900-3ceed9aaa660
INFISICAL_CLIENT_SECRET=d0d09ace072d03670826f4d99673f8b6469f61e39470f15103f0510eea3e89ae
INFISICAL_PROJECT_ID=47146ccc-417e-4b8b-9a4c-443e2534b82e
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
      CLOUD_ROUTING_KEY: b07812d7-4641-cfcb-ea03-02d39079a1eb
      CLOUD_ROUTING_SECRET: x2kw2qewx262369fnybk
      CLOUD_RPC_HOST: 192.168.0.126
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
      POSTGRES_PASSWORD: postgres
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
      test: ["CMD", "curl", "-f", "http://127.0.0.1:80/pilti/piltiUrls"]
      interval: 10s       # How often to check
      timeout: 5s         # How long to wait for a response before failing
      retries: 5          # How many consecutive failures mean it's "unhealthy"
      start_period: 45s
`

const RawPiltiServicesEnv = `INFISICAL_CLIENT_ID=0c4a1a06-34d9-4500-b900-3ceed9aaa660
INFISICAL_CLIENT_SECRET=d0d09ace072d03670826f4d99673f8b6469f61e39470f15103f0510eea3e89ae
INFISICAL_PROJECT_ID=47146ccc-417e-4b8b-9a4c-443e2534b82e
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
      KAFKA_ADVERTISED_LISTENERS: "PLAINTEXT://10.70.70.4:9092"
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

const RawPiltiCloudCompose = `services:
  piltiCloud:
    image: piltismartsolutions/pilticloud:v8.4.41
    container_name: piltiCloud
    ports:
      - "8088:80"
    env_file:
      - ./.pmx.env
#    working_dir: /app
#    volumes:
#      - ./logs:/app/logs
    restart: always
`

const RawPiltiCloudEnv = `INFISICAL_CLIENT_ID=0c4a1a06-34d9-4500-b900-3ceed9aaa660
INFISICAL_CLIENT_SECRET=d0d09ace072d03670826f4d99673f8b6469f61e39470f15103f0510eea3e89ae
INFISICAL_PROJECT_ID=47146ccc-417e-4b8b-9a4c-443e2534b82e
INFISICAL_SITE_URL=https://eu.infisical.com
INFISICAL_ENV=test
KAFKA_SERVER=10.70.70.4:9092
`

// WriteTemplates provisions the configuration and compose files for a given software.
func WriteTemplates(s *Software, targetDir string) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	switch s.ID {
	case "tb-app":
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(RawTbAppCompose), 0644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(targetDir, ".tb.env"), []byte(RawTbEnv), 0644); err != nil {
			return err
		}
	case "tb-db":
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(RawTbDbCompose), 0644); err != nil {
			return err
		}
		tsDir := filepath.Join(targetDir, "timescale_data")
		_ = os.MkdirAll(tsDir, 0777)
		_ = os.Chmod(tsDir, 0777)
	case "tb-edge":
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(RawTbEdgeCompose), 0644); err != nil {
			return err
		}
	case "jenkins":
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(RawJenkinsCompose), 0644); err != nil {
			return err
		}
		dataDir := filepath.Join(targetDir, "data")
		_ = os.MkdirAll(dataDir, 0777)
		_ = os.Chmod(dataDir, 0777)
	case "piltiservices":
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(RawPiltiServicesCompose), 0644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(targetDir, ".piltiservices.env"), []byte(RawPiltiServicesEnv), 0644); err != nil {
			return err
		}
		logsDir := filepath.Join(targetDir, "logs")
		_ = os.MkdirAll(logsDir, 0777)
		_ = os.Chmod(logsDir, 0777)
	case "kafka":
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(RawKafkaCompose), 0644); err != nil {
			return err
		}
	case "pilticloud":
		if err := os.WriteFile(filepath.Join(targetDir, "docker-compose.yml"), []byte(RawPiltiCloudCompose), 0644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(targetDir, ".pmx.env"), []byte(RawPiltiCloudEnv), 0644); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown software template ID: %s", s.ID)
	}

	return nil
}
