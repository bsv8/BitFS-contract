#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

"${ROOT_DIR}/scripts/gen_ent_v1.sh"

if [[ -n "$(git -C "${ROOT_DIR}" status --porcelain -- ent/v1/gen 2>/dev/null || true)" ]]; then
  echo "[bitfs-contract] generated ent v1 files are not up to date" >&2
  echo "[bitfs-contract] run: BitFS-contract/scripts/gen_ent_v1.sh" >&2
  exit 1
fi

echo "[bitfs-contract] ent v1 generated files are up to date"
