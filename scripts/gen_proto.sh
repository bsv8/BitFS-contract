#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROTO_DIR="${ROOT_DIR}/proto"
OUT_DIR="${ROOT_DIR}/gen/go"

if ! command -v protoc >/dev/null 2>&1; then
  echo "[bitfs-contract] protoc not found; please install protoc first" >&2
  exit 1
fi

if ! command -v /home/david/.gvm/gos/go1.26.0/bin/go >/dev/null 2>&1; then
  echo "[bitfs-contract] go1.26.0 toolchain not found" >&2
  exit 1
fi

/home/david/.gvm/gos/go1.26.0/bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6

PROTOC_GEN_GO_BIN="$(/home/david/.gvm/gos/go1.26.0/bin/go env GOPATH)/bin/protoc-gen-go"
if [[ ! -x "${PROTOC_GEN_GO_BIN}" ]]; then
  echo "[bitfs-contract] protoc-gen-go not found at ${PROTOC_GEN_GO_BIN}" >&2
  exit 1
fi

mkdir -p "${OUT_DIR}"

protoc \
  --proto_path="${PROTO_DIR}" \
  --plugin=protoc-gen-go="${PROTOC_GEN_GO_BIN}" \
  --go_out="${OUT_DIR}" \
  --go_opt=paths=source_relative \
  "${PROTO_DIR}/v1/common.proto" \
  "${PROTO_DIR}/v1/direct_transfer.proto" \
  "${PROTO_DIR}/v1/live.proto"

echo "[bitfs-contract] proto generated into ${OUT_DIR}"
