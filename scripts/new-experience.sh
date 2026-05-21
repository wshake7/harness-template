#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "用法: $0 <context-scope> <经验名>" >&2
  echo "示例: $0 project/example pagination-limit" >&2
  exit 1
fi

scope="$1"
slug="$2"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

case "${scope}" in
  team|harness-framework|project/*)
    ;;
  *)
    echo "context-scope 必须是 team、harness-framework 或 project/<name>" >&2
    exit 1
    ;;
esac

if [[ ! "${slug}" =~ ^[a-z0-9][a-z0-9-]*$ ]]; then
  echo "经验名只能使用小写字母、数字和连字符，并且必须以字母或数字开头" >&2
  exit 1
fi

target_dir="${repo_root}/context/${scope}/experience"
mkdir -p "${target_dir}"

timestamp="$(date +%Y%m%d-%H%M)"
target="${target_dir}/${timestamp}-${slug}.md"

cat > "${target}" <<EOF
# ${slug}

## 触发场景

记录这条经验来自哪次纠正、事故、评审或重复问题。

## 可复用结论

写清下次同类任务应如何判断和行动。

## 适用范围

说明这是团队级、框架级还是具体项目级经验。

## 验证方式

说明这条经验如何被脚本、门禁、测试或评审检查复用。
EOF

echo "已创建经验记录: ${target#${repo_root}/}"
