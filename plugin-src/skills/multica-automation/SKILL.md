---
name: multica-automation
description: 使用 Multica Skill 与文件库、导入、autopilot、触发器、插件和其他完整 API 功能；便捷工具不足时按实际 API schema 完成操作。
---

# Multica 能力库与自动化

先读 [执行约定](references/execution.md)。这里的 Multica Skill 库是 Agent 运行时使用的业务能力；本插件里的 Skill 是当前助手的操作指南，二者不能混为已安装状态。

- 用 `multica_api_catalog` 按 `/api/skills`、`/api/autopilots`、`/api/runtimes`、`/api/workspaces` 或具体功能查询，读取匹配的 HTTP 方法、路径参数、输入 schema 和来源版本后调用。
- 操作 Multica Skill：先读取 Skill 与文件，确认内容和归属；用户要求时创建、更新、导入或分配给指定 Agent。导入请求可能异步，沿 request ID 查询结果；接受请求不是导入完成。
- 附件与文件上传使用工具实际 `files` schema，文件内容以 base64 传入，每文件最多 4 MiB、请求最多 8 MiB。不把本机路径传给远端当文件内容。下载返回 `data_base64` 时按宿主支持保存真实文件，不伪造下载链接。
- autopilot：先读配置、时区、触发条件、执行动作、近期 delivery/run；定时表达式可先查预览。用户授权后再创建/更改/启停，回读配置，触发后跟踪 run/delivery，不以 HTTP 2xx 代替业务执行。
- webhook、插件安装、成员与计费操作要核对目标与用户要求；密钥不回显，不写进文档或插件包。诊断不默认轮换密钥、扩权限或更换账号。
- 处理 PATCH/PUT 时按写入 schema 组织数据，不能原样复制详情响应。保留有意的 null、空字符串、false、0。有 ETag/revision 时遵循 API 的冲突处理规则，冲突后重新读取目标。
- 对未知功能先查当前 catalog，不猜工具名或任意 URL。完整 API 仍由 Multica 的账号权限和功能可用性约束。

没有明确要求定期执行时，不创建 autopilot；没有明确要求发送时，不把成功结果自动推送到飞书或其他聊天。
