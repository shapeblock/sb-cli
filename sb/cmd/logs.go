/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/shapeblock/sb-cli/sb/config"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var ClusterId string

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logs called")
		kubeconfigFile := path.Join(config.GetConfigDir(), fmt.Sprintf("%s.yaml", ClusterId))
		config, _ := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
		// creates the clientset
		clientset, _ := kubernetes.NewForConfig(config)
		count := int64(100)
		podLogOptions := corev1.PodLogOptions{
			Follow:    true,
			TailLines: &count,
		}
		podLogRequest := clientset.CoreV1().
			Pods("rain").
			GetLogs("rain-rain-drupal-7c8cb87b8-gfv8t", &podLogOptions)
		stream, err := podLogRequest.Stream(context.TODO())
		if err != nil {
			fmt.Println("Unable to stream logs")
		}
		defer stream.Close()

		for {
			buf := make([]byte, 2000)
			numBytes, err := stream.Read(buf)
			if numBytes == 0 {
				continue
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Unable to print logs to console")
			}
			message := string(buf[:numBytes])
			fmt.Print(message)
		}

	},
}

func init() {
	appsCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringVarP(&ClusterId, "cluster", "c", "", "The cluster ID for which the projects need to be listed")
	logsCmd.MarkFlagRequired("cluster")
}
