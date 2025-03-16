package cmd

import (
	"fmt"

	"github.com/pablobfonseca/codemate/internal"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat [message]",
	Short: "Chat with codemate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]

		err := internal.InitDB()
		if err != nil {
			fmt.Errorf("Error initiating DB: %v\n", err)
		}

		context, err := internal.GetProjectContext()
		if err != nil {
			fmt.Println("Error loading project context:", err)
			return
		}

		internal.SendMessage(message, context)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
