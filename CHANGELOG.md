# Changelog

## 0.3.2

- 修复浏览器授权表单因 no-referrer 发送 Origin:null 而被拒绝的问题；保留严格来源校验。
- CSP 允许跳转至已注册回调地址的来源，登录 Cookie 按授权请求隔离，避免多标签互相覆盖。
- 支持同一 Multica 账号的新 PAT 登录，通过 /api/me 校验身份；临时登录 PAT 不保存、不替换后端凭据。
- 授权失败向浏览器显示可理解的错误和重新登录入口。

## 0.3.1

- 增加 ChatGPT OAuth、私有 OSS 状态存储和 FC CLI 部署。
- Streamable HTTP 改为无状态 JSON 响应，加入 namespaced 路由与健康检查。
- 增加分页、状态目录、结构化结果及读写工具注解。
- 修正搜索筛选、清空描述、详情部分失败、批量部分成功和错误字段丢失。
- 保留 MIT 来源信息和独立 Git 历史。


## Unreleased

- 从已包含 PR #14 的来源快照建立独立 `multica-mcp` 项目。
- 独立初始化 Git，仅保留本项目历史和分支。
- 调整 Go module、构建路径和项目文档；保留 MIT 许可与来源记录。
- 统一程序默认版本与 `VERSION` 为 0.3.0。

导入之前的来源发布记录见 [历史 CHANGELOG](docs/upstream/CHANGELOG.md)。本项目尚未发布。
