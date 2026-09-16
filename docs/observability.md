# FC 可观测性与发布

通过 `deploy/deploy.sh` 从干净且已提交的当前源码构建，健康响应包含源码 revision。
函数日志写入新加坡 `wildflow-mcp-sg/multica-mcp`，请求指标和实例指标开启。

所有已注册工具通过同一接收中间件记录 event、operation、outcome、request_id 和 duration_ms；协议错误及 `isError=true` 都标记为 error，即使 HTTP 返回 200。不记录工具参数、返回正文、令牌或错误详情。HTTP 日志使用随机关联 ID；保留流式响应能力，不记录查询参数或请求正文。

共享 SLS 索引、告警、不可变版本和生产别名由独立 `feishu-mcp-serverless` 仓库的 `deploy/fc-ops.py` 管理：

1. 首次改动前 `snapshot --function multica-mcp --snapshot-dir <private-path>`，保存私有配置及回滚版本。
2. 本仓执行 `go test ./...`、`go vet ./...`、`go build -o bin/multica-mcp .` 并提交。
3. 执行 `deploy/deploy.sh` 更新 LATEST。
4. 共享入口执行 `release --function multica-mcp --revision <commit>`，发布版本并切换 prod。域名只调整本函数路由，保持 TLS 和飞书路由。
5. 执行 `python3 scripts/smoke_http.py`，验证 OAuth、395 工具列表及真实只读调用；不要使用会写业务数据的 `--write` 或 `--full-api`。
6. 如业务验收失败，执行 `rollback --function multica-mcp --version <previous-version>`。

生产域名及 HTTP 触发器使用 prod，后续部署 LATEST 不立即影响生产。健康检查原有参数保持不变；不启用预留实例、追踪或改变并发。

每个 Logstore 保留 7 天，1 shard。共享配置维护字段索引以及函数错误、平台错误、HTTP 5xx、两类限流、接近超时 6 条规则。复用现有通知联系人，空闲无数据不告警。HTTP 200 内的业务失败通过 SLS 工具统计查询，不会自动计入 FC 函数错误。通知渠道的实际投递需独立验证。
