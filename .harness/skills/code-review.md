# Code Review Skill

目标：从 Harness 角度审查代码变更，而不仅是语法或风格。

优先级：

1. 行为是否满足需求和验收标准。
2. 是否违反架构、服务矩阵或安全默认约束。
3. 是否缺少测试、文档、history 或 release note。
4. 是否引入难以维护的隐式耦合。
5. 如果改了 `backend/go/admin/internal/services/orm/models/`，是否运行 `backend/go/admin` 下的 `make script-orm` 并提交生成的 ORM 查询代码。

输出时先列风险和缺陷，再列摘要。
