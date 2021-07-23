package main

import (
	"context"
	"fmt"
	"path/filepath"

	workflowv1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	workflowclientset "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var home = homedir.HomeDir()
	var kubeconfig = filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// get clientset
	clientset := workflowclientset.NewForConfigOrDie(config)
	api := clientset.WorkflowV1().WorkflowTasks("")

	// create watch channel
	watch, err := api.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	ch := watch.ResultChan()

	for {
		event := <-ch
		workflowtasks, ok := event.Object.(*workflowv1.WorkflowTask)
		if !ok {
			panic("Could not cast to WorkflowTask")
		}
		fmt.Printf("%v\n", workflowtasks.ObjectMeta.Name)
	}
}
