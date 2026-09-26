#!/usr/bin/env bash
# ==============================================================================
# PiltiSmart 'pilti' CLI Universal Installer for Linux & macOS
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/PiltiSmart/tb-stack-ps-cli/main/install.sh | bash
# ==============================================================================

set -e

REPO="PiltiSmart/tb-stack-ps-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="pilti"

# Color helpers
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${CYAN}${BOLD}"
echo "=========================================================="
echo "    PiltiSmart 'pilti' Enterprise CLI Installer"
echo "    Target: Linux & macOS (Darwin)"
echo "=========================================================="
echo -e "${NC}"

# 1. Detect Operating System
OS="$(uname -s)"
case "${OS}" in
    Linux*)     PLATFORM="linux";;
    Darwin*)    PLATFORM="darwin";;
    *)          
        echo -e "${RED}[ERROR] Unsupported Operating System: ${OS}${NC}"
        echo "The 'pilti' CLI currently supports Linux and macOS (Darwin)."
        exit 1
        ;;
esac

# 2. Detect CPU Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64|amd64)   TARGET_ARCH="amd64";;
    aarch64|arm64)  TARGET_ARCH="arm64";;
    armv7l)         TARGET_ARCH="arm";;
    *)
        echo -e "${RED}[ERROR] Unsupported architecture: ${ARCH}${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}[✔] Detected Environment:${NC} OS=${PLATFORM}, ARCH=${TARGET_ARCH}"

# 3. Elevate permissions helper
SUDO=""
if [ "$(id -u)" -ne 0 ]; then
    if command -v sudo >/dev/null 2>&1; then
        SUDO="sudo"
    else
        echo -e "${YELLOW}[!] Warning: Not running as root and sudo is not installed.${NC}"
        INSTALL_DIR="${HOME}/.local/bin"
        mkdir -p "${INSTALL_DIR}"
        echo -e "${CYAN}[i] Falling back to installation path: ${INSTALL_DIR}${NC}"
    fi
fi

# 4. Resolve latest release or compile locally if Go is available
TMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

RELEASE_BINARY="pilti-${PLATFORM}-${TARGET_ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${RELEASE_BINARY}"

echo -e "${CYAN}[1/3] Fetching '${BINARY_NAME}' for ${PLATFORM}/${TARGET_ARCH}...${NC}"

DOWNLOAD_SUCCESS=0
if curl -sLf "${DOWNLOAD_URL}" -o "${TMP_DIR}/${BINARY_NAME}"; then
    DOWNLOAD_SUCCESS=1
else
    echo -e "${YELLOW}[!] Pre-compiled binary not yet found on GitHub Releases.${NC}"
    # Check if we are inside the source directory or Go is available to build
    if command -v go >/dev/null 2>&1; then
        echo -e "${CYAN}[i] Detected Go toolchain. Compiling from source...${NC}"
        if [ -f "main.go" ]; then
            go mod tidy
            go build -ldflags="-s -w" -o "${TMP_DIR}/${BINARY_NAME}" main.go
            DOWNLOAD_SUCCESS=1
        else
            echo -e "${CYAN}[i] Fetching repository to compile...${NC}"
            git clone --depth 1 "https://github.com/${REPO}.git" "${TMP_DIR}/source"
            (cd "${TMP_DIR}/source" && go mod tidy && go build -ldflags="-s -w" -o "${TMP_DIR}/${BINARY_NAME}" main.go)
            DOWNLOAD_SUCCESS=1
        fi
    fi
fi

if [ ${DOWNLOAD_SUCCESS} -ne 1 ]; then
    echo -e "${RED}[ERROR] Failed to download or build 'pilti' binary.${NC}"
    echo "Please ensure internet connectivity or install Go: https://go.dev/dl/"
    exit 1
fi

# 5. Install to system path
echo -e "${CYAN}[2/3] Installing binary into ${INSTALL_DIR}/${BINARY_NAME}...${NC}"
chmod +x "${TMP_DIR}/${BINARY_NAME}"
${SUDO} mkdir -p "${INSTALL_DIR}"
${SUDO} cp "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
${SUDO} chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

# 6. Install MinIO Client ('mc') and configure S3 alias
echo -e "${CYAN}[3/4] Installing MinIO Client ('mc') for S3 operations...${NC}"
MC_URL="https://dl.min.io/client/mc/release/${PLATFORM}-${TARGET_ARCH}/mc"
if curl -sLf "${MC_URL}" -o "${TMP_DIR}/mc"; then
    chmod +x "${TMP_DIR}/mc"
    ${SUDO} cp "${TMP_DIR}/mc" "${INSTALL_DIR}/mc"
    ${SUDO} chmod +x "${INSTALL_DIR}/mc"
    echo -e "${GREEN}[✔] Installed 'mc' to ${INSTALL_DIR}/mc${NC}"

    # Configure myminio alias
    echo -e "${CYAN}[i] Configuring MinIO alias 'myminio' (http://145.241.237.108:9000)...${NC}"
    "${INSTALL_DIR}/mc" alias set myminio http://145.241.237.108:9000 minioadmin minioadmin123 >/dev/null 2>&1 || true
    if [ -n "${SUDO}" ]; then
        ${SUDO} "${INSTALL_DIR}/mc" alias set myminio http://145.241.237.108:9000 minioadmin minioadmin123 >/dev/null 2>&1 || true
    fi
    echo -e "${GREEN}[✔] S3 / MinIO alias 'myminio' configured successfully!${NC}"
else
    echo -e "${YELLOW}[!] Note: Could not pre-fetch 'mc' from ${MC_URL}. 'pilti s3' will auto-install it on first run.${NC}"
fi

# 7. Verify installation
echo -e "${CYAN}[4/4] Verifying installation...${NC}"
if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
    echo -e "${GREEN}${BOLD}✔ Successfully installed '${BINARY_NAME}' to ${INSTALL_DIR}/${BINARY_NAME}!${NC}"
    echo ""
    "${BINARY_NAME}" version
    echo ""
    echo -e "Run ${BOLD}'pilti list'${NC} to view all software components."
    echo -e "Run ${BOLD}'pilti s3 ls'${NC} to view S3 buckets."
    echo -e "Run ${BOLD}'pilti s3 --help'${NC} to explore AWS S3-compatible commands."
    echo -e "Run ${BOLD}'pilti doctor'${NC} to check system health."
else
    echo -e "${YELLOW}[!] '${BINARY_NAME}' installed to ${INSTALL_DIR}, but ${INSTALL_DIR} is not in your current PATH.${NC}"
    echo -e "Add it to your shell configuration (.bashrc, .zshrc):"
    echo -e "    export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
