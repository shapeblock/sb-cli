package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var logsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "A brief description of your command",
	Aliases: []string{"log"},
	Run:     appLogs,
}

var tail bool

func appLogs(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

	// API call
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
		return
	}

	/*token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}*/
token, err := GetToken(sbUrl)
if err != nil {
    fmt.Printf("error getting token: %v\n", err)
    return
}


	fullUrl := fmt.Sprintf("%s/apps/%s/shell-info/", sbUrl, app.UUID)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	var shellInfo ShellInfo
	if err := json.NewDecoder(resp.Body).Decode(&shellInfo); err != nil {
		fmt.Printf("Unable to decode podinfo from response: %v\n", err)
		os.Exit(1)
	}

	decodedKubeConfig, err := base64.StdEncoding.DecodeString(shellInfo.KubeConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode kubeconfig: %v\n", err)
		os.Exit(1)
	}

	label := fmt.Sprintf("appUuid=%s", app.UUID)

	if err = streamLogsFromPods(label, shellInfo.Namespace, string(decodedKubeConfig), tail); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stream logs: %v\n", err)
		os.Exit(1)
	}

}

func streamLogsFromPods(label, namespace, kubeconfig string, stream bool) error {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		if stream {
			fmt.Printf("Streaming logs for pod %s\n", pod.Name)
			err := streamPodLogs(clientset, pod, namespace)
			if err != nil {
				fmt.Printf("Error streaming logs for pod %s: %v\n", pod.Name, err)
			}
		} else {
			fmt.Printf("Getting last 100 lines of logs for pod %s\n", pod.Name)
			err := getLastPodLogs(clientset, pod, namespace)
			if err != nil {
				fmt.Printf("Error getting logs for pod %s: %v\n", pod.Name, err)
			}
		}
	}

	return nil
}

func streamPodLogs(clientset *kubernetes.Clientset, pod v1.Pod, namespace string) error {
	req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		Follow: true,
	})

	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer podLogs.Close()

	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func getLastPodLogs(clientset *kubernetes.Clientset, pod v1.Pod, namespace string) error {
	req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		TailLines: int64Ptr(100),
	})

	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer podLogs.Close()

	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func int64Ptr(i int64) *int64 {
	return &i
}
func init() {
	appsCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolVarP(&tail, "follow", "f", false, "Follow the pod logs")
}
