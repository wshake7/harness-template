# /requirement:gate-check

用途：按 `docs/HARNESS_PROCESS.md` 检查当前需求是否能进入下一阶段。

执行步骤：

1. 定位需求文档或用户提供的需求上下文。
2. 检查是否有目标、非目标和验收标准。
3. 如果进入实现阶段，检查是否已有 execution plan 或轻量设计说明。
4. 运行 `scripts/validate-service-matrix.sh`。
5. 如果准备交付，确认 `make ci` 和 history/release note 状态。

输出要求：

- 明确当前阶段。
- 明确通过、阻塞、建议补强三类结论。
- 阻塞项必须引用具体文件或缺失产物。
