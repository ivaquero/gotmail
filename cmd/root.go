var rootCmd = &cobra.Command{
	Use:   "tmail",
	Short: "tmail is a CLI tool to manage temporary email.",
	Long:  `A CLI tool to generate and manage temporary email addresses.`, Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
