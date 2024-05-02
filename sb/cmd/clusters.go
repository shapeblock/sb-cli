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
	Short: "Manage clusters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

var scaleClusterCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a cluster up or down",
}

func init() {
	rootCmd.AddCommand(clustersCmd)
	clustersCmd.AddCommand(scaleClusterCmd)
}
