package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ivaquero/gotmail/utils"
)

// validateAccountID validates account ID format
func validateAccountID(accountID string) error {
	if accountID == "" {
		return fmt.Errorf("account ID cannot be empty")
	}
	// Account ID should be alphanumeric with length between 10-50 characters
	if len(accountID) < 10 || len(accountID) > 50 {
		return fmt.Errorf("account ID length must be between 10 and 50 characters")
	}
	// Simple alphanumeric validation
	matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", accountID)
	if err != nil {
		return fmt.Errorf("failed to validate account ID format: %w", err)
	}
	if !matched {
		return fmt.Errorf("account ID can only contain letters and numbers")
	}
	return nil
}

// validateEmailIndex validates email index input
func validateEmailIndex(indexStr string) (int, error) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return 0, fmt.Errorf("email index must be a valid number")
	}
	if index < 1 {
		return 0, fmt.Errorf("email index must be greater than 0")
	}
	return index, nil
}

// validateExportPath validates export folder path
func validateExportPath(path string) error {
	if path == "" {
		return fmt.Errorf("export path cannot be empty")
	}
	// Check if path is absolute or relative
	if filepath.IsAbs(path) {
		// For absolute paths, check if parent directory exists
		parent := filepath.Dir(path)
		if _, err := os.Stat(parent); err != nil {
			return fmt.Errorf("parent directory does not exist: %s", parent)
		}
	}
	return nil
}

func main() {
	// Get current executable directory
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	// Set data file path to ~/.gotmail.json
	dataPath := filepath.Join(homeDir, ".gotmail.json")

	// Create mail manager
	mailManager := utils.NewMailManager(dataPath)

	// If no arguments, show help
	if len(os.Args) < 2 {
		utils.ShowHelp()
		fmt.Println("\nTip: Use 'gotmail help <command>' for detailed command information")
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
			fmt.Fprintf(os.Stderr, "Error creating account: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Account created successfully!")

	case "ls":
		if err := mailManager.ListAccounts(); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing accounts: %v\n", err)
			os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Error fetching messages: %v\n", err)
			os.Exit(1)
		}
		if messages != nil {
			if len(messages) == 0 {
				fmt.Println("No messages found")
			} else {
				fmt.Printf("Found %d message(s):\n", len(messages))
				for i, msg := range messages {
					fmt.Printf("  %d. ID: %s, From: %s\n", i+1, msg.ID, msg.From)
				}
			}
		}

	case "del":
		if hasAccountID {
			// Validate account ID format
			if err := validateAccountID(accountID); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid account ID: %v\n", err)
				os.Exit(1)
			}
			if err := mailManager.DeleteAccountByID(accountID); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting account: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := mailManager.DeleteAccount(); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting account: %v\n", err)
				os.Exit(1)
			}
		}

	case "show":
		if hasAccountID {
			// 指定 id 时，显示 accounts.json 中的 id 对应信息
			if err := mailManager.ShowAccountDetails(accountID); err != nil {
				fmt.Fprintf(os.Stderr, "Error showing account details: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 未指定 id 时，打印整个 accounts.json
			jsonData, err := mailManager.GetAllAccountsJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting accounts JSON: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(jsonData)
			fmt.Println("Account data displayed successfully!")
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
			fmt.Fprintf(os.Stderr, "Invalid email number: %v\n", err)
			fmt.Fprintf(os.Stderr, "💡 Usage: gotmail open <email-number> [--id <account-id>]\n")
			os.Exit(1)
		}

		if hasAccountID {
			// Validate account ID format
			if err := validateAccountID(accountID); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid account ID: %v\n", err)
				os.Exit(1)
			}
			if err := mailManager.OpenEmailByAccountID(accountID, emailNum); err != nil {
				fmt.Fprintf(os.Stderr, "Error opening email: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := mailManager.OpenEmail(emailNum); err != nil {
				fmt.Fprintf(os.Stderr, "Error opening email: %v\n", err)
				os.Exit(1)
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
			// Validate account ID format
			if err := validateAccountID(accountID); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid account ID: %v\n", err)
				os.Exit(1)
			}
			if err := mailManager.ExportAccountByID(accountID, exportFolder); err != nil {
				fmt.Fprintf(os.Stderr, "Error exporting account: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Validate export path
			if err := validateExportPath(exportFolder); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid export path: %v\n", err)
				os.Exit(1)
			}
			if err := mailManager.ExportAccount(exportFolder); err != nil {
				fmt.Fprintf(os.Stderr, "Error exporting account: %v\n", err)
				os.Exit(1)
			}
		}

	case "help":
		if len(os.Args) > 2 {
			utils.ShowCommandHelp(os.Args[2])
		} else {
			utils.ShowHelp()
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintf(os.Stderr, "Available commands: new, ls, msg, del, show, open, export, help\n")
		fmt.Fprintf(os.Stderr, "   Use 'gotmail help' for more information\n")
		utils.ShowHelp()
	}
}
