package cmd

import (
	"fmt"

	"github.com/pablobfonseca/codemate/internal"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat [message]",
	Short: "Chat with codemate",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.InitDB()
		if err != nil {
			fmt.Errorf("Error initiating DB: %v\n", err)
			return
		}

		internal.RunChatUI()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
