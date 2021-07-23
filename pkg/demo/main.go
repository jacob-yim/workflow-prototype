package main

import (
	"context"
	"fmt"
	"path/filepath"

	workflowv1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	workflowcs "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// get config
	var home = homedir.HomeDir()
	var kubeconfig = filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// get WorkflowTasks clientset
	clientset := workflowcs.NewForConfigOrDie(config)
	api := clientset.WorkflowV1().WorkflowTasks("")

	// create watch channel
	watch, err := api.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	events := watch.ResultChan()

	// handle watch events
	for {
		event := <-events
		workflowtasks, ok := event.Object.(*workflowv1.WorkflowTask)
		if !ok {
			panic("Could not cast to WorkflowTask")
		}

		if event.Type == "ADDED" {
			fmt.Printf("%v\n", workflowtasks.ObjectMeta.Name)
		}
	}
}
