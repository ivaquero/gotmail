package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Copy copies data to the system clipboard
// Supports Windows (clip), macOS (pbcopy), and Linux (xclip/xsel)
func Copy(data string) error {
	platform := runtime.GOOS

	switch platform {
	case "windows":
		return copyToClipboard("clip", data)
	case "darwin":
		return copyToClipboard("pbcopy", data)
	case "linux":
		// Try xclip first, if not available try xsel
		if err := copyToClipboard("xclip", data); err == nil {
			return nil
		}
		if err := copyToClipboard("xsel", data); err == nil {
			return nil
		}
		return fmt.Errorf("no clipboard utility found (please install xclip or xsel)")
	default:
		return fmt.Errorf("platform not supported: %s", platform)
	}
}

// copyToClipboard performs the actual copy operation
func copyToClipboard(command string, data string) error {
	cmd := exec.Command(command)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", command, err)
	}

	if _, err := stdin.Write([]byte(data)); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to write data: %w", err)
	}

	if err := stdin.Close(); err != nil {
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s command failed: %w", command, err)
	}

	return nil
}
