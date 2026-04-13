package utils_test

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/ivaquero/gotmail/utils"
)

func TestCopy(t *testing.T) {
	// 检查当前平台是否支持剪贴板操作
	if !isClipboardSupported() {
		t.Skip("剪贴板功能在当前平台不支持，跳过测试")
		return
	}

	// 测试多种数据类型
	testCases := []struct {
		name string
		data string
	}{
		{"英文文本", "Hello, World!"},
		{"中文文本", "这是测试数据。"},
		{"混合文本", "Hello, World! 这是测试数据。"},
		{"特殊字符", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
		{"多行文本", "第一行\n第二行\n第三行"},
	}

	fmt.Println("=== 剪贴板复制功能测试 ===")

	for i, tc := range testCases {
		fmt.Printf("测试 %d: %s\n", i+1, tc.name)
		fmt.Printf("数据: %s\n", tc.data)
		fmt.Print("正在复制到剪贴板...\n")

		if err := utils.Copy(tc.data); err != nil {
			log.Printf("复制失败: %v\n", err)
			t.Errorf("测试 %s 失败: %v\n", tc.name, err)
			continue
		}

		fmt.Println(" 成功！")
	}

	fmt.Println("=== 所有测试完成 ===")
	fmt.Println("提示: 您可以手动粘贴验证剪贴板内容")
}

// isClipboardSupported 检查当前平台是否支持剪贴板操作
func isClipboardSupported() bool {
	switch runtime.GOOS {
	case "windows", "darwin":
		return true
	default:
		return false
	}
}
