# Delivery Check Skill

目标：交付前确认变更已经形成可审计闭环。

检查项：

- `make ci` 是否已运行并通过。
- 实质性代码或流程变更是否写入 `docs/histories/`。
- 用户可感知变化是否更新 release note。
- 流程、脚本或目录变化是否同步更新 `docs/HARNESS_PROCESS.md` 和相关入口。
- 如果修改了 `backend/go/admin/internal/services/orm/models/` 下的 model 文件，是否已在 `backend/go/admin` 运行 `make script-orm` 并提交生成的 ORM 查询代码。

没有验证证据时，不得声明交付完成。
