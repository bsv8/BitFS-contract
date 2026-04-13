#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

"${ROOT_DIR}/scripts/gen_proto.sh"

if [[ -n "$(git -C "${ROOT_DIR}" status --porcelain -- gen/go 2>/dev/null || true)" ]]; then
  echo "[bitfs-contract] generated proto files are not up to date" >&2
  echo "[bitfs-contract] run: BitFS-contract/scripts/gen_proto.sh" >&2
  exit 1
fi

echo "[bitfs-contract] proto generated files are up to date"
