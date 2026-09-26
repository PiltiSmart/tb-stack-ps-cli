#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "======================================================="
echo "   BUILDING PILTISMART 'pilti' CLI (GO + COBRA)"
echo "======================================================="

mkdir -p bin
echo "[1/2] Downloading Go dependencies..."
go mod tidy

echo "[2/2] Compiling standalone static binary..."
go build -ldflags="-s -w" -o bin/pilti main.go
chmod +x bin/pilti

echo "[SUCCESS] Compiled binary available at: bin/pilti"
echo ""
echo "To install globally, run:"
echo "   cp bin/pilti /usr/local/bin/pilti"
