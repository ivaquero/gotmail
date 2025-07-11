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

	"github.com/manifoldco/promptui"
)

type Account struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Token    struct {
		Token string `json:"token"`
	} `json:"token"`
}

type Database struct {
	Data *Account `json:"data"`
}

var db Database
var dirname string

func init() {
	rand.Seed(time.Now().UnixNano())
	dirname, _ = os.UserHomeDir()
}

func createAccount() {
	spinner := promptui.Prompt{
		Label: "Creating...",
	}

	spinner.Run()

	filePath := filepath.Join(dirname, "../data/account.json")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Account already exists")
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&db); err == nil && db.Data != nil {
		spinner.Stop()
		fmt.Println("Account already exists")
		return
	}

	response, err := http.Get("https://api.mail.tm/domains?page=1")
	if err != nil {
		fmt.Println("Error fetching domains:", err)
		return
	}
	defer response.Body.Close()

	var domainResponse struct {
		HydraMember []struct {
			Domain string `json:"domain"`
		} `json:"hydra:member"`
	}
	json.NewDecoder(response.Body).Decode(&domainResponse)

	domain := domainResponse.HydraMember[0].Domain
	email := fmt.Sprintf("%s@%s", randomString(7), domain)
	password := randomString(7)

	accountData := Account{
		Address:  email,
		Password: password,
	}

	resp, err := http.Post("https://api.mail.tm/accounts", "application/json", bytes.NewBuffer(accountData))
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&accountData)
	accountData.Password = password

	db.Data = &accountData

	err = writeToFile(filePath, db)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	spinner.Stop()
	fmt.Printf("Account created: %s\n", email)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func writeToFile(filePath string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, file, 0644)
}

// Additional functions (fetchMessages, deleteAccount, showDetails, openEmail) would be defined similarly...

func main() {
	// Implement main logic and function calls here
}
