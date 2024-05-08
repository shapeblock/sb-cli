/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (

	"github.com/spf13/cobra"
)

// buildEnvCmd represents the buildEnv command
var buildCmd=&cobra.Command{
	Use: "build-env",
	Short: "Manage build env tasks",
   
   }

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildEnvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildEnvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
