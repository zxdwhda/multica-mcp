---
name: multica-local
description: 在本机通过 CLI 操作野生流动的自建 Multica，查询项目、任务、评论和 Agent。用于本地 Codex 或 Claude 的 Multica 日常交互及飞书资料关联。
---

# 本地 Multica

本机日常使用现有 `multica` CLI。业务操作方式读取已安装的 `multica-cli` Skill；执行前为其中的每条命令添加 `--profile wildflow-sg`。该 profile 将配置和运行时状态隔离，避免误用历史官方云端默认账号。

当前指定目标（2026-10-02 核验）：

- 服务：`https://multica.wildflow.cc`
- 工作区：野生流动；从当前 profile 读取 ID，与 `workspace list` 返回值核对。
- Profile：`wildflow-sg`

首次使用时检查当前配置中的服务地址和工作区 ID，再实际只读验证：

```bash
multica --profile wildflow-sg workspace list --output json
multica --profile wildflow-sg project list --output json
multica --profile wildflow-sg issue list --limit 20 --output json
```

配置可能包含 Token；只读取必要字段，不将配置或凭据输出到聊天、仓库或任务评论。目标不匹配时先说明差异；用户点名其他账号或工作区时以该请求为准，使用对应的已有 profile，不自动改默认配置。

先查任务基本信息；讨论采用 `issue comment list <id> --roots-only --summary --compact --output json`，再展开相关线程 `--thread <comment-id> --tail 30`。记录已有工作的状态或指派时使用 `--no-start`；是否派发新工作由当前用户授权决定。继承当前会话已有授权，不为同一个已授权动作重复提问。

结合飞书时：文档、会议和日历提供来源，Multica 管理执行状态，PR/CI 由所属代码托管平台核验。保留实际链接；跨工具写入、发消息或派发仍按用户本次请求完成，不自动双向同步。

ChatGPT 网页使用原 `multica-mcp` 插件及其远程 Skills。插件包不包含本地 profile、凭据或本 Skill；本地优先 CLI，网页按实际 MCP schema 操作。
