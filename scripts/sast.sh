#!/usr/bin/env bash

set -Eeuo pipefail

# Locate the project root and backend directory.
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"

# Add Go-installed tools to PATH.
GOPATH_BIN="$(go env GOPATH)/bin"
export PATH="$PATH:$GOPATH_BIN"

echo "======================================"
echo " SoilCycle - Static Security Analysis"
echo "======================================"

# Verify backend directory.
if [[ ! -d "$BACKEND_DIR" ]]; then
    echo "[ERROR] Backend directory not found: $BACKEND_DIR"
    exit 1
fi

# Verify Gosec installation.
if ! command -v gosec >/dev/null 2>&1; then
    echo "[ERROR] Gosec is not installed or not available in PATH."
    echo "Install it using:"
    echo "go install github.com/securego/gosec/v2/cmd/gosec@latest"
    exit 1
fi

cd "$BACKEND_DIR"

echo
echo "[INFO] Scanning backend Go source code..."
echo "[INFO] Scanner: $(gosec -version 2>&1 | head -n 1)"
echo

if gosec -fmt=text ./...; then
    echo
    echo "[PASS] Gosec scan completed successfully."
else
    status=$?
    echo
    echo "[FAIL] Gosec detected issues or could not complete."
    exit "$status"
fi

echo
echo "======================================"
echo " SAST scan completed"
echo "======================================"
