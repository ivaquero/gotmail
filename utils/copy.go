package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Copy 将数据复制到系统剪贴板
// 支持 Windows (clip) 和 macOS (pbcopy)
func Copy(data string) error {
	platform := runtime.GOOS

	switch platform {
	case "windows":
		return copyToClipboard("clip", data)
	case "darwin":
		return copyToClipboard("pbcopy", data)
	default:
		return fmt.Errorf("platform not supported: %s", platform)
	}
}

// copyToClipboard 执行具体的复制操作
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
