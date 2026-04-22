#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
go_bin="${GO_BIN:-/home/david/.gvm/gos/go1.26.0/bin/go}"

cd "$repo_root"

echo "[bitfs-contract] gate: proto generated check"
"$repo_root/scripts/check_proto_generated.sh"

echo "[bitfs-contract] gate: go test"
"$go_bin" test ./...

echo "[bitfs-contract] gate: ok"
