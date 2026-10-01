---
name: multica-tasks
description: 在 Multica 搜索、读取、创建和更新任务，查看评论线程、子任务、状态、标签、附件和关联 PR；记录进度或回复指定讨论时使用。
---

# Multica 任务与讨论

先读 [执行约定](references/execution.md)。这里的任务是 issue 工作项；`/api/tasks` 是 Agent 执行记录，两者不能混用。

## 查找与读取

- 用 `multica_list_projects` 解析项目名，保存真实 ID；同名项目按当前工作区和用户上下文消歧。
- 用 `multica_search_tasks` 搜索标题、正文和评论，或 `multica_list_tasks` 按项目、状态、负责人筛选。沿 `next_offset` 翻页，直到 `has_more=false`；搜索后的本地筛选可能产生空页但仍有下一页。
- 只读详情用 `multica_api_get_issue`，参数 `path:{id:真实ID或任务编号}`，避免自动加载评论与子任务。便捷工具 `multica_get_task` 的 schema 支持时也可传 `include_comments:false, include_subtasks:false`；不要给旧 schema 传入未声明字段。跳过的部分不代表不存在。
- 按 [评论读取流程](references/comments.md) 先摘要、后相关线程，避免默认拉取整段历史。
- 自定义状态先查询 `multica_list_statuses`。标签、属性、附件、linked PR、历史等用 `multica_api_catalog` 查询对应 issue 路由及实际 schema，不把便捷工具不支持误报成 Multica 不支持。

## 创建、更新与记录

- 从当前用户请求确定目标、标题、内容、负责人和需要的执行行为。普通可逆选择自行完成；缺失会实质影响目标或授权时再问。
- 创建时默认 todo；指定 Agent 或 squad 可能立即运行。只登记后续工作可用 backlog。建依赖树时先以 backlog 建完，再按用户要求启动可执行阶段，避免树未完成就开始执行。
- 仅整理状态、日期或负责人、不需要新一轮执行时，更新使用 `suppress_run:true`。需要开始执行时先结合当前状态和 `multica_preview_issue_triggers` 检查实际影响。
- 只留记录用 `/note` 前缀，避免评论触发 Agent。要继续 Agent 工作时可先 `multica_preview_comment_triggers`；提及 Agent/squad 能触发运行，提及成员可能发通知。
- 回复指定问题时，把 `parent_id` 指向那条真实评论；不要随意挂到最新一条。
- `multica_create_task_with_subtasks` 顺序执行，可能部分成功。保留返回 ID，检查 `complete` 与 `failures`，只补失败项，不能整批重试。
- `multica_plan_task_breakdown` 只返回固定四步模板。针对业务的拆解应结合上下文自行推理，不声称此工具理解了项目。
- `dry_run` 只检查本地输入，不证明远端 ID、权限、状态规则或创建成功。写入后读取目标并报告实际改变。

关联 PR 是任务线索；代码评审、CI、合并和发布结果须从所属 GitHub/Forgejo 核实。任务标记完成不能替代业务验收。
