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

// MailManager mail manager
type MailManager struct {
	db     *Database
	client *http.Client
	color  *Color
}

// NewMailManager creates new mail manager
func NewMailManager(dataPath string) *MailManager {
	return &MailManager{
		db: NewDatabase(dataPath),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		color: &Color{},
	}
}

// CreateAccount creates new account
func (m *MailManager) CreateAccount() error {
	spinner := NewSpinner("creating...")
	spinner.Start()
	defer spinner.Stop()

	// Read account data
	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	// Check if maximum accounts limit reached (optional, can be adjusted)
	accounts := m.db.GetData()
	if len(accounts) >= 10 {
		fmt.Printf("%s\n", m.color.Red("Maximum 10 accounts allowed"))
		return nil
	}

	// Get available domain
	domain, err := m.getDomain()
	if err != nil {
		return fmt.Errorf("failed to get domain: %w", err)
	}

	// Generate random email and password
	email := fmt.Sprintf("%s@%s", GenerateRandomString(7), domain)
	password := GenerateRandomString(10)

	// Create account
	accountData, err := m.createAccountAPI(email, password)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	// Get JWT token
	tokenData, err := m.getToken(email, password)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Build account information
	account := &Account{
		ID:        accountData["id"].(string),
		Address:   email,
		Password:  password,
		Token:     TokenData{Token: tokenData["token"].(string)},
		CreatedAt: time.Now(),
	}

	// Copy email to clipboard
	if err := Copy(email); err != nil {
		fmt.Printf("Warning: failed to copy email to clipboard: %v\n", err)
	}

	// Save account data
	if err := m.db.AddAccount(account); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}
	if err := m.db.Write(); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	fmt.Printf("\r%s: %s (ID: %s)\n", m.color.Blue("Account created"), m.color.Underline(m.color.Green(email)), m.color.Green(account.ID))
	return nil
}

// ExportAccount exports account data to specified folder
func (m *MailManager) ExportAccount(exportFolder string) error {
	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	accounts := m.db.GetData()
	if len(accounts) == 0 {
		return fmt.Errorf("no accounts found")
	}

	// Ensure export directory exists
	if err := os.MkdirAll(exportFolder, 0755); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	exportPath := filepath.Join(exportFolder, "accounts.json")

	// Read the original account file
	originalData, err := os.ReadFile(m.db.dataPath)
	if err != nil {
		return fmt.Errorf("failed to read original account file: %w", err)
	}

	// Write to export path
	if err := os.WriteFile(exportPath, originalData, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	fmt.Printf("Account data exported to: %s\n", exportPath)
	return nil
}

// ExportAccountByID exports specific account data to specified path
func (m *MailManager) ExportAccountByID(accountID string, exportFolder string) error {
	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetAccount(accountID)
	if account == nil {
		return fmt.Errorf("account with ID %s not found", accountID)
	}

	// Ensure export directory exists
	if err := os.MkdirAll(exportFolder, 0755); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	exportPath := filepath.Join(exportFolder, fmt.Sprintf("account_%s.json", accountID))

	// Create single account data for export
	singleAccountData := map[string]*Account{accountID: account}
	exportData, err := json.MarshalIndent(singleAccountData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal account data: %w", err)
	}

	// Write to export path
	if err := os.WriteFile(exportPath, exportData, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	fmt.Printf("Account %s exported to: %s\n", accountID, exportPath)
	return nil
}

// GetAllAccountsJSON returns all accounts data as JSON string
func (m *MailManager) GetAllAccountsJSON() (string, error) {
	if err := m.db.Read(); err != nil {
		return "", fmt.Errorf("failed to read database: %w", err)
	}

	accounts := m.db.GetData()
	if len(accounts) == 0 {
		return "[]", nil
	}

	// Convert accounts to JSON
	jsonData, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal accounts data: %w", err)
	}

	return string(jsonData), nil
}

// ListAccounts lists all accounts
func (m *MailManager) ListAccounts() error {
	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	accounts := m.db.GetData()
	if len(accounts) == 0 {
		fmt.Printf("%s\n", m.color.Red("No accounts found"))
		return nil
	}

	fmt.Println(FormatAccountList(accounts))
	return nil
}

// FetchMessagesByAccountID fetches messages for specific account
func (m *MailManager) FetchMessagesByAccountID(accountID string) ([]Message, error) {
	spinner := NewSpinner("fetching messages...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return nil, fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetAccount(accountID)
	if account == nil {
		return nil, fmt.Errorf("account with ID %s not found", accountID)
	}

	messages, err := m.fetchMessagesAPI(account.Token.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	if len(messages) == 0 {
		fmt.Printf("\r\033[K%s\n", m.color.Red("No Emails"))
		return nil, nil
	}

	return messages, nil
}

// ShowAccountDetails shows details for specific account
func (m *MailManager) ShowAccountDetails(accountID string) error {
	spinner := NewSpinner("fetching details...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetAccount(accountID)
	if account == nil {
		return fmt.Errorf("account with ID %s not found", accountID)
	}

	// Use existing account data instead of API call for basic info
	fmt.Printf("\n    Account ID: %s\n    Email: %s\n    Created: %s\n",
		m.color.Green(account.ID),
		m.color.Underline(m.color.Green(account.Address)),
		m.color.Green(account.CreatedAt.Format("2006-01-02 15:04:05")))

	return nil
}

// DeleteAccountByID deletes specific account by ID
func (m *MailManager) DeleteAccountByID(accountID string) error {
	spinner := NewSpinner("deleting account...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetAccount(accountID)
	if account == nil {
		return fmt.Errorf("account with ID %s not found", accountID)
	}

	if err := m.deleteAccountAPI(account.ID, account.Token.Token); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if err := m.db.DeleteAccount(accountID); err != nil {
		return fmt.Errorf("failed to delete account data: %w", err)
	}

	if err := m.db.Write(); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	fmt.Printf("%s\n", m.color.Blue("Account deleted"))
	return nil
}

// FetchMessages fetches email messages for the first account (backward compatibility)
func (m *MailManager) FetchMessages() ([]Message, error) {
	spinner := NewSpinner("fetching...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return nil, fmt.Errorf("failed to read database: %w", err)
	}

	accounts := m.db.GetData()
	if len(accounts) == 0 {
		fmt.Printf("%s\n", m.color.Red("No accounts found"))
		return nil, nil
	}

	// Get first account for backward compatibility
	var account *Account
	for _, acc := range accounts {
		account = acc
		break
	}

	messages, err := m.fetchMessagesAPI(account.Token.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	if len(messages) == 0 {
		fmt.Printf("\r\033[K%s\n", m.color.Red("No Emails"))
		return nil, nil
	}

	return messages, nil
}

// DeleteAccount deletes the first account (backward compatibility)
func (m *MailManager) DeleteAccount() error {
	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	accounts := m.db.GetData()
	if len(accounts) == 0 {
		fmt.Printf("%s\n", m.color.Red("No accounts found"))
		return nil
	}

	// Use interactive selection like SelectAccount
	account, err := SelectAccount(accounts)
	if err != nil {
		return err
	}

	spinner := NewSpinner("deleting...")
	spinner.Start()
	defer spinner.Stop()

	// Find the account ID from the map
	var accountID string
	for id, acc := range accounts {
		if acc == account {
			accountID = id
			break
		}
	}

	if err := m.deleteAccountAPI(account.ID, account.Token.Token); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if err := m.db.DeleteAccount(accountID); err != nil {
		return fmt.Errorf("failed to delete account data: %w", err)
	}

	if err := m.db.Write(); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	fmt.Printf("%s\n", m.color.Blue("Account deleted"))
	return nil
}

// ShowDetails shows account details (backward compatibility - shows first account)
func (m *MailManager) ShowDetails() error {
	spinner := NewSpinner("fetching details...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	accounts := m.db.GetData()
	if len(accounts) == 0 {
		fmt.Printf("%s\n", m.color.Red("No accounts found"))
		return nil
	}

	// Get first account for backward compatibility
	var account *Account
	for _, acc := range accounts {
		account = acc
		break
	}

	// Use existing account data instead of API call for basic info
	fmt.Printf("\n    Account ID: %s\n    Email: %s\n    Created: %s\n",
		m.color.Green(account.ID),
		m.color.Underline(m.color.Green(account.Address)),
		m.color.Green(account.CreatedAt.Format("2006-01-02 15:04:05")))

	return nil
}

// OpenEmail opens specified email (backward compatibility - uses first account)
func (m *MailManager) OpenEmail(emailIndex int) error {
	emailFilePath := filepath.Join(getCurrentDir(), "../data/email.html")
	spinner := NewSpinner("opening...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	accounts := m.db.GetData()
	if len(accounts) == 0 {
		return fmt.Errorf("no accounts found")
	}

	// Get first account for backward compatibility
	var account *Account
	for _, acc := range accounts {
		account = acc
		break
	}

	messages, err := m.fetchMessagesAPI(account.Token.Token)
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
		fmt.Printf("%s\n", m.color.Red("No HTML content found"))
		return nil
	}

	// Write HTML file
	dir := filepath.Dir(emailFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(emailFilePath, []byte(emailDetail.HTML[0]), 0644); err != nil {
		return fmt.Errorf("failed to write email file: %w", err)
	}

	// Open file in browser
	if err := openInBrowser(emailFilePath); err != nil {
		return fmt.Errorf("failed to open email in browser: %w", err)
	}

	return nil
}

// OpenEmailByAccountID opens specified email for specific account
func (m *MailManager) OpenEmailByAccountID(accountID string, emailIndex int) error {
	emailFilePath := filepath.Join(getCurrentDir(), "../data/email.html")
	spinner := NewSpinner("opening...")
	spinner.Start()
	defer spinner.Stop()

	if err := m.db.Read(); err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	account := m.db.GetAccount(accountID)
	if account == nil {
		return fmt.Errorf("account with ID %s not found", accountID)
	}

	messages, err := m.fetchMessagesAPI(account.Token.Token)
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
		fmt.Printf("%s\n", m.color.Red("No HTML content found"))
		return nil
	}

	// Write HTML file
	dir := filepath.Dir(emailFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(emailFilePath, []byte(emailDetail.HTML[0]), 0644); err != nil {
		return fmt.Errorf("failed to write email file: %w", err)
	}

	// Open file in browser
	if err := openInBrowser(emailFilePath); err != nil {
		return fmt.Errorf("failed to open email in browser: %w", err)
	}

	return nil
}

// getDomain gets available domain
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

// createAccountAPI calls API to create account
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

// getToken gets JWT token
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

// fetchMessagesAPI fetches email messages API
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

// deleteAccountAPI deletes account API
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

// getAccountDetailsAPI gets account details API
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

// getEmailDetailAPI gets email detail API
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

// getCurrentDir gets current execution directory
func getCurrentDir() string {
	ex, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(ex)
}

// openInBrowser opens file in browser
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
