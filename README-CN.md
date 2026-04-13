# GoTMail

一个用 Go 语言编写的 [mail.tm](https://mail.tm/) 临时邮箱 CLI 工具，提供跨平台的临时邮箱管理功能。

## 🌟 功能特性

- **临时邮箱创建** - 快速创建 Mail.tm 临时邮箱账户
- **消息管理** - 获取和列出收到的邮件消息
- **邮件查看** - 在浏览器中打开特定邮件
- **账户管理** - 查看账户详情和删除账户
- **跨平台支持** - 支持 Windows、macOS 和 Linux
- **剪贴板集成** - 自动复制邮箱地址到剪贴板
- **命令行界面** - 简单易用的 CLI 操作

## 🚀 快速开始

### 安装

确保您的系统已安装 Go 1.18 或更高版本，然后运行：

```bash
go install github.com/ivaquero/gotmail
```

### 基本用法

创建新的临时邮箱账户：

```bash
gotmail create
```

查看收到的邮件：

```bash
gotmail messages
```

在浏览器中打开特定邮件：

```bash
gotmail open 1
```

查看账户详情：

```bash
gotmail details
```

删除账户：

```bash
gotmail delete
```

## 📖 命令说明

|      命令       |          描述          |        示例        |
| :-------------: | :--------------------: | :----------------: |
|    `create`     |  创建新的临时邮箱账户  |  `gotmail create`  |
|   `messages`    |   获取并列出所有邮件   | `gotmail messages` |
| `open <number>` | 在浏览器中打开指定邮件 |  `gotmail open 1`  |
|    `details`    |    显示当前账户详情    | `gotmail details`  |
|    `delete`     |      删除当前账户      |  `gotmail delete`  |

## 🔧 开发指南

### 环境要求

- Go 1.18 或更高版本
- 网络连接（用于 Mail.tm API）
- 支持的浏览器（用于打开邮件）

### 构建项目

```bash
git clone https://github.com/ivaquero/gotmail
cd gotmail
go build
```

### 运行测试

```bash
go test ./tests/... -v
```

### 代码规范

- 遵循 Go 标准代码格式
- 使用结构化的错误处理
- 提供详细的错误信息
- 保持跨平台兼容性

## 🔒 安全特性

- **加密随机数生成** - 使用 `crypto/rand` 生成安全的随机字符串
- **错误回退机制** - 在加密随机数生成失败时提供回退方案
- **输入验证** - 对 API 响应和用户输入进行验证
- **安全的数据存储** - 账户数据以 JSON 格式安全存储

## 🌐 API 集成

本项目使用 Mail.tm API 提供临时邮箱服务：

- **API 端点**: `https://api.mail.tm`
- **功能支持**: 账户创建、邮件获取、账户删除
- **数据格式**: JSON
- **认证方式**: Bearer Token

## 📝 数据存储

账户数据存储在本地文件中：

- **文件路径**: `<执行目录>/data/account.json`
- **数据格式**: JSON
- **包含信息**: 账户 ID、邮箱地址、密码、认证令牌

## 🐛 错误处理

项目实现了完善的错误处理机制：

- **网络错误** - 处理 API 连接失败
- **文件操作错误** - 处理数据读写失败
- **剪贴板错误** - 处理跨平台剪贴板操作失败
- **API 响应错误** - 处理 API 返回的错误状态

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目：

1. Fork 项目仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可。

## 🙏 致谢

- [mail.tm](https://mail.tm) 提供临时邮箱服务

---

**注意**: 这是一个临时邮箱工具，请勿用于接收重要或敏感信息。
