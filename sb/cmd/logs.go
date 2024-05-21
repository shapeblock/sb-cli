package cmd

import (
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	appsCmd.AddCommand(logsCmd)
}
