#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "用法: $0 <项目名>" >&2
  exit 1
fi

project_name="$1"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

required_files=(
  "README.md"
  "AGENTS.md"
  ".service-matrix/dependencies.yaml"
  "docs/product/specs/examples/minimal-requirement-lifecycle.md"
)

missing=0

for path in "${required_files[@]}"; do
  if [[ ! -f "${repo_root}/${path}" ]]; then
    echo "初始化失败，缺少必要文件: ${path}" >&2
    missing=1
  fi
done

if [[ "${missing}" -ne 0 ]]; then
  echo "请先运行 make ci 修复模板骨架，再执行初始化。" >&2
  exit 1
fi

targets=()
for path in "${required_files[@]}"; do
  targets+=("${repo_root}/${path}")
done

perl -0pi -e "s/harness-template-cn/${project_name}/g; s/harness-template/${project_name}/g" "${targets[@]}"

echo "已将模板名称初始化为: ${project_name}"
echo "下一步建议: 补齐 docs/build/architecture.md、context/project/ 和 docs/product/specs/ 中的真实项目信息"
