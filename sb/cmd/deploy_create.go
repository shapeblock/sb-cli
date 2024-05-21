package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var createDeployCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deployment",
	Run:   createDeployment,
}

type PodInfo struct {
	Name       string `json:"name"`
	KubeConfig string `json:"kubeconfig"`
	Namespace  string `json:"namespace"`
}

type DeploymentResponse struct {
	UUID string `json:"uuid"`
}

var follow bool

func createDeployment(cmd *cobra.Command, args []string) {
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

	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}

	fullUrl := fmt.Sprintf("%s/api/apps/%s/deployments/", sbUrl, app.UUID)

	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Deployment created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
		os.Exit(1)
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create deployment, bad request.")
		os.Exit(1)
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create deployment, internal server error.")
		os.Exit(1)
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var deploymentResponse DeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&deploymentResponse); err != nil {
		fmt.Printf("Unable to decode deployment response: %v\n", err)
		os.Exit(1)
	}
	// TODO: if follow flag is given, stream deployment logs
	if follow {
		time.Sleep(5 * time.Second)
		fullUrl = fmt.Sprintf("%s/deployments/%s/pod-info/", sbUrl, deploymentResponse.UUID)

		req, err = http.NewRequest("GET", fullUrl, nil)
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

		client = &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		var podInfo PodInfo
		if err := json.NewDecoder(resp.Body).Decode(&podInfo); err != nil {
			fmt.Printf("Unable to decode podinfo from response: %v\n", err)
			os.Exit(1)
		}

		decodedKubeConfig, err := base64.StdEncoding.DecodeString(podInfo.KubeConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode kubeconfig: %v\n", err)
			os.Exit(1)
		}

		if err := tailPodLogs(podInfo.Name, string(decodedKubeConfig), podInfo.Namespace); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to tail logs: %v\n", err)
			os.Exit(1)
		}

	}

}

func tailPodLogs(podName, kubeConfig, namespace string) error {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for _, container := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
		if err := streamLogsWithRetry(clientset, namespace, podName, container.Name); err != nil {
			return err
		}
	}

	for {
		pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if pod.Status.Phase == corev1.PodSucceeded {
			fmt.Println("Pod has completed.")
			break
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}

func streamLogsWithRetry(clientset *kubernetes.Clientset, namespace, podName, containerName string) error {
	for {
		err := streamLogs(clientset, namespace, podName, containerName)
		if err == nil {
			break
		}
		fmt.Printf("Retrying log stream for container %s in pod %s: %v\n", containerName, podName, err)
		time.Sleep(2 * time.Second)
	}
	return nil
}

func streamLogs(clientset *kubernetes.Clientset, namespace, podName, containerName string) error {
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{Container: containerName, Follow: true})
	stream, err := req.Stream(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()

	_, err = io.Copy(os.Stdout, stream)
	return err
}

func init() {
	deployCmd.AddCommand(createDeployCmd)
	createDeployCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow the pod logs")
}
