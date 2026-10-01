# 按需要读取任务评论

用 `multica_api_get_issue` 的 `path:{id:真实ID或任务编号}` 获取 issue 的真实 ID；便捷工具 schema 支持时可用 `multica_get_task` 并传 `include_comments:false, include_subtasks:false`。评论使用 `multica_api_list_comments`，参数结构为 `path:{id:ISSUE_ID}` 和 `query:{...}`。

## 两步读取

1. 先用 `query:{roots_only:true, summary:true}` 获取根线程的短摘要、`reply_count` 和 `last_activity_at`。
2. 根据用户问题选线程，使用 `query:{thread:ROOT_ID, tail:10, summary:true}` 阅读根评论和最新若干回复。正文不足时再次读取同一范围、设置 `summary:false`，不要默认展开所有无关线程。

需要最近几个活跃线程时用 `query:{recent:3, summary:true}`。`recent` 数的是线程，包含它们的回复，返回条数不等于 3；一个长线程仍可能很多，优先摘要再选线程 `tail`。

## 分页与完整性

- `X-Multica-Next-Before` 与 `X-Multica-Next-Before-Id` 是复合游标；下一页同时传 `before`、`before_id`，保持原 recent 或 thread+tail 范围。不能只传其中一个。
- `X-Comments-Truncated:true` 表示服务端保护上限导致结果不全。`content_truncated:true` 表示该条摘要省略正文；这两者不能混用。
- `roots_only` 不能与 thread、recent、tail 或游标混用。thread 与 recent 互斥；tail 仅用于 thread，允许 0（只取根），游标只能用于 recent 或 thread+tail。不要给该路由加不支持的 offset/limit。
- `since` 使用 RFC3339 时间；是增量评论读取，不等于完整线程。
- 原生 HTTP API 默认 `fold:false`；完整阅读已解决线程时可显式 `fold:true`，检查 `thread_resolved` 和 `folded_count`。fold 不能与 roots_only、tail 或 since 混用。CLI 的默认折叠行为不能直接套到 HTTP 默认。
- 无下一页游标也可能有截断标记。遇到截断且不能继续分页，明确标注资料不全；缩小为相关线程/时间范围或提供任务链接让用户查看，不声称读完全部。
- 旧 MCP 服务可能未转发这些上游响应头；未出现标记不能证明完整。旧连接优先采用摘要与单线程有界读取，报告范围；需要完整历史时说明连接的完整性限制，不把缺少游标当成已读完。

回复时引用对应评论 ID，`parent_id` 指向真实待答评论。记录而不触发 Agent 用 `/note`；希望执行则结合 trigger preview 与用户授权。

来源：官方 multica-cli Skill `f391633862eed012760bb905d8c6d6df17f0e484`；Multica 服务端 `server/internal/handler/comment.go`，版本 `c7f259c70`（本机 CLI 0.4.44 对应提交）。远端部署版本若不支持某参数，报告实际错误并按其 schema 收窄，不能悄悄改为无限读取。
