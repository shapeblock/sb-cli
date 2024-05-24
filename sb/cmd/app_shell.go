package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type ShellInfo struct {
	Name       string `json:"name"`
	KubeConfig string `json:"kubeconfig"`
	Namespace  string `json:"namespace"`
}

var shellCmd = &cobra.Command{
	Use:     "shell",
	Short:   "App shell",
	Aliases: []string{"sh"},
	Run:     appShell,
}

func appShell(cmd *cobra.Command, args []string) {
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

	fullUrl := fmt.Sprintf("%s/apps/%s/shell-info/", sbUrl, app.UUID)

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

	if err = execIntoPod(shellInfo.Name, shellInfo.Namespace, string(decodedKubeConfig)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to exec into pod: %v\n", err)
		os.Exit(1)
	}

}

func execIntoPod(podName, namespace, kubeConfig string) error {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// TODO: send EOF signal correctly
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", pod.Spec.Containers[0].Name).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true").
		Param("command", "launcher").
		Param("command", "bash")

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return err
	}

	return nil
}
func init() {
	appsCmd.AddCommand(shellCmd)
}
