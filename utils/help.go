package utils

import (
	"fmt"
)

// ShowHelp displays general help information
func ShowHelp() {
	fmt.Println("Your Temporary Email Accounts Manager")
	fmt.Println("\nUsage:")
	fmt.Println("  new                           Create a new account")
	fmt.Println("  ls                            List all accounts")
	fmt.Println("  msg [--id <id>]               Fetch and list messages")
	fmt.Println("  del [--id <id>]               Delete account")
	fmt.Println("  show [--id <id>]              Show account details or all accounts in JSON")
	fmt.Println("  open <number> [--id <id>]     Open specific email in browser")
	fmt.Println("  export <folder> [--id <id>]   Export account data to specified folder")
	fmt.Println("  help                          Show this help message")
	fmt.Println("\nOptions:")
	fmt.Println("  --id <id>                     Specify account ID for operations")
}

// ShowCommandHelp displays detailed help for specific command
func ShowCommandHelp(command string) {
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
		fmt.Println("  gotmail msg [--id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Retrieves messages from the specified account or the default account")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail msg                    # Fetch from default account")
		fmt.Println("  gotmail msg --id abc123        # Fetch from specific account")

	case "del":
		fmt.Println("Delete account")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail del [--id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Removes the specified account or the default account from storage")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail del                    # Delete default account")
		fmt.Println("  gotmail del --id abc123        # Delete specific account")

	case "show":
		fmt.Println("Show account details or all accounts")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail show [--id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  When --id is specified, shows detailed information about the specific account")
		fmt.Println("  When --id is not specified, shows all accounts data in JSON format")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail show                   # Show all accounts in JSON format")
		fmt.Println("  gotmail show --id abc123       # Show specific account details")

	case "open":
		fmt.Println("Open specific email in browser")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail open <number> [--id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Opens the specified email message in your default web browser")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail open 1                 # Open first message from default account")
		fmt.Println("  gotmail open 3 --id abc123     # Open third message from specific account")

	case "export":
		fmt.Println("Export account data to specified folder")
		fmt.Println("\nUsage:")
		fmt.Println("  gotmail export <folder> [--id <account_id>]")
		fmt.Println("\nDescription:")
		fmt.Println("  Exports account data to the specified folder")
		fmt.Println("\nExamples:")
		fmt.Println("  gotmail export /tmp/backup          # Export default account")
		fmt.Println("  gotmail export /tmp/backup --id abc123 # Export specific account")

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
