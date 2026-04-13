#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_BIN="${GO_BIN:-/home/david/.gvm/gos/go1.26.0/bin/go}"
SCHEMA_DIR="${ROOT_DIR}/ent/v1/schema"
OUT_DIR="${ROOT_DIR}/ent/v1/gen"

if [[ ! -d "${SCHEMA_DIR}" ]]; then
  echo "[bitfs-contract] ent v1 schema dir not found: ${SCHEMA_DIR}" >&2
  exit 1
fi

mkdir -p "${OUT_DIR}"
# 每次生成前先清空产物目录，避免旧实体残留（例如已删除的 settle_*）。
rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}"
cat > "${OUT_DIR}/doc.go" <<'DOC'
// Package gen 存放由 ent v1 schema 生成的代码。
package gen
DOC

"${GO_BIN}" run entgo.io/ent/cmd/ent generate --target "${OUT_DIR}" "${SCHEMA_DIR}"

echo "[bitfs-contract] ent v1 generated into ${OUT_DIR}"
