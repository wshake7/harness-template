# 技术债追踪

这里记录那些暂时不阻塞当前任务、但已经值得留档的技术债。

| 日期 | 区域 | 债务描述 | 为什么会存在 | 计划中的后续动作 |
| --- | --- | --- | --- | --- |
| 2026-05-21 | 产品面 | 已接入管理后台和前端模板，但还没有真实产品需求文档。 | 本轮目标是代码接入后的工程文档同步，不定义产品范围。 | 在 `docs/product/specs/` 补管理后台真实用户路径、目标、非目标和验收标准。 |
| 2026-05-21 | 测试 | `make ci` 尚未运行 `pnpm ready` 或 `cd backend/go && go test ./...`。 | 项目级命令需要先确认外部依赖和本地环境稳定，避免把偶发失败直接放入基础门禁。 | 稳定前后端验证后接入 `scripts/ci.sh`。 |
| 2026-05-21 | 本地依赖 | Postgres、Redis、Temporal 需要手动准备，没有一键启动脚本。 | 新接入代码还没有配套 Docker Compose 或本地 dev stack。 | 补 `infra/` 或脚本化本地依赖启动方式，并更新 `docs/operate/reliability.md`。 |
| 2026-05-21 | 可观测性 | 健康检查、metrics、traces、dashboard 和告警入口尚未定义。 | 当前只从依赖和代码中确认了 zap、Prometheus client、Fiber monitor 等基础能力。 | 定义健康检查端点和指标访问路径，接入本地/CI 验证。 |
| 2026-05-21 | 安全 | 生产级密钥注入、Cookie/CORS/CSRF、replay 防护和审计保留策略未定稿。 | 本地默认配置可运行，但不能直接代表生产安全边界。 | 补安全设计，启用或明确替代 Nonce/replay 防护，并更新 `docs/operate/security.md`。 |
| 2026-05-21 | 交付 | release 仍产出仓库元数据制品，没有真实前端/后端构建产物。 | 部署目标和制品格式尚未确定。 | 用真实构建产物替换 `scripts/release-package.sh` 占位逻辑。 |
