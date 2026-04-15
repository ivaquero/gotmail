package utils_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ivaquero/gotmail/utils"
)

func TestDatabase(t *testing.T) {
	fmt.Println("=== Test Database operations ===")

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "gotmail_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data file path
	testDataPath := filepath.Join(tempDir, "test_accounts.json")

	// Test creating new database
	db := utils.NewDatabase(testDataPath)
	if db == nil {
		t.Fatal("NewDatabase returned nil")
	}

	// Test empty database read
	if err := db.Read(); err != nil {
		t.Errorf("Failed to read empty database: %v", err)
	}

	// Test adding account
	testAccount := &utils.Account{
		ID:        "test123",
		Address:   "test@example.com",
		Password:  "testpass123",
		Token:     utils.TokenData{Token: "test_token_123"},
		CreatedAt: time.Now(),
	}

	if err := db.AddAccount(testAccount); err != nil {
		t.Errorf("Failed to add account: %v", err)
	}

	// Test writing database
	if err := db.Write(); err != nil {
		t.Errorf("Failed to write database: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}

	// Test reading database with data
	newDb := utils.NewDatabase(testDataPath)
	if err := newDb.Read(); err != nil {
		t.Errorf("Failed to read database with data: %v", err)
	}

	// Test getting account
	retrievedAccount := newDb.GetAccount("test123")
	if retrievedAccount == nil {
		t.Error("Failed to retrieve account")
	} else if retrievedAccount.Address != "test@example.com" {
		t.Errorf("Retrieved account has wrong address: %s", retrievedAccount.Address)
	}

	// Test getting all account IDs
	accountIDs := newDb.GetAllAccountIDs()
	if len(accountIDs) != 1 {
		t.Errorf("Expected 1 account ID, got %d", len(accountIDs))
	}

	// Test deleting account
	if err := newDb.DeleteAccount("test123"); err != nil {
		t.Errorf("Failed to delete account: %v", err)
	}

	// Verify account was deleted
	if newDb.GetAccount("test123") != nil {
		t.Error("Account was not properly deleted")
	}

	fmt.Println("Database operations test passed!")
}

func TestDatabaseMultipleAccounts(t *testing.T) {
	fmt.Println("=== Test Database with multiple accounts ===")

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "gotmail_test_multi")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data file path
	testDataPath := filepath.Join(tempDir, "test_multi_accounts.json")
	db := utils.NewDatabase(testDataPath)

	// Add multiple accounts
	accounts := []*utils.Account{
		{
			ID:        "acc1",
			Address:   "user1@example.com",
			Password:  "pass1",
			Token:     utils.TokenData{Token: "token1"},
			CreatedAt: time.Now(),
		},
		{
			ID:        "acc2",
			Address:   "user2@example.com",
			Password:  "pass2",
			Token:     utils.TokenData{Token: "token2"},
			CreatedAt: time.Now(),
		},
		{
			ID:        "acc3",
			Address:   "user3@example.com",
			Password:  "pass3",
			Token:     utils.TokenData{Token: "token3"},
			CreatedAt: time.Now(),
		},
	}

	for _, account := range accounts {
		if err := db.AddAccount(account); err != nil {
			t.Errorf("Failed to add account %s: %v", account.ID, err)
		}
	}

	// Test writing multiple accounts
	if err := db.Write(); err != nil {
		t.Errorf("Failed to write database with multiple accounts: %v", err)
	}

	// Test reading multiple accounts
	newDb := utils.NewDatabase(testDataPath)
	if err := newDb.Read(); err != nil {
		t.Errorf("Failed to read database with multiple accounts: %v", err)
	}

	// Verify all accounts were saved
	accountIDs := newDb.GetAllAccountIDs()
	if len(accountIDs) != 3 {
		t.Errorf("Expected 3 account IDs, got %d", len(accountIDs))
	}

	// Test getting data
	allData := newDb.GetData()
	if len(allData) != 3 {
		t.Errorf("Expected 3 accounts in data, got %d", len(allData))
	}

	fmt.Println("Multiple accounts test passed!")
}

func TestDatabaseErrors(t *testing.T) {
	fmt.Println("=== Test Database error handling ===")

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "gotmail_test_errors")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data file path
	testDataPath := filepath.Join(tempDir, "test_errors.json")
	db := utils.NewDatabase(testDataPath)

	// Test deleting non-existent account
	err = db.DeleteAccount("nonexistent")
	if err == nil {
		t.Error("Expected error when deleting non-existent account")
	}

	// Test adding duplicate account
	testAccount := &utils.Account{
		ID:        "duplicate",
		Address:   "test@example.com",
		Password:  "testpass",
		Token:     utils.TokenData{Token: "test_token"},
		CreatedAt: time.Now(),
	}

	if err := db.AddAccount(testAccount); err != nil {
		t.Errorf("Failed to add account: %v", err)
	}

	// Try to add the same account again
	err = db.AddAccount(testAccount)
	if err == nil {
		t.Error("Expected error when adding duplicate account")
	}

	fmt.Println("Database error handling test passed!")
}
