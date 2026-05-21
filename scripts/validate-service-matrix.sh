#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
matrix="${repo_root}/.service-matrix/dependencies.yaml"

if [[ ! -f "${matrix}" ]]; then
  echo "缺少服务矩阵: .service-matrix/dependencies.yaml" >&2
  exit 1
fi

required_patterns=(
  '^project_name:'
  '^placeholders:'
  '^roots:'
  '^modules:'
  '^services:'
  'repo_path:'
)

failed=0

for pattern in "${required_patterns[@]}"; do
  if ! grep -Eq "${pattern}" "${matrix}"; then
    echo "服务矩阵缺少必要字段或段落: ${pattern}" >&2
    failed=1
  fi
done

allowed_placeholders=(
  app-root
  package-root
  infra-root
  project-name
)

while IFS= read -r placeholder; do
  name="${placeholder#\{}"
  name="${name%\}}"
  found=0
  for allowed in "${allowed_placeholders[@]}"; do
    if [[ "${name}" == "${allowed}" ]]; then
      found=1
      break
    fi
  done
  if [[ "${found}" -eq 0 ]]; then
    echo "服务矩阵包含未登记占位符: ${placeholder}" >&2
    failed=1
  fi
done < <(grep -Eo '\{[^}]+\}' "${matrix}" | sort -u)

if ! grep -Eq '^[[:space:]]{2}[a-zA-Z0-9_-]+:' "${matrix}"; then
  echo "服务矩阵至少需要一个模块或服务条目" >&2
  failed=1
fi

if [[ "${failed}" -ne 0 ]]; then
  exit 1
fi

echo "服务矩阵检查通过"
