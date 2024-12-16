package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "attorney-toolkit",
	Short: "Attorney is a lightweight forward proxy written in Go",
	Long:  `Allows HTTP/HTTPS transparent proxying on your machine.`,
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
