package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ivaquero/gotmail/utils"
)

func main() {
	// Get current executable directory
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}
	execDir := filepath.Dir(execPath)

	// Set data file path
	dataPath := filepath.Join(execDir, "data", "account.json")

	// Create mail manager
	mailManager := utils.NewMailManager(dataPath)

	// If no arguments, show help
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "new":
		if err := mailManager.CreateAccount(); err != nil {
			log.Fatal("Error creating account:", err)
		}

	case "msg":
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

	case "del":
		if err := mailManager.DeleteAccount(); err != nil {
			log.Fatal("Error deleting account:", err)
		}

	case "show":
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

	case "export":
		if len(os.Args) < 3 {
			fmt.Println("Please provide export path")
			return
		}
		exportPath := os.Args[2]
		if err := mailManager.ExportAccount(exportPath); err != nil {
			log.Fatal("Error exporting account:", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Mail.tm CLI Tool - Go Version")
	fmt.Println("\nUsage:")
	fmt.Println("  new             - Create a new account")
	fmt.Println("  msg             - Fetch and list messages")
	fmt.Println("  del             - Delete the account")
	fmt.Println("  show            - Show account details")
	fmt.Println("  open <number>   - Open specific email in browser")
	fmt.Println("  export <path>   - Export account data to specified path")
	fmt.Println("\nExamples:")
	fmt.Println("gotmail new")
	fmt.Println("gotmail msg")
	fmt.Println("gotmail show")
	fmt.Println("gotmail open 1")
	fmt.Println("gotmail export /path/to/account.json")
}
