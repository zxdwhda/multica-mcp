# multica-mcp 工作约定

- 本项目是 zxdwhda 的个人开源项目，正式仓库为 https://github.com/zxdwhda/multica-mcp 。不要放入 wildsyn 或其他组织名下；真实密钥、OAuth 状态和私有配置不得提交。
- 本目录是独立 Git 根和独立维护的产品，项目名为 `multica-mcp`。保留来源署名及 MIT 许可。
- 产品目标和具体范围以用户当前要求为准；当前接入场景是 ChatGPT 网页插件调用远程 Multica，部署方向为阿里云 FC。
- 来源与更新方式以 `UPSTREAM.md` 为准。不要引入来源仓库的分支、标签或历史，不创建 GitHub Fork 关系；按需导入选定补丁。
- 复用现有 Go MCP SDK、HTTP client 和代码分层。修改前核对工作区变化，保留已有工作。
- 写工具遵守 `MULTICA_READ_ONLY`；正确说明任务创建、分配和评论可能触发 Agent 执行的语义。
- 验证与改动相称：Go 代码或模块路径变化运行 `go test ./...`、`go vet ./...`、`go build .`；纯文档无需新增测试。
- `docs/upstream/` 是原项目文档的历史材料，里面的安装地址、规则、版本和功能承诺不作为本项目当前操作依据。
- 版本、仓库链接和发布方式由本项目维护；不得把原作者的发布结果当成本项目已发布。
