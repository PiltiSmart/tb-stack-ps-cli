# PiltiSmart Enterprise CLI (`pilti`) & Stack Catalog

A standalone, zero-dependency CLI written in **Go** using the **Cobra Framework** to discover, configure, inspect, and orchestrate all software components and microservices across the **PiltiSmart** ecosystem:

1. `tb-db` : TimescaleDB / PostgreSQL database (**Install 1st**)
2. `tb-app` : ThingsBoard Core Application (**Requires `tb-db`**)
3. `tb-edge` : ThingsBoard Edge Gateway (**Requires `tb-app`**)
4. `jenkins` : Jenkins CI/CD Automation Engine
5. `piltiservices` : PiltiSmart Microservices Backend
6. `kafka` : Apache Kafka Distributed Streaming Broker
7. `pulseX` : PulseX Cloud Gateway & Sync Tunnel (PMX - formerly `pilticloud`)

---

## ⚡ Quick 1-Line Installation (Linux & macOS)

Install the pre-compiled `pilti` binary on any fresh Linux VM, LXC container, or macOS terminal with a single command:

### Production / Main Branch:
```bash
curl -fsSL https://raw.githubusercontent.com/PiltiSmart/stack-catalog/main/install.sh | bash
```

### Feature Branch (`feat/add-enterprise-tools`):
```bash
curl -fsSL https://raw.githubusercontent.com/PiltiSmart/stack-catalog/feat/add-enterprise-tools/install.sh | bash
```

### Direct Static Binary Download (Manual):
If you cannot use `curl ... | bash`, download the static binary directly:
```bash
# Linux (amd64 / x86_64):
sudo curl -sSL https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-linux-amd64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti

# Linux (ARM64 / aarch64):
sudo curl -sSL https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-linux-arm64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti

# macOS (Apple Silicon M1/M2/M3/M4):
sudo curl -sSL https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-darwin-arm64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti
```

---

## 🚀 Available Stacks & Software Catalog

| Software ID | Name | Version | Default Ports | Dependencies | Description |
|---|---|---|---|---|---|
| **[`tb-db`](stacks/tb-db/)** | TimescaleDB / Postgres | `pg17` | `5432` | *None* (**Install 1st**) | High-performance telemetry & time-series DB |
| **[`tb-app`](stacks/tb-app/)** | ThingsBoard Core | `3.8.1` | `8080`, `1883`, `7070` | **`tb-db`** | ThingsBoard enterprise IoT server |
| **[`tb-edge`](stacks/tb-edge/)** | ThingsBoard Edge | `3.9.1EDGE` | `8082`, `1884`, `5683-5688/udp` | **`tb-app`** | Remote autonomous ThingsBoard Edge gateway |
| **[`jenkins`](stacks/jenkins/)** | Jenkins CI/CD | `lts` | `8085:8080`, `50000` | *None* | CI/CD build automation controller |
| **[`piltiservices`](stacks/piltiservices/)** | PiltiSmart Microservices | `v7.10.7` | `9000:80` | *None* | Specialized API microservices backend |
| **[`kafka`](stacks/kafka/)** | Apache Kafka Broker | `4.1.1` | `9092` | *None* | KRaft distributed event streaming broker |
| **[`pulseX`](stacks/pulsex/)** | PulseX Cloud Gateway | `v8.4.41` (dynamic selector) | `8088:80` | *None* | Hybrid cloud connector (PMX / formerly PiltiCloud) |

---

## 🛑 Installation Workflow & Dependency Hierarchy

ThingsBoard software components **must** be deployed in prerequisite order. If a prerequisite dependency is missing, `pilti` will halt and display a clear dependency error:

### Step 1: Deploy TimescaleDB Database (`tb-db`) — Install 1st
```bash
pilti install tb-db
# or: pilti tb-db install
```

### Step 2: Deploy ThingsBoard Core Application (`tb-app`) — Install 2nd
> **Prerequisite:** `tb-db` must be running!
```bash
pilti install tb-app
# or: pilti tb-app install
```
*If `tb-db` is not running, the installer will immediately abort with:*
```text
[✖] Cannot install 'ThingsBoard Core Application' (tb-app): Missing required prerequisite!
[✖] Prerequisite dependency 'TimescaleDB / PostgreSQL' (tb-db) is NOT running (Current status: NOT INSTALLED).
```

### Step 3: Deploy ThingsBoard Edge Gateway (`tb-edge`) — Install 3rd
> **Prerequisite:** `tb-app` (and `tb-db`) must be running!
```bash
pilti install tb-edge
# or: pilti tb-edge install
```

---

## 📦 Deploying Standalone Enterprise Tools

These services have no inter-dependencies and can be deployed at any time:

```bash
# Jenkins CI/CD Engine
pilti install jenkins

# Apache Kafka Event Streaming Broker
pilti install kafka

# PiltiSmart Microservices Backend
pilti install piltiservices

# PulseX Cloud Gateway (PMX - formerly PiltiCloud)
pilti install pulseX
# or: pilti install pilticloud
# or: pilti pulseX install
```

---

## 🛠️ CLI Management Commands

### 1. View Software Catalog & Live Status
```bash
pilti list
```

### 2. Pre-Flight Diagnostics
Checks Docker daemon, Docker Compose plugin, RAM, disk space, and network port availability:
```bash
pilti doctor
```

### 3. Check Service Status
```bash
pilti status            # Global status
pilti tb-app status     # Specific service status
pilti jenkins status
pilti kafka status
pilti pulseX status
```

### 4. Stream Service Logs
```bash
pilti tb-app logs -f
pilti jenkins logs -f
pilti kafka logs -f
pilti piltiservices logs -f
pilti pulseX logs -f
```

### 5. Restart, Stop, or Remove Services
```bash
pilti tb-db restart
pilti kafka restart
pilti pulseX restart
pilti pulseX stop
pilti jenkins remove
```

### 6. Interactive Terminal Selector
Launch an interactive menu to choose software and actions:
```bash
pilti install
```

### 7. Interactive Port Check & PulseX Version Selector
When installing any software, `pilti` now:
1. **Interactive Port Verification**: Displays default host port first, asks the tech user `[Y/n]` to confirm or enter a custom port, and actively tests live port availability on the target system.
2. **PulseX Version Selector**: Dynamically connects to GitHub / Docker registry to list all available versions (e.g. `v8.4.41`, `v8.4.38`, `v8.4.4`, `latest`) and allows selecting or inputting custom tags.
3. **Non-Interactive / Scripted Flags**: Supports `--port <number>`, `--version <tag>`, and `-y`/`--yes` for automated deployments:
```bash
pilti pulseX install --version v8.4.41 --port 8088 -y
pilti tb-app install --port 8080 -y
```

---

## 🪣 AWS S3 Cloud Storage CLI (`pilti s3`)

`pilti` includes an AWS S3-compatible cloud storage CLI powered by the **MinIO Client (`mc`)** engine.
It connects out-of-the-box to the PiltiSmart MinIO Object Storage server at `http://localhost:9000` with the `myminio` alias, allowing full CRUD operations using familiar AWS S3 CLI syntax.

### 1. Show S3 Overview & Available Commands
```bash
pilti s3
# or
pilti s3 --help
```

### 2. S3 Operations & CRUD Examples

| Operation | Command | Description |
|---|---|---|
| **List Buckets** | `pilti s3 ls` | Lists all buckets on MinIO |
| **List Objects** | `pilti s3 ls s3://mybucket` | Lists objects inside `mybucket` |
| **Make Bucket** | `pilti s3 mb s3://mybucket` | Creates a new S3 bucket |
| **Remove Bucket** | `pilti s3 rb s3://mybucket [--force]` | Deletes an S3 bucket |
| **Upload File** | `pilti s3 cp ./data.csv s3://mybucket/` | Copies local file to S3 |
| **Download File** | `pilti s3 cp s3://mybucket/data.csv ./` | Copies S3 object to local disk |
| **Recursive Copy** | `pilti s3 cp ./dist s3://mybucket/dist -r` | Recursively uploads directory |
| **Read Object** | `pilti s3 cat s3://mybucket/config.json` | Displays object content to stdout |
| **Object Metadata** | `pilti s3 stat s3://mybucket/data.csv` | Displays size, ETags, and metadata |
| **Sync / Mirror** | `pilti s3 sync ./backup s3://mybucket/backup` | Synchronizes directory with S3 |
| **Delete Object** | `pilti s3 rm s3://mybucket/data.csv` | Removes an object from S3 |
| **Setup & Verify** | `pilti s3 setup` | Verifies `mc` binary and registers alias |
| **Configuration** | `pilti s3 config` | Displays active endpoint & credentials |

