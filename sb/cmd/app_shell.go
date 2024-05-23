package cmd;

import (
	//"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	//"io/ioutil"
	"os"
	//"context"
	//"time"
	v1 "k8s.io/api/core/v1" // Import for PodExecOptions
	"k8s.io/client-go/tools/remotecommand" // Import for NewSPDYExecutor and Stream

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var displayShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Shell access",
	Run:   shellAccess,
}

type ShellInfo struct {
	KubeConfig string `json:"kubeconfig"`
	Namespace  string `json:"namespace"`
	//Label string `json:"label"`
}
//var access bool

func shellAccess(cmd *cobra.Command, args []string) {

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
    fullUrl:= fmt.Sprintf("%s/apps/%s/shell-info/", sbUrl, app.UUID)

		req, err := http.NewRequest("GET", fullUrl, nil)
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


		var shellInfo ShellInfo
		if err := json.NewDecoder(resp.Body).Decode(&shellInfo); err != nil {
			fmt.Printf("Unable to decode podinfo from response: %v\n", err)
			os.Exit(1)
		}
		//fmt.Printf("Kubeconfig: %s\n", podInfo.KubeConfig)
	   //fmt.Printf("Namespace: %s\n", podInfo.Namespace)
       // fmt.Printf("Name:%s\n",shellInfo.Name)
		decodedKubeConfig, err := base64.StdEncoding.DecodeString(shellInfo.KubeConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode kubeconfig: %v\n", err)
			os.Exit(1)
		}
        //fmt.Printf("decoded: %s\n",decodedKubeConfig)
		if err := podAccess(string(decodedKubeConfig), shellInfo.Namespace); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to tail logs: %v\n", err)
			os.Exit(1)
		}
	
}

func podAccess(kubeconfig,namespace string) error{

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %s\n", err.Error())
		os.Exit(1)
	}

	podName := "final-test-python-659b486586-mdbqs"

	execCommand := []string{"sh"} // Command to execute, e.g., "sh" for shell access

	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: execCommand,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, metav1.ParameterCodec)
//fmt.Println(namespace)
//fmt.Println(podName)
//fmt.Println(kubeconfig)
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Printf("Error creating Executor: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Starting stream to pod")
	fmt.Printf("Pod: %s\n", podName)
	fmt.Printf("Namespace: %s\n", namespace)
	fmt.Printf("Command: %v\n", execCommand)
	fmt.Printf("Request URL: %s\n", req.URL())

	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		fmt.Printf("Error executing command in container: %s\n", err.Error())
		os.Exit(1)
	}
	return nil
}

func init() {
	appsCmd.AddCommand(displayShellCmd)
}
