# /harness:health

用途：检查当前仓库是否满足 Harness 模板的基础可用性。

执行步骤：

1. 读取 `AGENTS.md`、`docs/HARNESS_PROCESS.md` 和 `context/harness-framework/INDEX.md`。
2. 运行 `make ci`。
3. 检查 `.service-matrix/dependencies.yaml` 是否能通过 `scripts/validate-service-matrix.sh`。
4. 报告缺失的上下文、门禁或文档同步问题。

输出要求：

- 先列阻塞项。
- 再列建议项。
- 不要把未验证的状态描述为已通过。
