package stack

const RawTbDbCompose = `version: '3.8'
services:
  timescaledb:
    image: timescale/timescaledb-ha:pg17
    container_name: Test-timescaledb
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

const RawTbEnv = `INFISICAL_CLIENT_ID=0c4a1a06-34d9-4500-b900-3ceed9aaa660
INFISICAL_CLIENT_SECRET=d0d09ace072d03670826f4d99673f8b6469f61e39470f15103f0510eea3e89ae
INFISICAL_PROJECT_ID=47146ccc-417e-4b8b-9a4c-443e2534b82e
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
      - "5432"
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
