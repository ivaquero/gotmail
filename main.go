package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ivaquero/gotmail/utils"
)

func main() {
	// 获取当前执行目录
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}
	execDir := filepath.Dir(execPath)

	// 设置数据文件路径
	dataPath := filepath.Join(execDir, "data", "account.json")

	// 创建邮件管理器
	mailManager := utils.NewMailManager(dataPath)

	// 如果没有参数，显示帮助
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		if err := mailManager.CreateAccount(); err != nil {
			log.Fatal("Error creating account:", err)
		}

	case "messages":
		messages, err := mailManager.FetchMessages()
		if err != nil {
			log.Fatal("Error fetching messages:", err)
		}
		if messages != nil {
			fmt.Println("Messages:")
			for i, msg := range messages {
				fmt.Printf("%d. ID: %s, From: %s\n", i+1, msg.ID, msg.From)
			}
		}

	case "delete":
		if err := mailManager.DeleteAccount(); err != nil {
			log.Fatal("Error deleting account:", err)
		}

	case "details":
		if err := mailManager.ShowDetails(); err != nil {
			log.Fatal("Error showing details:", err)
		}

	case "open":
		if len(os.Args) < 3 {
			fmt.Println("Please provide email number to open")
			return
		}
		var emailNum int
		if _, err := fmt.Sscanf(os.Args[2], "%d", &emailNum); err != nil {
			log.Fatal("Invalid email number:", err)
		}
		if err := mailManager.OpenEmail(emailNum); err != nil {
			log.Fatal("Error opening email:", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Mail.tm CLI Tool - Go Version")
	fmt.Println("\nUsage:")
	fmt.Println("  create          - Create a new account")
	fmt.Println("  messages        - Fetch and list messages")
	fmt.Println("  delete          - Delete the account")
	fmt.Println("  details         - Show account details")
	fmt.Println("  open <number>   - Open specific email in browser")
	fmt.Println("\nExamples:")
	fmt.Println("gotmail create")
	fmt.Println("gotmail messages")
	fmt.Println("gotmail open 1")
}
