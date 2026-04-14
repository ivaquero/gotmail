package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Account account data structure
type Account struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Password  string    `json:"password"`
	Token     TokenData `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

// TokenData JWT token data structure
type TokenData struct {
	Token string `json:"token"`
}

// Domain API returned domain data structure
type Domain struct {
	Domain string `json:"domain"`
}

// DomainResponse API domain response
type DomainResponse struct {
	HydraMember []Domain `json:"hydra:member"`
}

// Message email message data structure
type Message struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
}

// MessageResponse API message response
type MessageResponse struct {
	HydraMember []Message `json:"hydra:member"`
}

// EmailDetail email detail data structure
type EmailDetail struct {
	ID      string   `json:"id"`
	HTML    []string `json:"html"`
	Subject string   `json:"subject"`
}

// Database database operation structure for multiple accounts
type Database struct {
	dataPath string
	accounts map[string]*Account // Key: account ID, Value: account data
}

// NewDatabase creates new database instance
func NewDatabase(dataPath string) *Database {
	return &Database{
		dataPath: dataPath,
		accounts: make(map[string]*Account),
	}
}

// Read reads accounts data
func (db *Database) Read() error {
	data, err := os.ReadFile(db.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			db.accounts = make(map[string]*Account)
			return nil
		}
		return fmt.Errorf("failed to read accounts file: %w", err)
	}

	var accounts map[string]*Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		// Try to read as single account (backward compatibility)
		var singleAccount Account
		if err := json.Unmarshal(data, &singleAccount); err != nil {
			return fmt.Errorf("failed to unmarshal accounts data: %w", err)
		}
		// Convert single account to map format
		db.accounts = map[string]*Account{singleAccount.ID: &singleAccount}
		return nil
	}

	db.accounts = accounts
	return nil
}

// Write writes accounts data
func (db *Database) Write() error {
	if db.accounts == nil {
		return nil
	}

	data, err := json.MarshalIndent(db.accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal accounts data: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(db.dataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(db.dataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write accounts file: %w", err)
	}

	return nil
}

// GetData gets all accounts data
func (db *Database) GetData() map[string]*Account {
	return db.accounts
}

// GetAccount gets specific account by ID
func (db *Database) GetAccount(accountID string) *Account {
	if db.accounts == nil {
		return nil
	}
	return db.accounts[accountID]
}

// SetAccount sets specific account data
func (db *Database) SetAccount(accountID string, data *Account) {
	if db.accounts == nil {
		db.accounts = make(map[string]*Account)
	}
	db.accounts[accountID] = data
}

// AddAccount adds new account
func (db *Database) AddAccount(data *Account) error {
	if db.accounts == nil {
		db.accounts = make(map[string]*Account)
	}
	if _, exists := db.accounts[data.ID]; exists {
		return fmt.Errorf("account with ID %s already exists", data.ID)
	}
	db.accounts[data.ID] = data
	return nil
}

// DeleteAccount deletes specific account
func (db *Database) DeleteAccount(accountID string) error {
	if db.accounts == nil {
		return fmt.Errorf("no accounts found")
	}
	if _, exists := db.accounts[accountID]; !exists {
		return fmt.Errorf("account with ID %s not found", accountID)
	}
	delete(db.accounts, accountID)
	return nil
}

// GetAllAccountIDs gets all account IDs
func (db *Database) GetAllAccountIDs() []string {
	if db.accounts == nil {
		return []string{}
	}
	ids := make([]string, 0, len(db.accounts))
	for id := range db.accounts {
		ids = append(ids, id)
	}
	return ids
}

// DeleteData deletes accounts data file
func (db *Database) DeleteData() error {
	if err := os.Remove(db.dataPath); err != nil {
		return fmt.Errorf("failed to delete accounts file: %w", err)
	}
	db.accounts = make(map[string]*Account)
	return nil
}

// GenerateRandomString generates random string
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// If cryptographic random number generation fails, fall back to time seed
			b[i] = charset[i%len(charset)]
		} else {
			b[i] = charset[num.Int64()]
		}
	}
	return string(b)
}

// Spinner simple loading animation
type Spinner struct {
	message string
	done    chan bool
}

// NewSpinner creates new loading animation
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		done:    make(chan bool),
	}
}

// Start starts loading animation
func (s *Spinner) Start() {
	go func() {
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.done:
				fmt.Printf("\r%s\r", strings.Repeat(" ", len(s.message)+10))
				return
			default:
				fmt.Printf("\r%s %s", chars[i%len(chars)], s.message)
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

// Stop stops loading animation
func (s *Spinner) Stop() {
	close(s.done)
}

// Color simple color output functions
type Color struct{}

// Red red output
func (c Color) Red(text string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", text)
}

// Green green output
func (c Color) Green(text string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", text)
}

// Blue blue output
func (c Color) Blue(text string) string {
	return fmt.Sprintf("\033[34m%s\033[0m", text)
}

// Underline underline output
func (c Color) Underline(text string) string {
	return fmt.Sprintf("\033[4m%s\033[0m", text)
}

// Helper functions for multi-account management

// ParseAccountID parses account ID from command arguments
func ParseAccountID(args []string) (string, bool) {
	for i, arg := range args {
		if arg == "--id" && i+1 < len(args) {
			return args[i+1], true
		}
	}
	return "", false
}

// ParseEmailID parses email ID from command arguments
func ParseEmailID(args []string) (string, bool) {
	for i, arg := range args {
		if arg == "--email" && i+1 < len(args) {
			return args[i+1], true
		}
	}
	return "", false
}

// FormatAccountList formats account list for display
func FormatAccountList(accounts map[string]*Account) string {
	if len(accounts) == 0 {
		return "No accounts found"
	}

	var result strings.Builder
	color := &Color{}
	result.WriteString("Available accounts:\n")

	for id, account := range accounts {
		result.WriteString(fmt.Sprintf("  %s. %s (created: %s)\n",
			color.Green(id),
			color.Underline(account.Address),
			account.CreatedAt.Format("2006-01-02 15:04:05")))
	}

	return result.String()
}

// SelectAccount prompts user to select an account if no ID provided
func SelectAccount(accounts map[string]*Account) (*Account, error) {
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found")
	}

	if len(accounts) == 1 {
		// Auto-select if only one account
		for _, account := range accounts {
			return account, nil
		}
	}

	// Multiple accounts - show list and ask for selection
	fmt.Println(FormatAccountList(accounts))
	fmt.Print("Enter account ID: ")

	var selectedID string
	fmt.Scanln(&selectedID)

	account, exists := accounts[selectedID]
	if !exists {
		return nil, fmt.Errorf("account with ID %s not found", selectedID)
	}

	return account, nil
}
