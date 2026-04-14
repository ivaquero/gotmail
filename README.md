# GoTMail

[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org/doc/go1.18)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivaquero/gotmail)](https://goreportcard.com/report/github.com/ivaquero/gotmail)
![code size](https://img.shields.io/github/languages/code-size/ivaquero/gotmail.svg)
![repo size](https://img.shields.io/github/repo-size/ivaquero/gotmail.svg)

GoTMail (Go Temporary Mail) is temporary email CLI tool written in Go for [mail.tm](https://mail.tm/), providing cross-platform temporary email management functionality.

**[中文版本](README-CN.md)**

## 🌟 Features

- **Temporary Email Creation** - Quickly create Mail.tm temporary email accounts
- **Message Management** - Fetch and list received email messages
- **Email Viewing** - Open specific emails in your browser
- **Account Management** - View account details and delete accounts
- **Cross-platform Support** - Support for Windows, macOS, and Linux
- **Clipboard Integration** - Automatically copy email addresses to clipboard
- **Command Line Interface** - Simple and easy-to-use CLI operations

## 🚀 Quick Start

### Installation

- Windows

```bash
scoop bucket add scoopforge/main-plus
scoop install gotmail
```

- macOS and Linux

```bash
brew tap brewforge/more
brew install gotmail
```

- Alternative method (requires Go 1.18 or higher)

```bash
go install github.com/ivaquero/gotmail
```

### Basic Usage

Create a new temporary email account:

```bash
gotmail create
```

View received emails:

```bash
gotmail messages
```

Open a specific email in your browser:

```bash
gotmail open 1
```

View account details:

```bash
gotmail details
```

Delete account:

```bash
gotmail delete
```

## 📖 Command Reference

|     Command     |             Description              |      Example       |
| :-------------: | :----------------------------------: | :----------------: |
|    `create`     | Create a new temporary email account |  `gotmail create`  |
|   `messages`    |      Fetch and list all emails       | `gotmail messages` |
| `open <number>` |   Open specified email in browser    |  `gotmail open 1`  |
|    `details`    |   Display current account details    | `gotmail details`  |
|    `delete`     |        Delete current account        |  `gotmail delete`  |

## 🔧 Development Guide

### Requirements

- Go 1.18 or higher
- Network connection (for Mail.tm API)
- Supported browser (for opening emails)

### Building the Project

```bash
git clone https://github.com/ivaquero/gotmail
cd gotmail
go build
```

### Running Tests

```bash
go test ./tests/... -v
```

### Code Standards

- Follow Go standard code formatting
- Use structured error handling
- Provide detailed error information
- Maintain cross-platform compatibility

## 🔒 Security Features

- **Cryptographic Random Generation** - Use `crypto/rand` to generate secure random strings
- **Error Fallback Mechanism** - Provide fallback solutions when cryptographic random generation fails
- **Input Validation** - Validate API responses and user inputs
- **Secure Data Storage** - Store account data securely in JSON format

## 🌐 API Integration

This project uses the Mail.tm API to provide temporary email services:

- **API Endpoint**: `https://api.mail.tm`
- **Feature Support**: Account creation, email retrieval, account deletion
- **Data Format**: JSON
- **Authentication**: Bearer Token

## 📝 Data Storage

Account data is stored in a local file:

- **File Path**: `<execution directory>/data/account.json`
- **Data Format**: JSON
- **Contains**: Account ID, email address, password, authentication token

## 🐛 Error Handling

The project implements comprehensive error handling mechanisms:

- **Network Errors** - Handle API connection failures
- **File Operation Errors** - Handle data read/write failures
- **Clipboard Errors** - Handle cross-platform clipboard operation failures
- **API Response Errors** - Handle API error status returns

## 🤝 Contributing Guidelines

Issues and Pull Requests are welcome to improve the project:

1. Fork the project repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add some amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- [mail.tm](https://mail.tm) for providing temporary email services

---

**Note**: This is a temporary email tool, please do not use it to receive important or sensitive information.
