package utils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ivaquero/gotmail/utils"
)

func TestColor(t *testing.T) {
	fmt.Println("=== Test Color functions ===")

	color := &utils.Color{}

	// Test Red function
	redText := color.Red("test text")
	if !strings.Contains(redText, "\033[31m") {
		t.Error("Red function should contain ANSI red color code")
	}
	if !strings.Contains(redText, "\033[0m") {
		t.Error("Red function should contain ANSI reset code")
	}
	if !strings.Contains(redText, "test text") {
		t.Error("Red function should contain the original text")
	}

	// Test Green function
	greenText := color.Green("test text")
	if !strings.Contains(greenText, "\033[32m") {
		t.Error("Green function should contain ANSI green color code")
	}
	if !strings.Contains(greenText, "\033[0m") {
		t.Error("Green function should contain ANSI reset code")
	}
	if !strings.Contains(greenText, "test text") {
		t.Error("Green function should contain the original text")
	}

	// Test Blue function
	blueText := color.Blue("test text")
	if !strings.Contains(blueText, "\033[34m") {
		t.Error("Blue function should contain ANSI blue color code")
	}
	if !strings.Contains(blueText, "\033[0m") {
		t.Error("Blue function should contain ANSI reset code")
	}
	if !strings.Contains(blueText, "test text") {
		t.Error("Blue function should contain the original text")
	}

	// Test Underline function
	underlineText := color.Underline("test text")
	if !strings.Contains(underlineText, "\033[4m") {
		t.Error("Underline function should contain ANSI underline code")
	}
	if !strings.Contains(underlineText, "\033[0m") {
		t.Error("Underline function should contain ANSI reset code")
	}
	if !strings.Contains(underlineText, "test text") {
		t.Error("Underline function should contain the original text")
	}

	// Test empty string handling
	emptyRed := color.Red("")
	if emptyRed != "\033[31m\033[0m" {
		t.Error("Color functions should handle empty strings properly")
	}

	// Test special characters
	specialText := "!@#$%^&*()"
	specialRed := color.Red(specialText)
	if !strings.Contains(specialRed, specialText) {
		t.Error("Color functions should handle special characters properly")
	}

	fmt.Println("Color functions test passed!")
}

func TestColorReset(t *testing.T) {
	fmt.Println("=== Test color reset functionality ===")

	color := &utils.Color{}

	// Test that all functions end with reset code
	testCases := []struct {
		name     string
		function func(string) string
	}{
		{"Red", color.Red},
		{"Green", color.Green},
		{"Blue", color.Blue},
		{"Underline", color.Underline},
	}

	for _, tc := range testCases {
		result := tc.function("test")
		if !strings.HasSuffix(result, "\033[0m") {
			t.Errorf("%s function should end with reset code", tc.name)
		}
	}

	fmt.Println("Color reset functionality test passed!")
}
