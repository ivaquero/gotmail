package utils_test

import (
	"fmt"
	"testing"

	"github.com/ivaquero/gotmail/utils"
)

func TestGenerateRandomString(t *testing.T) {
	fmt.Println("=== Test random string generation ===")

	// Generate multiple random strings to verify functionality
	for i := 0; i < 5; i++ {
		randomStr := utils.GenerateRandomString(10)
		fmt.Printf("Random string %d: %s\n", i+1, randomStr)

		// Verify string length
		if len(randomStr) != 10 {
			t.Errorf("Generated string length incorrect: expected 10, actual %d", len(randomStr))
		}

		// Verify string only contains allowed characters
		const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
		for _, char := range randomStr {
			found := false
			for _, allowed := range charset {
				if char == allowed {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Generated string contains invalid character: %c", char)
			}
		}
	}

	fmt.Println("\n=== Test completed ===")
}
