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
	dataPath := filepath.Join(execDir, "accounts.json")

	// Create mail manager
	mailManager := utils.NewMailManager(dataPath)

	// If no arguments, show help
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	// Parse account ID from arguments
	accountID, hasAccountID := utils.ParseAccountID(os.Args)

	switch command {
	case "new":
		if err := mailManager.CreateAccount(); err != nil {
			log.Fatal("Error creating account:", err)
		}

	case "list":
		if err := mailManager.ListAccounts(); err != nil {
			log.Fatal("Error listing accounts:", err)
		}

	case "msg":
		var messages []utils.Message
		var err error

		if hasAccountID {
			messages, err = mailManager.FetchMessagesByAccountID(accountID)
		} else {
			messages, err = mailManager.FetchMessages()
		}

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
		if hasAccountID {
			if err := mailManager.DeleteAccountByID(accountID); err != nil {
				log.Fatal("Error deleting account:", err)
			}
		} else {
			if err := mailManager.DeleteAccount(); err != nil {
				log.Fatal("Error deleting account:", err)
			}
		}

	case "show":
		if hasAccountID {
			if err := mailManager.ShowAccountDetails(accountID); err != nil {
				log.Fatal("Error showing details:", err)
			}
		} else {
			// Show all accounts if no ID specified
			if err := mailManager.ListAccounts(); err != nil {
				log.Fatal("Error showing accounts:", err)
			}
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
		if hasAccountID {
			if err := mailManager.OpenEmailByAccountID(accountID, emailNum); err != nil {
				log.Fatal("Error opening email:", err)
			}
		} else {
			if err := mailManager.OpenEmail(emailNum); err != nil {
				log.Fatal("Error opening email:", err)
			}
		}

	case "export":
		if len(os.Args) < 3 {
			fmt.Println("Please provide export path")
			return
		}
		exportPath := os.Args[2]
		if hasAccountID {
			if err := mailManager.ExportAccountByID(accountID, exportPath); err != nil {
				log.Fatal("Error exporting account:", err)
			}
		} else {
			if err := mailManager.ExportAccount(exportPath); err != nil {
				log.Fatal("Error exporting account:", err)
			}
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Mail.tm CLI Tool - Go Version")
	fmt.Println("\nUsage:")
	fmt.Println("  new                           Create a new account")
	fmt.Println("  list                          List all accounts")
	fmt.Println("  msg [id <id>]                 Fetch and list messages")
	fmt.Println("  del [id <id>]                 Delete account")
	fmt.Println("  show [id <id>]                Show account details")
	fmt.Println("  open <number> [id <id>]       Open specific email in browser")
	fmt.Println("  export <path> [id <id>]       Export account data to specified path")
	fmt.Println("\nOptions:")
	fmt.Println("  id <id>                       Specify account ID for operations")
	fmt.Println("\nExamples:")
	fmt.Println("  gotmail new")
	fmt.Println("  gotmail list")
	fmt.Println("  gotmail msg")
	fmt.Println("  gotmail msg id abc123")
	fmt.Println("  gotmail show")
	fmt.Println("  gotmail show id abc123")
	fmt.Println("  gotmail open 1")
	fmt.Println("  gotmail open 1 id abc123")
	fmt.Println("  gotmail export /path/to/backup")
	fmt.Println("  gotmail export /path/to/backup id abc123")
}
