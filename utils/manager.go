package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// MailManager 邮件管理器
type MailManager struct {
	db     *Database
	client *http.Client
	color  *Color
}

// NewMailManager 创建新的邮件管理器
func NewMailManager(dataPath string) *MailManager {
	return &MailManager{
		db: NewDatabase(dataPath),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		color: &Color{},
	}
}

// CreateAccount 创建新账户
func (m *MailManager) CreateAccount() error {
	spinner := NewSpinner("creating...")
	spinner.Start()
	defer spinner.Stop()

	// 读取账户数据
	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	// 如果账户已存在
	if m.db.GetData() != nil {
		fmt.Printf("\r%s\n", m.color.Red("Account already exists"))
		return nil
	}

	// 获取可用域名
	domain, err := m.getDomain()
	if err != nil {
		return fmt.Errorf("failed to get domain: %w", err)
	}

	// 生成随机邮箱和密码
	email := fmt.Sprintf("%s@%s", GenerateRandomString(7), domain)
	password := GenerateRandomString(10)

	// 创建账户
	accountData, err := m.createAccountAPI(email, password)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	// 获取JWT令牌
	tokenData, err := m.getToken(email, password)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// 构建账户信息
	account := &Account{
		ID:        accountData["id"].(string),
		Address:   email,
		Password:  password,
		Token:     TokenData{Token: tokenData["token"].(string)},
		CreatedAt: time.Now(),
	}

	// 复制邮箱到剪贴板
	if err := Copy(email); err != nil {
		fmt.Printf("Warning: failed to copy email to clipboard: %v\n", err)
	}

	// 保存账户数据
	m.db.SetData(account)
	if err := m.db.Write(); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	fmt.Printf("\r%s: %s\n", m.color.Blue("Account created"), m.color.Underline(m.color.Green(email)))
	return nil
}

// FetchMessages 获取邮件消息
func (m *MailManager) FetchMessages() ([]Message, error) {
	spinner := NewSpinner("fetching...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return nil, fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetData()
	if account == nil {
		fmt.Printf("\r%s\n", m.color.Red("Account not created yet"))
		return nil, nil
	}

	messages, err := m.fetchMessagesAPI(account.Token.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	if len(messages) == 0 {
		fmt.Printf("\r%s\n", m.color.Red("No Emails"))
		return nil, nil
	}

	return messages, nil
}

// DeleteAccount 删除账户
func (m *MailManager) DeleteAccount() error {
	spinner := NewSpinner("deleting...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetData()
	if account == nil {
		fmt.Printf("\r%s\n", m.color.Red("Account not created yet"))
		return nil
	}

	if err := m.deleteAccountAPI(account.ID, account.Token.Token); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if err := m.db.DeleteData(); err != nil {
		return fmt.Errorf("failed to delete account data: %w", err)
	}

	fmt.Printf("\r%s\n", m.color.Blue("Account deleted"))
	return nil
}

// ShowDetails 显示账户详情
func (m *MailManager) ShowDetails() error {
	spinner := NewSpinner("fetching details...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetData()
	if account == nil {
		fmt.Printf("\r%s\n", m.color.Red("Account not created yet"))
		return nil
	}

	details, err := m.getAccountDetailsAPI(account.ID, account.Token.Token)
	if err != nil {
		return fmt.Errorf("failed to get account details: %w", err)
	}

	fmt.Printf("\n    Email: %s\n    createdAt: %s\n",
		m.color.Underline(m.color.Green(details["address"].(string))),
		m.color.Green(time.Now().Format("2006-01-02 15:04:05")))

	return nil
}

// OpenEmail 打开指定邮件
func (m *MailManager) OpenEmail(emailIndex int) error {
	emailFilePath := filepath.Join(getCurrentDir(), "../data/email.html")
	spinner := NewSpinner("opening...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetData()
	if account == nil {
		return fmt.Errorf("account not found")
	}

	messages, err := m.FetchMessages()
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	if len(messages) == 0 || emailIndex < 1 || emailIndex > len(messages) {
		return fmt.Errorf("invalid email index")
	}

	mailToOpen := messages[emailIndex-1]
	emailDetail, err := m.getEmailDetailAPI(mailToOpen.ID, account.Token.Token)
	if err != nil {
		return fmt.Errorf("failed to get email detail: %w", err)
	}

	if len(emailDetail.HTML) == 0 {
		fmt.Printf("\r%s\n", m.color.Red("No HTML content found"))
		return nil
	}

	// 写入HTML文件
	dir := filepath.Dir(emailFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(emailFilePath, []byte(emailDetail.HTML[0]), 0644); err != nil {
		return fmt.Errorf("failed to write email file: %w", err)
	}

	// 在浏览器中打开文件
	if err := openInBrowser(emailFilePath); err != nil {
		return fmt.Errorf("failed to open email in browser: %w", err)
	}

	return nil
}

// getDomain 获取可用域名
func (m *MailManager) getDomain() (string, error) {
	resp, err := m.client.Get("https://api.mail.tm/domains?page=1")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var domainResp DomainResponse
	if err := json.NewDecoder(resp.Body).Decode(&domainResp); err != nil {
		return "", err
	}

	if len(domainResp.HydraMember) == 0 {
		return "", fmt.Errorf("no domains available")
	}

	return domainResp.HydraMember[0].Domain, nil
}

// createAccountAPI 调用API创建账户
func (m *MailManager) createAccountAPI(email, password string) (map[string]interface{}, error) {
	payload := map[string]string{
		"address":  email,
		"password": password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Post("https://api.mail.tm/accounts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// getToken 获取JWT令牌
func (m *MailManager) getToken(email, password string) (map[string]interface{}, error) {
	payload := map[string]string{
		"address":  email,
		"password": password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Post("https://api.mail.tm/token", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// fetchMessagesAPI 获取邮件消息API
func (m *MailManager) fetchMessagesAPI(token string) ([]Message, error) {
	req, err := http.NewRequest("GET", "https://api.mail.tm/messages", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var msgResp MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, err
	}

	return msgResp.HydraMember, nil
}

// deleteAccountAPI 删除账户API
func (m *MailManager) deleteAccountAPI(accountID, token string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.mail.tm/accounts/%s", accountID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete account, status: %d", resp.StatusCode)
	}

	return nil
}

// getAccountDetailsAPI 获取账户详情API
func (m *MailManager) getAccountDetailsAPI(accountID, token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.mail.tm/accounts/%s", accountID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// getEmailDetailAPI 获取邮件详情API
func (m *MailManager) getEmailDetailAPI(messageID, token string) (*EmailDetail, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.mail.tm/messages/%s", messageID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmailDetail
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// getCurrentDir 获取当前执行目录
func getCurrentDir() string {
	ex, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(ex)
}

// openInBrowser 在浏览器中打开文件
func openInBrowser(filePath string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", filePath}
	case "darwin":
		cmd = "open"
		args = []string{filePath}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{filePath}
	}

	return exec.Command(cmd, args...).Start()
}
