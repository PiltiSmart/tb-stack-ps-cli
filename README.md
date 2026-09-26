# PiltiSmart Enterprise CLI (`pilti`) & Stack Catalog

A standalone, zero-dependency CLI written in **Go** using the **Cobra Framework** to discover, configure, inspect, and orchestrate all software components and microservices across the **PiltiSmart** ecosystem:

1. `tb-db` : TimescaleDB / PostgreSQL database (**Install 1st**)
2. `tb-app` : ThingsBoard Core Application (**Requires `tb-db`**)
3. `tb-edge` : ThingsBoard Edge Gateway (**Requires `tb-app`**)
4. `jenkins` : Jenkins CI/CD Automation Engine
5. `piltiservices` : PiltiSmart Microservices Backend
6. `kafka` : Apache Kafka Distributed Streaming Broker
7. `pilticloud` : PiltiSmart Cloud Gateway & Sync Tunnel (PMX)

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
| **[`tb-app`](stacks/tb-app/)** | ThingsBoard Core | `3.8.1` | `80`, `8080`, `1883`, `7070` | **`tb-db`** | ThingsBoard enterprise IoT server |
| **[`tb-edge`](stacks/tb-edge/)** | ThingsBoard Edge | `3.9.1EDGE` | `8082`, `1884`, `5683-5688/udp` | **`tb-app`** | Remote autonomous ThingsBoard Edge gateway |
| **[`jenkins`](stacks/jenkins/)** | Jenkins CI/CD | `lts` | `80:8080`, `50000` | *None* | CI/CD build automation controller |
| **[`piltiservices`](stacks/piltiservices/)** | PiltiSmart Microservices | `v7.10.7` | `80` | *None* | Specialized API microservices backend |
| **[`kafka`](stacks/kafka/)** | Apache Kafka Broker | `4.1.1` | `9092` | *None* | KRaft distributed event streaming broker |
| **[`pilticloud`](stacks/pilticloud/)** | PiltiSmart Cloud Gateway | `v8.4.41` | `80` | *None* | Hybrid cloud synchronization connector (PMX) |

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

# PiltiSmart Cloud Gateway (PMX)
pilti install pilticloud
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
```

### 4. Stream Service Logs
```bash
pilti tb-app logs -f
pilti jenkins logs -f
pilti kafka logs -f
pilti piltiservices logs -f
```

### 5. Restart, Stop, or Remove Services
```bash
pilti tb-db restart
pilti kafka restart
pilti pilticloud stop
pilti jenkins remove
```

### 6. Interactive Terminal Selector
Launch an interactive menu to choose software and actions:
```bash
pilti install
```
