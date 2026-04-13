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

// Database database operation structure
type Database struct {
	dataPath string
	data     *Account
}

// NewDatabase creates new database instance
func NewDatabase(dataPath string) *Database {
	return &Database{
		dataPath: dataPath,
	}
}

// Read reads account data
func (db *Database) Read() error {
	data, err := os.ReadFile(db.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			db.data = nil
			return nil
		}
		return fmt.Errorf("failed to read account file: %w", err)
	}

	var account Account
	if err := json.Unmarshal(data, &account); err != nil {
		return fmt.Errorf("failed to unmarshal account data: %w", err)
	}

	db.data = &account
	return nil
}

// Write writes account data
func (db *Database) Write() error {
	if db.data == nil {
		return nil
	}

	data, err := json.MarshalIndent(db.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal account data: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(db.dataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(db.dataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write account file: %w", err)
	}

	return nil
}

// GetData gets account data
func (db *Database) GetData() *Account {
	return db.data
}

// SetData sets account data
func (db *Database) SetData(data *Account) {
	db.data = data
}

// DeleteData deletes account data file
func (db *Database) DeleteData() error {
	if err := os.Remove(db.dataPath); err != nil {
		return fmt.Errorf("failed to delete account file: %w", err)
	}
	db.data = nil
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
