/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clustersCmd represents the clusters command
var clustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

func init() {
	rootCmd.AddCommand(clustersCmd)
}
