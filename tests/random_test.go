package utils_test

import (
	"fmt"
	"testing"

	"github.com/ivaquero/gotmail/utils"
)

func TestGenerateRandomString(t *testing.T) {
	fmt.Println("=== 测试随机字符串生成 ===")

	// 生成多个随机字符串来验证功能
	for i := 0; i < 5; i++ {
		randomStr := utils.GenerateRandomString(10)
		fmt.Printf("随机字符串 %d: %s\n", i+1, randomStr)

		// 验证字符串长度
		if len(randomStr) != 10 {
			t.Errorf("生成的字符串长度不正确: 期望 10, 实际 %d", len(randomStr))
		}

		// 验证字符串只包含允许的字符
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
				t.Errorf("生成的字符串包含非法字符: %c", char)
			}
		}
	}

	fmt.Println("\n=== 测试完成 ===")
}
