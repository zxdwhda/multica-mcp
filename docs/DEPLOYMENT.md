# 部署与验收

## 路由契约

地域：ap-southeast-1（新加坡）；函数：multica-mcp；运行时：custom.debian11；端口：9000；当前别名/版本：LATEST。

| 公网路径 | 方法 | 用途 |
| --- | --- | --- |
| /multica/mcp | POST、GET、DELETE | Streamable HTTP MCP（无状态模式不提供常驻 GET SSE） |
| /multica/healthz | GET | 进程健康 |
| /multica/authorize | GET、POST | PAT 登录与同意授权 |
| /multica/register | POST | OAuth 公共客户端动态注册 |
| /multica/token | POST | 授权码兑换和刷新 |
| /multica/revoke | POST | 撤销授权链 |
| /.well-known/oauth-authorization-server/multica | GET | issuer 元数据 |
| /.well-known/oauth-protected-resource/multica/mcp | GET | 资源元数据 |

issuer 为 `https://mcp.wildflow.cn/multica`，resource 为 `https://mcp.wildflow.cn/multica/mcp`。域名路由直接转发，不改写路径。`/multica/.well-known/` 同时提供兼容元数据。

## ChatGPT 接入

在自定义 MCP 插件中填入 `https://mcp.wildflow.cn/multica/mcp`，选择 OAuth。服务提供动态客户端注册，客户端 ID/secret 不必预填。授权页面填写部署账号的任意有效 PAT，确认页面显示的返回地址属于当前 ChatGPT 客户端，然后授权连接。服务通过 `/api/me` 比对账号身份；登录 PAT 不保存、不替换部署凭据，也不会放入回调 URL 或交给 ChatGPT。连接成功后如工具列表为空，点击插件详情中的“刷新”。

可通过本机配置文件取 PAT；不要把 PAT 粘贴到聊天、Issue、截图或 Git。具体 ChatGPT UI 可能变更，以实际授权流程为准。

## 状态与运维

- 私有 OSS：wildflow-mcp-state-sg，独立 `multica/` 前缀。
- 复用 AliyunFcDefaultRole，运行时只读取 FC 注入的临时凭据环境变量；不信任请求方提供的云凭据头。
- AES-GCM 加密，派生密钥绑定当前 Multica PAT；对象键作为附加认证数据。OSS 对象名中仅有随机标识的哈希。
- 单次使用依赖 `x-oss-forbid-overwrite`；桶不能启用版本控制。不要删除仍有效授权的 claim 对象。
- 当前无全桶生命周期规则，避免清理永久 DCR 注册记录。后续清理可区分 clients 与过期 pending/code/token，且 claim 存活必须覆盖对应授权有效期。
- 每次云端部署代码和状态与本地 commit 是不同的验收层；`deploy/build/deployment.json` 记录最近 CLI 上传结果，不代表 ChatGPT 验收通过。
- 默认 FC 测试域同样受应用 OAuth 保护。共享域名与证书由统一域名任务维护。
- 2026-09-15 实际 TLS 握手核对：共享域名证书到期时间为 2026-12-14 06:03:53 UTC（北京时间 14:03:53）。共享域名维护任务确认目前通过手动续期及域名部署入口更新，尚未配置定时续期；需要在到期前由域名维护方处理。Multica 函数代码部署不会续期证书。

## 测试

`go test -race ./...` 覆盖鉴权、PKCE、授权码并发重放、refresh 重放撤销、同意页 CSRF、分页和 HTTP 无状态请求；部署后还需完成真实 OAuth 和 MCP 调用以及 ChatGPT 页面接入。

浏览器回归另覆盖来源策略、回调 CSP、多标签 Cookie 隔离、同账号新 PAT 通过、其他账号/过期 PAT 拒绝及上游异常。原生表单测试不可用手工设置 Origin 的 HTTP 脚本替代：`no-referrer` 会使浏览器表单 POST 带 `Origin:null`。当前使用 `strict-origin`，且仍拒绝空值或 null 来源。

## 2026-09-15 已定位的 FC 行为

- 即使 GetFunction 回显 `customRuntimeConfig.port=8080`，新加坡实例启动错误仍提示检查 9000。部署脚本让程序直接监听 9000，实测恢复 200。
- fcapp.run 默认域名拦截 OAuth 的绝对重定向（ExternalRedirectForbidden），正式 OAuth 必须使用自定义域名。
- 自定义运行时已实测能通过环境变量中的 FC 角色临时凭据向私有 OSS 写入注册记录。

## 重跑验收

```sh
python3 scripts/smoke_http.py --write
```

该脚本完成 OAuth→MCP 握手→工具目录→真实读取→创建一个临时未分配 backlog 任务→更新/评论/读取/搜索→删除刚创建的任务→刷新/撤销。访问令牌仅留内存，报告不含凭据。普通读取验证可省略 `--write`。

## 当前验收结果（2026-09-15）

0.3.2 已通过阿里云 CLI 部署。正式域名 OAuth 授权/兑换/重放拒绝/刷新/撤销、MCP 初始化、16 工具目录和真实读取均通过。

ChatGPT 网页已使用同账号临时 PAT 完成 OAuth，显示已安装；刷新后加载全部 16 工具。ChatGPT 实际读取项目和任务，官网 wildflow 工作区基线为空（已通过 /api/workspaces 核对）。随后仅创建临时未分配 backlog 任务 WILD-2，完成描述更新、/note 评论和详情读取；通过官网 API 独立核对描述、评论与未分配状态后删除（HTTP 204），再次读取为 404。临时登录 PAT 未保存或部署，可撤销；原后端 PAT 保持不变。

逐项状态见 [verification.json](verification.json)，浏览器对话见 [真实验收记录](https://chatgpt.com/c/6aa8f3eb-916c-83ec-84d5-63c8724f5e0c)。此前 0.3.1 的协议级创建、清空描述、搜索及清理结果单独保留。
