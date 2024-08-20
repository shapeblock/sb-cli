package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout current user",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("endpoint", "")
		viper.Set("token", "")
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
