# Installing PiltiSmart 'pilti' CLI

This guide provides instructions to install the **PiltiSmart Enterprise CLI (pilti)** on **Linux** (Ubuntu, Debian, RHEL, CentOS, Arch, Proxmox/LXC) and **macOS** (Intel & Apple Silicon M1/M2/M3/M4).

---

## ⚡ Method 1: Automatic 1-Line Global Installation (Recommended)

Run the universal installer script in your terminal:

### From Main Branch:
`ash
curl -fsSL https://raw.githubusercontent.com/PiltiSmart/stack-catalog/main/install.sh | bash
`

### From Feature Branch (eat/add-enterprise-tools):
`ash
curl -fsSL https://raw.githubusercontent.com/PiltiSmart/stack-catalog/feat/add-enterprise-tools/install.sh | bash
`

### What This Script Does:
1. Detects your Operating System (**Linux** or **macOS / Darwin**).
2. Detects your CPU architecture (**x86_64 / md64** or **rm64**).
3. Downloads the official pre-compiled static pilti binary from GitHub Releases.
4. Places the executable into /usr/local/bin/pilti (and creates symlink /usr/local/bin/ps).
5. Installs MinIO Client (mc) and configures myminio alias to http://145.241.237.108:9000.
6. Verifies installation by running pilti version.

---

## 📦 Method 2: Direct Binary Download

| Operating System | Architecture | Binary Download Link |
|---|---|---|
| **Linux** | x86_64 / md64 | https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-linux-amd64 |
| **Linux** | ARM64 / arch64 | https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-linux-arm64 |
| **macOS (Apple Silicon)** | M1 / M2 / M3 / M4 (rm64) | https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-darwin-arm64 |
| **macOS (Intel)** | x86_64 | https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-darwin-amd64 |
| **Windows** | x86_64 | https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-windows-amd64.exe |

### Manual Installation Commands:
`ash
# On Linux (amd64):
sudo curl -sSL https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-linux-amd64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti

# On Linux (ARM64):
sudo curl -sSL https://github.com/PiltiSmart/stack-catalog/releases/latest/download/pilti-linux-arm64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti
`

---

## 🔍 Verifying Installation

`ash
pilti version
pilti doctor
pilti list
`

---

## 🪣 S3 Cloud Storage CLI Setup & Verification

The installer automatically downloads the official **MinIO Client (mc)** and configures the myminio alias to point to http://145.241.237.108:9000:

`ash
# Verify S3 connection and list buckets
pilti s3 ls

# Show S3 commands overview
pilti s3 --help

# Create a bucket and upload a test file
pilti s3 mb s3://mybucket
pilti s3 cp ./test.txt s3://mybucket/
`
