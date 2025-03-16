package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codemate",
	Short: "An AI assistant for your codebase",
	Long: `Codemate is a CLI tool that helps you chat with an AI assistant about your code.
It gathers context from your project files and uses AI to provide helpful answers
to your coding questions.

- Automatically collects context from your git repository
- Streams responses in real-time with a beautiful TUI
- Respects .gitignore rules when scanning your codebase
- Remembers conversation history for contextual responses`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		fmt.Println("Welcome to Codemate! Use 'codemate chat' to start chatting with your AI assistant.")
		fmt.Println("\nFor more information, run 'codemate --help'")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codemate.yaml)")
}
