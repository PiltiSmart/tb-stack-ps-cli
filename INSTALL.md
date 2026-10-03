# Installing PiltiSmart 'pilti' CLI

This guide provides instructions to install the **PiltiSmart Enterprise CLI (`pilti`)** on **Linux**, **macOS**, and **Windows**.

---

## ⚡ Method 1: Automatic 1-Line Installation (Recommended)

### 🐧 Linux & 🍎 macOS (Bash / Zsh):
```bash
curl -fsSL https://raw.githubusercontent.com/PiltiSmart/tb-stack-ps-cli/main/install.sh | bash
```

### 🪟 Windows (PowerShell):
```powershell
irm https://raw.githubusercontent.com/PiltiSmart/tb-stack-ps-cli/main/install.ps1 | iex
```

### What This Script Does Automatically:
1. Detects your Operating System (**Linux**, **macOS / Darwin**, or **Windows**).
2. Detects your CPU architecture (**amd64 / x86_64** or **arm64**).
3. Downloads the official pre-compiled static `pilti` binary from GitHub Releases.
4. Places the executable into `/usr/local/bin/pilti` (or `%LOCALAPPDATA%\Programs\pilti\pilti.exe` on Windows).
5. Installs the MinIO Client (`mc` / `mc.exe`) for full AWS S3-compatible cloud storage operations.
6. Verifies installation by running `pilti version`.

---

## 📦 Method 2: Direct Binary Download

| Operating System | Architecture | Binary Download Link |
|---|---|---|
| **Linux** | x86_64 / amd64 | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-linux-amd64` |
| **Linux** | ARM64 / aarch64 | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-linux-arm64` |
| **macOS (Apple Silicon)** | M1 / M2 / M3 / M4 (arm64) | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-darwin-arm64` |
| **macOS (Intel)** | x86_64 | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-darwin-amd64` |
| **Windows** | x86_64 | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-windows-amd64.exe` |

---

## 🔍 Verifying Installation

```bash
pilti version
pilti doctor
pilti list
```

---

## 🪣 S3 Cloud Storage CLI Setup & Verification

```bash
# Verify S3 connection and list buckets
pilti s3 ls

# Show S3 commands overview
pilti s3 --help

# Create a bucket and upload a file
pilti s3 mb s3://mybucket
pilti s3 cp ./test.txt s3://mybucket/
```
