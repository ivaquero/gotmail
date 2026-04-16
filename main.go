package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ivaquero/gotmail/utils"
)

func main() {
	// Get current executable directory
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get home directory:", err)
	}

	// Set data file path to ~/.gotmail.json
	dataPath := filepath.Join(homeDir, ".gotmail.json")

	// Create mail manager
	mailManager := utils.NewMailManager(dataPath)

	// If no arguments, show help
	if len(os.Args) < 2 {
		utils.ShowHelp()
		return
	}

	command := os.Args[1]

	// Parse flags
	var accountID string
	var hasAccountID bool

	// Create a new flag set for parsing command-specific flags
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.StringVar(&accountID, "id", "", "Account ID for operations")

	// Parse flags from os.Args[2:] (skip command name)
	if len(os.Args) > 2 {
		if err := fs.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing flags: %v\n", err)
			utils.ShowHelp()
			return
		}
	}

	if accountID != "" {
		hasAccountID = true
	}

	switch command {
	case "new":
		if err := mailManager.CreateAccount(); err != nil {
			log.Fatal("Error creating account: ", err)
		}

	case "ls":
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
			log.Fatal("\nError fetching messages: ", err)
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
				log.Fatal("Error deleting account: ", err)
			}
		} else {
			if err := mailManager.DeleteAccount(); err != nil {
				log.Fatal("Error deleting account: ", err)
			}
		}

	case "show":
		if hasAccountID {
			// 指定 id 时，显示 accounts.json 中的 id 对应信息
			if err := mailManager.ShowAccountDetails(accountID); err != nil {
				log.Fatal("Error showing details:", err)
			}
		} else {
			// 未指定 id 时，打印整个 accounts.json
			jsonData, err := mailManager.GetAllAccountsJSON()
			if err != nil {
				log.Fatal("Error getting accounts JSON:", err)
			}

			fmt.Println(jsonData)
		}

	case "open":
		// Parse remaining arguments for open command
		var emailNum int
		remainingArgs := fs.Args()

		if len(remainingArgs) < 1 {
			fmt.Println("Please provide email number to open")
			return
		}

		if _, err := fmt.Sscanf(remainingArgs[0], "%d", &emailNum); err != nil {
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
		remainingArgs := fs.Args()
		if len(remainingArgs) < 1 {
			fmt.Println("Please provide export folder")
			return
		}
		exportFolder := remainingArgs[0]
		if hasAccountID {
			if err := mailManager.ExportAccountByID(accountID, exportFolder); err != nil {
				log.Fatal("Error exporting account: ", err)
			}
		} else {
			if err := mailManager.ExportAccount(exportFolder); err != nil {
				log.Fatal("Error exporting account: ", err)
			}
		}

	case "help":
		if len(os.Args) > 2 {
			utils.ShowCommandHelp(os.Args[2])
		} else {
			utils.ShowHelp()
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		utils.ShowHelp()
	}
}
