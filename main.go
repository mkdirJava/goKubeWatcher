// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	outputFile := createOutputFile()
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	podsWatcher, err := clientset.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{}) //.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for {
		select {
		case result := <-podsWatcher.ResultChan():
			resultString := getEventString(result)
			if resultString != nil {
				outputFile.WriteString(*resultString)
			}
		}
	}
}

func createOutputFile() *os.File {
	if home, err := os.UserHomeDir(); err != nil {
		panic(err)
	} else {
		if file, err := os.Create(filepath.Join(home, ".kubeWatcher.result")); err != nil {
			panic(err)
		} else {
			return file
		}
	}

}

func getEventString(event watch.Event) *string {
	object := event.Object
	if pod, ok := object.(*v1.Pod); ok {
		formatedString := fmt.Sprintf(
			"EVENT: %s,Type: Pod,Name: %s,created: %s,Deleted : %s \n", event.Type, pod.Name, pod.CreationTimestamp, pod.DeletionTimestamp)
		return &formatedString
	} else {
		return nil
	}
}
