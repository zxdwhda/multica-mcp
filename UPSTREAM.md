# 来源与独立维护

## 导入来源

- 项目：[strider2038/multica-mcp](https://github.com/strider2038/multica-mcp)
- 原作者：Igor Lazarev
- 许可：MIT，原始 `LICENSE` 原文保留。
- 导入日期：2026-09-15
- 源码提交：`27e41208fa0e45265749c4a95a4bdcd4e3418692`
- 包含：[PR #14 — sync Multica REST API v0.4.43](https://github.com/strider2038/multica-mcp/pull/14)，导入时尚未在来源仓库合并。
- 文件哈希：[source-manifest.json](docs/upstream/source-manifest.json)，记录导入前源码快照。

仅复制该提交的 29 个 Git 跟踪文件，没有复制 `.git`、分支、标签、远端配置、构建产物或本地环境文件。在当前目录重新初始化独立 Git 历史。

原 README、CHANGELOG、AGENTS 和俄语使用文档保存在 `docs/upstream/`，用于说明来源历史，不是当前项目文档。

## 导入后的调整

- 项目名继续使用 `multica-mcp`。
- Go module、内部 import 和构建链接参数改为本项目路径 `multica-mcp`，避免依赖原作者的安装地址。
- 默认程序版本与导入快照的 `VERSION` 统一为 0.3.0。
- 新建本项目 README、工作约定与变更记录；来源版权不变。

## 后续同步

0.4.0 另外读取了官方 [multica-ai/multica](https://github.com/multica-ai/multica/tree/7ebe0bf58d99238ccf830c7ed1aa658d4f5807d0) 的服务端路由和请求声明，生成用户 API 接口目录。目录只记录路径、方法、字段和来源位置，不包含该项目的业务处理实现。生成器为本项目独立编写，API 目录与原 MCP 来源分开维护。完整覆盖见 [API-COVERAGE.md](docs/API-COVERAGE.md)。

1. 在本仓库之外的临时目录读取来源项目指定提交，确认需要吸收的内容。
2. 基于上次采用的来源提交生成文件差异，仅应用选定的源码、测试和必要文档补丁；适配本项目的 module 路径与已有修改。
3. 在本项目提交中记录来源 URL、原提交 SHA、采用范围及验证结果。
4. 不使用 mirror 同步，不拉取来源分支/标签到本仓库，不用来源快照覆盖本项目自己的功能。

这种同步方式保留来源可追溯性，提交和分支均由本项目维护。随着独立改动增加，来源补丁可能需要手工适配。

## GitHub 发布方式

创建同名的普通空仓库 `multica-mcp`，把它设置为本项目 `origin` 后推送自己的 `main`。不通过 GitHub 的 Fork 功能创建，也不推送其他仓库的分支或标签。

当前尚未创建 GitHub 仓库或设置远端；GitHub 上的归属与可见性在实际创建时确定。
