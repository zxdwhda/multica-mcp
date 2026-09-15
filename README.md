# multica-mcp

独立维护的 Multica MCP 项目，通过 MCP 让 ChatGPT 网页插件连接远程 Multica，后续产品能力在本仓库继续开发。

本项目基于 [Igor Lazarev / strider2038 的 multica-mcp](https://github.com/strider2038/multica-mcp)，导入了包含 [PR #14](https://github.com/strider2038/multica-mcp/pull/14) 的源码快照。保留原作者版权与 MIT 许可；本仓库独立维护，不代表原作者或 Multica 官方。详细来源见 [UPSTREAM.md](UPSTREAM.md)。

## 当前状态

- 已有 15 个 MCP 工具、stdio / Streamable HTTP 传输、PAT 和工作区配置。
- 源码适配基线为 Multica v0.4.43；实际远程实例兼容性以联调为准。
- 阿里云 FC 自定义容器是已评估可行的部署方式，尚未部署。
- ChatGPT 网页认证接入、分页、部分失败反馈等仍需完善；导入源码不代表这些问题已修复。

## 本地构建

需要 Go 1.25 或更高版本。在本目录执行：

```sh
make build
go test ./... -race -count=1
go vet ./...
```

输出为 `bin/multica-mcp`。当前 Go module 为 `multica-mcp`，不从原作者仓库安装；独立远端确定后可按实际地址更新 module 路径。

## 运行配置

程序从进程环境变量读取配置；不会自动加载 `.env`。

| 环境变量 | 用途 |
| --- | --- |
| `MULTICA_BASE_URL` | Multica API 所在服务地址 |
| `MULTICA_TOKEN` | 已有 Multica PAT |
| `MULTICA_WORKSPACE_ID` / `MULTICA_WORKSPACE_SLUG` | 工作区；slug 优先 |
| `MCP_TRANSPORT` | `stdio` 或 `http` |
| `MCP_HTTP_PORT` | HTTP 端口，默认 8080 |
| `MCP_API_KEY` | 当前 HTTP 模式的静态 Bearer Key 校验 |
| `MULTICA_READ_ONLY` | 为 true 时禁用写工具 |

启动命令为 `./bin/multica-mcp`。HTTP MCP 路径使用 `/mcp`。静态 Key 校验尚不等同于 ChatGPT 网页 OAuth 接入。

## 开发与来源同步

- `main.go`：启动与传输。
- `internal/mcp/`：工具定义和参数处理。
- `internal/app/`：调用流程。
- `internal/multica/`：Multica HTTP API 适配。
- `internal/domain/`：数据模型。
- `docs/upstream/`：来源快照文档与导入文件哈希，仅用于追溯。

本仓库拥有独立 Git 历史。后续按需吸收来源项目的补丁或文件差异，不同步其分支、标签或 GitHub Fork 关系。操作约定见 [UPSTREAM.md](UPSTREAM.md)。

## 许可

见 [LICENSE](LICENSE)。
