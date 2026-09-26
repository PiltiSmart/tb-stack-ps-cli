# Installing PiltiSmart 'pilti' CLI on Linux & macOS

This guide provides instructions to install the **PiltiSmart Enterprise CLI (`pilti`)** globally on **Linux** (Ubuntu, Debian, RHEL, CentOS, Arch, Proxmox/LXC) and **macOS** (Intel & Apple Silicon M1/M2/M3/M4).

---

## ⚡ Method 1: Automatic 1-Line Global Installation (Recommended)

Run the universal installer script in your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/PiltiSmart/tb-stack-ps-cli/main/install.sh | bash
```

### What This Script Does:
1. Automatically detects your Operating System (**Linux** or **macOS / Darwin**).
2. Detects your CPU architecture (**`x86_64` / `amd64`** or **`arm64` / Apple Silicon**).
3. Downloads the matching pre-compiled static `pilti` binary.
4. Places the executable into `/usr/local/bin/pilti` and sets executable permissions (`chmod +x`).
5. Verifies installation by running `pilti version`.

---

## 📦 Method 2: Manual Installation via Pre-compiled Binaries

If you prefer downloading the binary manually:

### 1. Choose Your Architecture:

| Operating System | Architecture | Binary Download Link |
|---|---|---|
| **Linux** | `x86_64` / `amd64` | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-linux-amd64` |
| **Linux** | `ARM64` / `aarch64` | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-linux-arm64` |
| **macOS (Apple Silicon)** | `M1 / M2 / M3 / M4` (`arm64`) | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-darwin-arm64` |
| **macOS (Intel)** | `x86_64` | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-darwin-amd64` |
| **Windows** | `x86_64` | `https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-windows-amd64.exe` |

### 2. Download and Move to `/usr/local/bin`:

#### On Linux (amd64 / x86_64):
```bash
sudo curl -fsSL https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-linux-amd64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti
```

#### On macOS (Apple Silicon M1/M2/M3/M4):
```bash
sudo curl -fsSL https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-darwin-arm64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti
```

#### On macOS (Intel):
```bash
sudo curl -fsSL https://github.com/PiltiSmart/tb-stack-ps-cli/releases/latest/download/pilti-darwin-amd64 -o /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti
```

---

## 🔨 Method 3: Build from Source (Requires Go 1.21+)

If you have Go installed on your machine:

```bash
# Clone the repository
git clone https://github.com/PiltiSmart/tb-stack-ps-cli.git
cd tb-stack-ps-cli

# Compile stripped production binary
go mod tidy
go build -ldflags="-s -w" -o bin/pilti main.go

# Install globally
sudo cp bin/pilti /usr/local/bin/pilti
sudo chmod +x /usr/local/bin/pilti
```

---

## 🔍 Verifying the Installation

Check that `pilti` is recognized in your terminal:

```bash
# Verify version
pilti version

# List all available software in PiltiSmart
pilti list

# Run pre-flight health diagnostics
pilti doctor

# Check individual software status
pilti tb-app status

# Install a software component interactively
pilti install
```

---

## 💡 Troubleshooting & Notes

### No Linux Command Collisions
By naming the binary `pilti`, it cleanly avoids any collision with the standard Linux process reporting tool `/bin/ps`.

### macOS Security Prompt (Gatekeeper)
If macOS blocks the binary on first execution because it was downloaded via the web, clear the quarantine attribute:
```bash
xattr -d com.apple.quarantine /usr/local/bin/pilti
```

---

## 🗑️ Uninstallation

To remove `pilti` from your system at any time:
```bash
sudo rm -f /usr/local/bin/pilti
```
