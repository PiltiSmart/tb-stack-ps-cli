package stack

const RawTbDbCompose = `version: '3.8'
services:
  timescaledb:
    image: timescale/timescaledb-ha:pg17
    container_name: Test-timescaledb
    restart: always
    environment:
      POSTGRES_PASSWORD: your_postgres_password
      POSTGRES_USER: your_postgres_user
      POSTGRES_DB: thingsboard
      PGDATA: /var/lib/postgresql/data/pgdata
    command: postgres -c shared_preload_libraries=pg_stat_statements,timescaledb
    ports:
      - "5432:5432"
    volumes:
      - ./timescale_data:/var/lib/postgresql/data
`

const RawTbCompose = `version: '3.8'

services:
  kabiPayment:
    image: piltismartsolutions/thingsboard-3.8.1:v-4.1.2
    container_name: Thingsboard-test
    ports:
      - "80:80"
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

const RawEdgeTbCompose = `version: '3.8'
services:
  mytbedge:
    restart: always
    image: "thingsboard/tb-edge:3.9.1EDGE"
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
      - "5432"
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
