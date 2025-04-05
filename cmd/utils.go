import (
	"fmt"
)

func CreateAccount() error {
	fmt.Printf("Creating new email address.\n")
	return nil // 实际应连接到服务创建地址
}

func FetchMessages() ([]Email, error) {
	fmt.Printf("Fetching messages from inbox.\n")
	// 这里模拟返回一些邮件数据
	return []Email{}, nil
}

type Email struct {
	Subject string
	From    FromInfo
}

type FromInfo struct {
	Address string
}
