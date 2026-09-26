# PiltiSmart 'pilti' Enterprise CLI (Go + Cobra Framework)

A standalone, zero-dependency CLI written in **Go (Golang)** using the **Cobra Framework** to discover, configure, inspect, and orchestrate all software components across the **PiltiSmart** ecosystem:

1. `tb-app` : ThingsBoard Core Application
2. `tb-db` : TimescaleDB / PostgreSQL database
3. `tb-edge` : ThingsBoard Edge Gateway
4. `jenkins` : Jenkins CI/CD Automation Engine
5. `piltiservices` : PiltiSmart Microservices Backend
6. `kafka` : Apache Kafka Distributed Event Streaming Broker
7. `pilticloud` : PiltiSmart Cloud Gateway & Sync Tunnel (PMX)

---

## 🚀 Key Commands

### 1. List All Available Software (`pilti list`)
Displays a formatted catalog of all PiltiSmart software components with live container statuses, port assignments, and descriptions:
```bash
pilti list
```

### 2. Software-Specific Management (`pilti <software-id> <subcommand>`)
Manage any software component directly with granular commands:
```bash
# Check status of specific software
pilti tb-app status
pilti jenkins status
pilti kafka status

# Deploy / install a software
pilti tb-app install
pilti jenkins install
pilti kafka install
pilti piltiservices install
pilti pilticloud install

# Restart, stop, or view logs
pilti tb-db restart
pilti kafka logs -f
pilti tb-app logs -f
pilti pilticloud stop
```

### 3. Interactive Installation Wizard (`pilti install`)
Run without arguments to launch an interactive terminal menu:
```bash
pilti install
```
Or specify the target directly:
```bash
pilti install tb-app
pilti install all
```

### 4. Pre-Flight Diagnostics (`pilti doctor`)
Checks host OS, Docker engine, Docker compose plugin, CPU/RAM/Disk metrics, and network port availability across all software components:
```bash
pilti doctor
```

### 5. Check Global Status (`pilti status`)
```bash
pilti status
```

---

## ⚡ Quick Global Installation (Linux & macOS)

Install the `pilti` CLI with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/PiltiSmart/tb-stack-ps-cli/main/install.sh | bash
```

> For manual installation options (Apple Silicon, Intel Mac, ARM64, Windows, or building from source), see the complete [Installation Guide](INSTALL.md).

---

## 🔨 How to Build Locally from Source

```bash
# Build binary:
./build.sh        # Linux / macOS
.\build.bat       # Windows

# Or install globally (Linux / macOS):
make install
```
