#!/usr/bin/env bash
# QingLong runner for faw-vw-auto
# This script is intended to be used in QingLong (青龙) task scripts.
# It prefers environment variables (FAW_*) and will export FAW_USE_ENV=1 to force the app
# to read settings from environment variables instead of env.yaml.

set -euo pipefail

# If you want the program to create an env.yaml instead of reading env vars,
# set CREATE_ENV_YAML=1 in QingLong and the script will write a minimal env.yaml
# from FAW_* variables.

export FAW_USE_ENV=1
export DUMP_ENV=1

# NO LEGACY VARS: use FAW_RUN_ONCE to run once and exit when used in scheduler

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Optional debug: if DUMP_ENV=1, print masked environment values and env.yaml info to help debugging in QingLong
if [ "${DUMP_ENV:-}" = "1" ]; then
  echo "--- DUMP ENV ---"
  echo "FAW_USE_ENV=${FAW_USE_ENV:-}" 
  echo "FAW_RUN_ONCE=${FAW_RUN_ONCE:-}"
  # mask mobile and password
  if [ -n "${FAW_MOBILE:-}" ]; then
    # print only last 4 digits
    echo "FAW_MOBILE=***${FAW_MOBILE: -4}"
  else
    echo "FAW_MOBILE="
  fi
  if [ -n "${FAW_PASSWORD:-}" ]; then
    echo "FAW_PASSWORD=***" 
  else
    echo "FAW_PASSWORD="
  fi
  echo "CREATE_ENV_YAML=${CREATE_ENV_YAML:-}"
  if [ -f "${ROOT_DIR}/env.yaml" ]; then
    echo "env.yaml exists:" && sed -n '1,200p' "${ROOT_DIR}/env.yaml" | sed -n '1,20p'
  else
    echo "env.yaml not present in ${ROOT_DIR}"
  fi
  echo "--- END DUMP ---"
fi

if [ "${CREATE_ENV_YAML:-}" = "1" ]; then
  cat > "$ROOT_DIR/env.yaml" <<EOF
BarkPushServerURL: "${FAW_BARK_PUSH_SERVER_URL:-}"
mobile: "${FAW_MOBILE:-}"
password: "${FAW_PASSWORD:-}"
WebDid: "${FAW_WEBDID:-}"
securityCode: "${FAW_SECURITY_CODE:-}"
did: "${FAW_DID:-}"
ExecFreq: "${FAW_EXEC_FREQ:-}"
EOF
  echo "wrote $ROOT_DIR/env.yaml"
  # unset FAW_USE_ENV so the app reads env.yaml
  unset FAW_USE_ENV
fi

# BINARY selection: prefer binary in current dir, otherwise try 'faw-vw-auto' in PATH
BINARY_LINE_NOTE="# BINARY selection: prefer binary in current dir, otherwise try 'faw-vw-auto' in PATH"
BINARY="$ROOT_DIR/faw-vw-auto"
if [ ! -x "$BINARY" ]; then
  BINARY="$(command -v faw-vw-auto || true)"
fi

if [ -z "$BINARY" ]; then
  echo "binary not found, please build the binary and place it next to this script or ensure 'faw-vw-auto' is in PATH"
  exit 2
fi

echo "Starting faw-vw-auto using $BINARY"
exec "$BINARY"
