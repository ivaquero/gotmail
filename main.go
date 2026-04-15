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

	case "help":
		if len(os.Args) > 2 {
			showCommandHelp(os.Args[2])
		} else {
			showHelp()
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Your Temporary Email Accounts Manager")
	fmt.Println("\nUsage:")
	fmt.Println("  new                           Create a new account")
	fmt.Println("  ls                            List all accounts")
	fmt.Println("  msg [id <id>]                 Fetch and list messages")
	fmt.Println("  del [id <id>]                 Delete account")
	fmt.Println("  show [id <id>]                Show account details")
	fmt.Println("  open <number> [id <id>]       Open specific email in browser")
	fmt.Println("  export <path> [id <id>]       Export account data to specified path")
	fmt.Println("  help                          Show this help message")
	fmt.Println("\nOptions:")
	fmt.Println("  id <id>                       Specify account ID for operations")
}

func showCommandHelp(command string) {
	switch command {
	case "new":
		fmt.Println("Create a new account")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail new")
		fmt.Println("\nDescription:")
		fmt.Println("  Creates a new temporary email account with a random address")
		fmt.Println("  The account credentials will be stored locally for future use")

	case "ls":
		fmt.Println("List all accounts")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail ls")
		fmt.Println("\nDescription:")
		fmt.Println("  Displays all stored email accounts with their IDs and addresses")

	case "msg":
		fmt.Println("Fetch and list messages")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail msg [id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Retrieves messages from the specified account or the default account")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail msg                    # Fetch from default account")
		fmt.Println("  gotmail msg id abc123          # Fetch from specific account")

	case "del":
		fmt.Println("Delete account")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail del [id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Removes the specified account or the default account from storage")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail del                    # Delete default account")
		fmt.Println("  gotmail del id abc123          # Delete specific account")

	case "show":
		fmt.Println("Show account details")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail show [id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Displays detailed information about the specified account")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail show                   # Show default account")
		fmt.Println("  gotmail show id abc123         # Show specific account")

	case "open":
		fmt.Println("Open specific email in browser")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail open <number> [id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Opens the specified email message in your default web browser")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail open 1                 # Open first message from default account")
		fmt.Println("  gotmail open 3 id abc123       # Open third message from specific account")

	case "export":
		fmt.Println("Export account data to specified path")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail export <path> [id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Exports account data to the specified file path")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail export /tmp/backup.json          # Export default account")
		fmt.Println("  gotmail export /tmp/backup.json id abc123 # Export specific account")

	case "help":
		fmt.Println("Show help information")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail help [command]")
		fmt.Println("\nDescription:")
		fmt.Println("  Shows general help or detailed help for a specific command")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail help                   # Show general help")
		fmt.Println("  gotmail help msg               # Show help for msg command")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("\nAvailable commands:")
		fmt.Println("  new, ls, msg, del, show, open, export, help")
		fmt.Println("\nUse 'gotmail help <command>' for detailed help on a specific command")
	}
}
