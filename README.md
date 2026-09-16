# multica-mcp

本项目由 [zxdwhda](https://github.com/zxdwhda) 独立维护，正式开源仓库：[zxdwhda/multica-mcp](https://github.com/zxdwhda/multica-mcp)。

面向 ChatGPT 的远程 Multica MCP 服务。以 strider2038/multica-mcp 的实现为基础独立维护，保留 MIT 许可和来源说明；本仓库不携带来源项目的 Git 历史、分支或 Fork 关系。

## 功能

395 个工具：官方当前 378 个用户 API 操作、原有 16 个便捷工具，以及 1 个 API 目录查询工具。覆盖项目、任务、评论、Agent、运行、小队、Skills、自动任务、工作区、成员、集成、插件、账号管理等，包含写入、删除和执行控制。完整清单见 [API-COVERAGE.md](docs/API-COVERAGE.md)。

- 对接官方 `https://api.multica.ai` 或自托管 Multica；支持 workspace ID/slug。
- stdio 与 Streamable HTTP；HTTP 无状态、JSON 响应，支持 FC 多实例。
- OAuth 授权页接受部署账号的有效 Multica PAT 验证身份；DCR 公共客户端、PKCE S256、授权码单次兑换、刷新令牌轮换与撤销。登录 PAT 仅用于确认同一账号，不保存或替换部署凭据。
- OAuth 数据加密存储于私有 OSS；使用 FC 角色临时凭据，代码不保存云 AK。
- 分页返回 `items`、`source_total`、`has_more`、`next_offset`。全文搜索的项目/状态/负责人筛选在单页结果中执行，必须继续翻页；空页不代表没有后续匹配。
- 任务详情的部分读取失败列入 `warnings`。批量创建失败返回已创建 ID、失败项及 `isError`，避免误报全量成功。
- 写操作可能触发 Agent；用触发预览、`suppress_run` 或 `backlog` 明确控制。批量创建是顺序操作，没有事务或自动去重。
- 新增 `multica_api_*` 工具使用 `path`、`query`、`body`、`headers`、`files` 参数；返回上游 `status` 和完整 JSON `body`，HTTP 错误标记为 `isError`。不移除上游新增字段，支持 PATCH/DELETE、multipart 上传、二进制下载和 204 空响应。

## 构建

需要 Go 1.26。

```sh
go test -race ./...
go vet ./...
go build -o bin/multica-mcp .
```

## 环境变量

| 变量 | 含义 |
| --- | --- |
| MULTICA_BASE_URL | API origin，例如 https://api.multica.ai |
| MULTICA_TOKEN | 当前 Multica PAT |
| MULTICA_WORKSPACE_ID / MULTICA_WORKSPACE_SLUG | 工作区范围 |
| MCP_TRANSPORT | stdio（默认）或 http |
| MCP_HTTP_PORT | 默认 8080 |
| MCP_HTTP_PREFIX | 默认 /multica；MCP 路径为 /multica/mcp |
| MULTICA_READ_ONLY | true 时不注册写工具 |
| MCP_API_KEY | 本地或支持静态头客户端使用；OAuth 启用时不接受此凭据 |
| MCP_OAUTH_ORIGIN | HTTPS origin；部署设为 https://mcp.wildflow.cn |
| MCP_OSS_BUCKET / MCP_OSS_ENDPOINT | OAuth 状态私有桶和 endpoint |

健康检查 `GET /multica/healthz` 仅代表进程可服务，不代表上游授权有效。

## 阿里云部署

```sh
python3 scripts/deploy_fc.py
```

脚本读取 `~/.multica/config.json` 和现有 Aliyun CLI 配置，构建 linux/amd64 静态二进制，通过 CLI 的受限临时 JSON 文件部署新加坡 FC `multica-mcp`。不修改共享 DNS、自定义域名、OSS 或其他函数。详细接入步骤见 [部署与验收](docs/DEPLOYMENT.md)。

## 边界

- 当前为固定 PAT、默认工作区的自用连接器；新增账号级和路径指定工作区的操作遵守该 PAT 在 Multica 的权限，不限于任务管理。不是多租户 Multica 登录服务。部署 PAT 轮换使 OAuth 存储解密失效，需要重新连接。
- OAuth access token 有效 1 小时，refresh token 有效 30 天；刷新令牌使用后失效，重放会撤销该授权链。
- dry_run 仅做本地预览，不验证远程 ID、权限和完整业务规则。
- 拆解工具返回固定四步模板，不调用模型生成计划。
- 0.4.0 的 API 目录固定到官方提交 `7ebe0bf58d99`（2026-09-15）。当前官网已做实际调用验证；未承诺所有历史或未来 Multica 版本支持每个接口。版本来源、覆盖范围与测试层次见 API 清单和部署文档。

来源说明见 [UPSTREAM.md](UPSTREAM.md)，原始资料在 `docs/upstream/`。
