import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 添加版本信息
	rootCmd.Flags().BoolP("version", "v", false, "Output the current version")
	rootCmd.Version(version)

	// 添加子命令
	addCommands()
}

func addCommands() {
	var cmd = &cobra.Command{
		Use:   "g",
		Short: "Generate a new email address.",
		Long:  `Generates a new temporary email address.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := utils.CreateAccount()
			if err != nil {
				fmt.Printf("Error generating email: %v\n", err)
			}
		},
	}

	cmd2 := &cobra.Command{
		Use:   "m",
		Short: "Fetch messages from inbox.",
		Long:  `Fetches and displays messages from the inbox.`,
		Run: func(cmd *cobra.Command, args []string) {
			emails, err := utils.FetchMessages()
			if err != nil {
				fmt.Printf("Error fetching messages: %v\n", err)
				return
			}

			if len(emails) == 0 {
				fmt.Println("No new messages.")
				return
			}

			// 在此处实现选择邮件的逻辑，例如打印列表并读取用户输入
			var selected int64
			for i, email := range emails {
				fmt.Printf("%d. %s - From: %s\n", i+1, email.Subject, email.From.Address)
			}
			fmt.Printf("\nEnter the index of the email to open (0 to cancel): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan(&selected)

			if selected <= 0 || int(selected) >= len(emails) {
				return
			}

			err = utils.OpenEmail(emails[int(selected)-1])
			if err != nil {
				fmt.Printf("Error opening email: %v\n", err)
			}
		},
	}

	rootCmd.AddCommand(cmd, cmd2)
}
