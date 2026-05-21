#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmp_dir="$(mktemp -d)"
init_output="${tmp_dir}/init-project.out"
check_docs_output="${tmp_dir}/check-docs.out"

cleanup() {
  rm -rf "${tmp_dir}"
}
trap cleanup EXIT

mkdir -p "${tmp_dir}/repo"

tar \
  --exclude='./.git' \
  --exclude='./.claude/worktrees' \
  --exclude='./.worktrees' \
  -cf - \
  -C "${repo_root}" . | tar -xf - -C "${tmp_dir}/repo"

(
  cd "${tmp_dir}/repo"
  ./scripts/init-project.sh demo-harness >"${init_output}"
  grep -q 'demo-harness' README.md
  grep -q 'demo-harness' AGENTS.md
  grep -q 'project_name: demo-harness' .service-matrix/dependencies.yaml
  ./scripts/check-docs.sh >"${check_docs_output}"
)

echo "初始化脚本 smoke test 通过"
