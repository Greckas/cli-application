package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check login status",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := loadToken()
		if err != nil {
			color.Red("User is not logged in.")
			return
		}

		// Check if the token is valid
		if token.Valid() {
			color.Green("User is logged in.")
		} else {
			color.Red("User is not logged in.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
