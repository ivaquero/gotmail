package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ivaquero/gotmail/utils"
)

func TestParseAccountID(t *testing.T) {
	fmt.Println("=== Test ParseAccountID function ===")

	testCases := []struct {
		name          string
		args          []string
		expectedID    string
		expectedFound bool
	}{
		{
			name:          "Valid account ID",
			args:          []string{"msg", "--id", "abc123"},
			expectedID:    "abc123",
			expectedFound: true,
		},
		{
			name:          "Valid account ID with more args",
			args:          []string{"msg", "--id", "xyz789", "--other", "value"},
			expectedID:    "xyz789",
			expectedFound: true,
		},
		{
			name:          "No account ID",
			args:          []string{"msg", "--other", "value"},
			expectedID:    "",
			expectedFound: false,
		},
		{
			name:          "Empty args",
			args:          []string{},
			expectedID:    "",
			expectedFound: false,
		},
		{
			name:          "--id at end without value",
			args:          []string{"msg", "--id"},
			expectedID:    "",
			expectedFound: false,
		},
		{
			name:          "--id in middle with flag as value",
			args:          []string{"msg", "--id", "--other"},
			expectedID:    "--other",
			expectedFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, found := utils.ParseAccountID(tc.args)

			if id != tc.expectedID {
				t.Errorf("Expected ID '%s', got '%s'", tc.expectedID, id)
			}

			if found != tc.expectedFound {
				t.Errorf("Expected found=%v, got %v", tc.expectedFound, found)
			}
		})
	}

	fmt.Println("ParseAccountID test passed!")
}

func TestParseEmailID(t *testing.T) {
	fmt.Println("=== Test ParseEmailID function ===")

	testCases := []struct {
		name          string
		args          []string
		expectedID    string
		expectedFound bool
	}{
		{
			name:          "Valid email ID",
			args:          []string{"open", "--email", "email123"},
			expectedID:    "email123",
			expectedFound: true,
		},
		{
			name:          "Valid email ID with more args",
			args:          []string{"open", "--email", "email456", "--other", "value"},
			expectedID:    "email456",
			expectedFound: true,
		},
		{
			name:          "No email ID",
			args:          []string{"open", "--other", "value"},
			expectedID:    "",
			expectedFound: false,
		},
		{
			name:          "Empty args",
			args:          []string{},
			expectedID:    "",
			expectedFound: false,
		},
		{
			name:          "--email at end without value",
			args:          []string{"open", "--email"},
			expectedID:    "",
			expectedFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, found := utils.ParseEmailID(tc.args)

			if id != tc.expectedID {
				t.Errorf("Expected ID '%s', got '%s'", tc.expectedID, id)
			}

			if found != tc.expectedFound {
				t.Errorf("Expected found=%v, got %v", tc.expectedFound, found)
			}
		})
	}

	fmt.Println("ParseEmailID test passed!")
}

func TestFormatAccountList(t *testing.T) {
	fmt.Println("=== Test FormatAccountList function ===")

	// Test with empty accounts
	emptyResult := utils.FormatAccountList(map[string]*utils.Account{})
	if emptyResult != "No accounts found" {
		t.Errorf("Expected 'No accounts found', got '%s'", emptyResult)
	}

	// Test with accounts
	testAccounts := map[string]*utils.Account{
		"acc1": {
			ID:        "acc1",
			Address:   "user1@example.com",
			Password:  "pass1",
			CreatedAt: time.Now(),
		},
		"acc2": {
			ID:        "acc2",
			Address:   "user2@example.com",
			Password:  "pass2",
			CreatedAt: time.Now(),
		},
	}

	result := utils.FormatAccountList(testAccounts)

	// Verify result contains expected information
	expectedParts := []string{
		"Available accounts:",
		"acc1",
		"user1@example.com",
		"acc2",
		"user2@example.com",
	}

	for _, part := range expectedParts {
		if !contains(result, part) {
			t.Errorf("Formatted account list missing expected part: %s", part)
		}
	}

	fmt.Println("FormatAccountList test passed!")
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
