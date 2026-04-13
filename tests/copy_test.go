package utils_test

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/ivaquero/gotmail/utils"
)

func TestCopy(t *testing.T) {
	// Check if current platform supports clipboard operations
	if !isClipboardSupported() {
		t.Skip("Clipboard function not supported on current platform, skipping test")
		return
	}

	// Test multiple data types
	testCases := []struct {
		name string
		data string
	}{
		{"English text", "Hello, World!"},
		{"Chinese text", "这是测试数据。"},
		{"Mixed text", "Hello, World! 这是测试数据。"},
		{"Special characters", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
		{"Multi-line text", "第一行\n第二行\n第三行"},
	}

	fmt.Println("=== Clipboard copy function test ===")

	for i, tc := range testCases {
		fmt.Printf("Test %d: %s\n", i+1, tc.name)
		fmt.Printf("Data: %s\n", tc.data)
		fmt.Print("Copying to clipboard...\n")

		if err := utils.Copy(tc.data); err != nil {
			log.Printf("Copy failed: %v\n", err)
			t.Errorf("Test %s failed: %v\n", tc.name, err)
			continue
		}

		fmt.Println(" Success!")
	}

	fmt.Println("=== All tests completed ===")
	fmt.Println("Tip: You can manually paste to verify clipboard content")
}

// isClipboardSupported checks if current platform supports clipboard operations
func isClipboardSupported() bool {
	switch runtime.GOOS {
	case "windows", "darwin":
		return true
	default:
		return false
	}
}
