package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"utils/index"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func main() {
	var version string

	file, err := os.Open("./package.json")
	if err != nil {
		fmt.Println("Error reading package.json:", err)
		return
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &version) // Assuming version is a string in package.json

	var rootCmd = &cobra.Command{
		Use:   "tmail",
		Short: "⚡️ Quickly generate a disposable email straight from terminal.",
	}

	rootCmd.Version = version

	// Generate a new email
	var generateCmd = &cobra.Command{
		Use:   "g",
		Short: "Generate a new email",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CreateAccount()
		},
	}

	// Fetch messages from the inbox
	var fetchMessagesCmd = &cobra.Command{
		Use:   "m",
		Short: "Fetch messages from the inbox",
		Run: func(cmd *cobra.Command, args []string) {
			emails, err := utils.FetchMessages()
			if err != nil {
				fmt.Println("Error fetching messages:", err)
				return
			}

			if len(emails) == 0 {
				return
			}

			// Show the emails using promptui
			var choices []string
			for index, email := range emails {
				choices = append(choices, fmt.Sprintf("%d. %s - %s: %s", index+1, color.BlueString(email.Subject), color.YellowString("From:"), email.From.Address))
			}

			prompt := promptui.Select{
				Label: "Select an email",
				Items: choices,
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Println("Prompt failed:", err)
				return
			}

			// Open the email
			utils.OpenEmail(index + 1)
		},
	}

	// Delete account
	var deleteCmd = &cobra.Command{
		Use:   "d",
		Short: "Delete account",
		Run: func(cmd *cobra.Command, args []string) {
			utils.DeleteAccount()
		},
	}

	// Show details of the account
	var detailsCmd = &cobra.Command{
		Use:   "me",
		Short: "Show details of the account",
		Run: func(cmd *cobra.Command, args []string) {
			utils.ShowDetails()
		},
	}

	rootCmd.AddCommand(generateCmd, fetchMessagesCmd, deleteCmd, detailsCmd)
	rootCmd.Execute()
}
