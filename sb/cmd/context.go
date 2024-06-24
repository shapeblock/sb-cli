package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)
var contextCmd = &cobra.Command{
	Use:   "context",
	Aliases: []string{"con"}, 
	Short: "Manage Context",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like set or list or get.")
	},
}


func init(){
	rootCmd.AddCommand(contextCmd)
}