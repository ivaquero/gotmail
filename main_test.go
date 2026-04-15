package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput captures stdout output for testing
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestMainFunction(t *testing.T) {
	// Test with no arguments (should show help)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"gotmail"}

	output := captureOutput(func() {
		main()
	})

	if !strings.Contains(output, "Your Temporary Email Accounts Manager") {
		t.Error("Expected help message when no arguments provided")
	}
}

func TestMainWithHelpCommand(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"gotmail", "help"}

	output := captureOutput(func() {
		main()
	})

	if !strings.Contains(output, "Your Temporary Email Accounts Manager") {
		t.Error("Expected help message when 'help' command provided")
	}
}

func TestMainWithInvalidCommand(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"gotmail", "invalidcommand"}

	output := captureOutput(func() {
		main()
	})

	if !strings.Contains(output, "Unknown command") {
		t.Error("Expected 'Unknown command' message when invalid command provided")
	}

	if !strings.Contains(output, "Your Temporary Email Accounts Manager") {
		t.Error("Expected help message when invalid command provided")
	}
}
