package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v3"
)

var basePath string

func init() {
	rand.Seed(time.Now().UnixNano())
	homeDir, _ := homedir.Dir()
	basePath = filepath.Join(homeDir, ".config", "app")
}

type Account struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

var accountData *Account

func createAccount() {
	// start the spinner
	sp := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithPrefix("Creating..."))
	sp.Start()

	// read the account data from file
	readAccountData()

	if accountData != nil {
		sp.Stop()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Account already exists"))
		return
	}

	// get the available email domains
	resp, err := http.Get("https://api.mail.tm/domains?page=1")
	if err != nil {
		sp.Stop()
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var domainsResponse struct {
		HydraMember []struct {
			Domain string `json:"domain"`
		} `json:"hydra:member"`
	}
	json.NewDecoder(resp.Body).Decode(&domainsResponse)

	// get the first domain
	domain := domainsResponse.HydraMember[0].Domain

	// generate a random email address
	email := fmt.Sprintf("%s@%s", randString(7), domain)

	// generate a random password
	password := randString(7)

	// create account
	account := Account{Address: email, Password: password}
	data, _ := json.Marshal(account)

	resp, err = http.Post("https://api.mail.tm/accounts", "application/json", bytes.NewBuffer(data))
	if err != nil {
		sp.Stop()
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&account)

	// copy the email to the clipboard
	// Clipboard copying is skipped for simplicity

	// get Jwt token
	tokenData := map[string]string{"address": email, "password": password}
	tokenJSON, _ := json.Marshal(tokenData)

	resp, err = http.Post("https://api.mail.tm/token", "application/json", bytes.NewBuffer(tokenJSON))
	if err != nil {
		sp.Stop()
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&tokenResponse)

	account.Token = tokenResponse.Token

	// write data to a JSON file
	writeAccountData(&account)

	// stop the spinner
	sp.Stop()
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render("Account created: " + lipgloss.NewStyle().Underline(true).Render(email)))
}

func fetchMessages() []interface{} {
	// start the spinner
	sp := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithPrefix("Fetching..."))
	sp.Start()

	readAccountData()

	if accountData == nil {
		sp.Stop()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Account not created yet"))
		return nil
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.mail.tm/messages", nil)
	req.Header.Set("Authorization", "Bearer "+accountData.Token)
	resp, err := client.Do(req)
	if err != nil {
		sp.Stop()
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	var messagesResponse struct {
		HydraMember []interface{} `json:"hydra:member"`
	}
	json.NewDecoder(resp.Body).Decode(&messagesResponse)

	// stop the spinner
	sp.Stop()

	if len(messagesResponse.HydraMember) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("No Emails"))
		return nil
	}
	return messagesResponse.HydraMember
}

func deleteAccount() {
	// start the spinner
	sp := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithPrefix("Deleting..."))
	sp.Start()

	readAccountData()

	if accountData == nil {
		sp.Stop()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Account not created yet"))
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("https://api.mail.tm/accounts/%s", accountData.ID), nil)
	req.Header.Set("Authorization", "Bearer "+accountData.Token)
	resp, err := client.Do(req)
	if err != nil {
		sp.Stop()
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// delete the account.json file
	os.Remove(filepath.Join(basePath, "account.json"))

	// stop the spinner
	sp.Stop()
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render("Account deleted"))
}

func showDetails() {
	// start the spinner
	sp := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithPrefix("Fetching details..."))
	sp.Start()

	readAccountData()

	if accountData == nil {
		sp.Stop()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Account not created yet"))
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.mail.tm/accounts/%s", accountData.ID), nil)
	req.Header.Set("Authorization", "Bearer "+accountData.Token)
	resp, err := client.Do(req)
	if err != nil {
		sp.Stop()
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var detailsResponse Account
	json.NewDecoder(resp.Body).Decode(&detailsResponse)

	// stop the spinner
	sp.Stop()
	fmt.Printf(`
	Email: %s
	CreatedAt: %s
	`, lipgloss.NewStyle().Underline(true).Render(detailsResponse.Address), lipgloss.NewStyle().Render(time.Unix(detailsResponse.CreatedAt, 0).Format(time.RFC1123)))
}

func openEmail(emailIndex int) {
	emailFilePath := filepath.Join(basePath, "email.html")
	sp := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithPrefix("Opening..."))
	sp.Start()

	readAccountData()
	fetchedEmails := fetchMessages()
	if fetchedEmails == nil {
		sp.Stop()
		return
	}
	mailToOpen := fetchedEmails[emailIndex-1]

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.mail.tm/messages/%s", mailToOpen.ID), nil)
	req.Header.Set("Authorization", "Bearer "+accountData.Token)
	resp, err := client.Do(req)
	if err != nil {
		sp.Stop()
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var emailResponse struct {
		HTML []string `json:"html"`
	}
	json.NewDecoder(resp.Body).Decode(&emailResponse)

	if len(emailResponse.HTML) == 0 {
		sp.Stop()
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("No HTML content found"))
		return
	}

	// write the email html content to a file
	os.WriteFile(emailFilePath, []byte(emailResponse.HTML[0]), 0644)

	// open the email html file in the browser
	// Opening the file is skipped for simplicity

	// stop the spinner
	sp.Stop()
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func readAccountData() {
	data, err := os.ReadFile(filepath.Join(basePath, "account.json"))
	if err != nil {
		accountData = nil
		return
	}
	json.Unmarshal(data, &accountData)
}

func writeAccountData(account *Account) {
	data, _ := json.Marshal(account)
	os.WriteFile(filepath.Join(basePath, "account.json"), data, 0644)
}

func main() {
	app := &cli.App{
		Name:  "Account Manager",
		Usage: "Manage email accounts",
		Commands: []*cli.Command{
			{
				Name:   "create",
				Usage:  "Create a new account",
				Action: func(c *cli.Context) error { createAccount(); return nil },
			},
			{
				Name:   "fetch",
				Usage:  "Fetch messages",
				Action: func(c *cli.Context) error { fetchMessages(); return nil },
			},
			{
				Name:   "delete",
				Usage:  "Delete the account",
				Action: func(c *cli.Context) error { deleteAccount(); return nil },
			},
			{
				Name:   "details",
				Usage:  "Show account details",
				Action: func(c *cli.Context) error { showDetails(); return nil },
			},
			{
				Name:   "open",
				Usage:  "Open a specific email",
				Action: func(c *cli.Context) error { openEmail(c.Int("index")); return nil },
			},
		},
	}

	app.Run(os.Args)
}
