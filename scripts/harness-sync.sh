#!/usr/bin/env bash

set -euo pipefail

mode="${1:---staged}"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ "${HARNESS_SYNC_SKIP:-}" == "1" ]]; then
  echo "Harness sync 检查已跳过: HARNESS_SYNC_SKIP=1"
  exit 0
fi

cd "${repo_root}"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Harness sync 检查跳过: 当前目录不是 git 仓库"
  exit 0
fi

changed_files=()

case "${mode}" in
  --staged)
    while IFS= read -r file; do
      [[ -n "${file}" ]] && changed_files+=("${file}")
    done < <(git diff --cached --name-only --diff-filter=ACMRT)
    ;;
  --workspace)
    while IFS= read -r file; do
      [[ -n "${file}" ]] && changed_files+=("${file}")
    done < <(
      {
        git diff --name-only --diff-filter=ACMRT
        git diff --cached --name-only --diff-filter=ACMRT
        git ls-files --others --exclude-standard
      } | sort -u
    )
    ;;
  *)
    echo "用法: $0 [--staged|--workspace]" >&2
    exit 1
    ;;
esac

if [[ "${#changed_files[@]}" -eq 0 ]]; then
  echo "Harness sync 检查通过: 没有需要检查的变更"
  exit 0
fi

has_changed_matching() {
  local pattern="$1"
  local file

  for file in "${changed_files[@]}"; do
    if [[ "${file}" =~ ${pattern} ]]; then
      return 0
    fi
  done

  return 1
}

add_failure() {
  failures+=("$1")
}

failures=()

require_companion() {
  local reason="$1"
  local companion_pattern="$2"
  local suggestion="$3"

  if ! has_changed_matching "${companion_pattern}"; then
    add_failure "${reason}: ${suggestion}"
  fi
}

require_history() {
  require_companion \
    "交付记录缺失" \
    '^docs/records/histories/[0-9]{4}-[0-9]{2}/[0-9]{8}-[0-9]{4}-[^/]+\.md$' \
    "运行 make new-history SLUG=<change-name>，补充并 stage 对应 history。"
}

process_change_pattern='^(scripts/|\.github/|\.vite-hooks/|Makefile$|package\.json$|pnpm-workspace\.yaml$|vite\.config\.ts$|eslint\.config\.js$)'
topology_change_pattern='^(\.service-matrix/|apps/|packages/|infra/)'
delivery_change_pattern='^(scripts/|apps/|packages/|infra/|\.github/|\.vite-hooks/|\.service-matrix/|Makefile$|package\.json$|pnpm-workspace\.yaml$|vite\.config\.ts$|eslint\.config\.js$)'

if has_changed_matching "${process_change_pattern}"; then
  require_companion \
    "流程或工具链变更未同步说明" \
    '^(docs/govern/harness-process\.md|docs/operate/cicd\.md|docs/govern/collaboration\.md|docs/operate/security\.md|context/harness-framework/INDEX\.md|\.harness/)' \
    "同步 stage docs/govern/harness-process.md、docs/operate/cicd.md、docs/govern/collaboration.md、docs/operate/security.md、context/harness-framework/INDEX.md 或 .harness/ 中的相关入口。"
fi

if has_changed_matching "${topology_change_pattern}"; then
  require_companion \
    "架构或服务拓扑变更未同步说明" \
    '^(docs/build/architecture\.md|\.service-matrix/dependencies\.yaml|context/project/)' \
    "同步 stage docs/build/architecture.md、.service-matrix/dependencies.yaml 或 context/project/ 中的项目上下文。"
fi

if has_changed_matching "${delivery_change_pattern}"; then
  require_history
fi

if has_changed_matching '^\.service-matrix/dependencies\.yaml$'; then
  "${repo_root}/scripts/validate-service-matrix.sh" >/dev/null
fi

if [[ "${#failures[@]}" -ne 0 ]]; then
  echo "Harness sync 检查失败。请让文档和上下文跟上 staged diff:"
  printf '变更文件:\n'
  printf '  - %s\n' "${changed_files[@]}"
  printf '需要处理:\n'
  printf '  - %s\n' "${failures[@]}"
  echo "如本次确实无需同步，可临时使用 HARNESS_SYNC_SKIP=1 git commit ...，但不建议长期跳过。"
  exit 1
fi

echo "Harness sync 检查通过"
