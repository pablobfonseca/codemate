package cmd

import (
	"fmt"
	"os"

	"github.com/pablobfonseca/codemate/internal"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive chat with Codemate",
	Long: `Start an interactive chat session with Codemate, your AI code assistant.
	
Codemate will analyze your project files and provide helpful answers to your coding questions.
It uses a local large language model through Ollama to generate responses.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.InitDB()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}

		internal.RunChatUI()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
