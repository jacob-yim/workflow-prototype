package main

import (
	"context"
	"fmt"
	"path/filepath"

	v1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	cs "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	clientv1 "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned/typed/workflow/v1"
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
	clientset := cs.NewForConfigOrDie(config)
	api := clientset.WorkflowV1().WorkflowTasks("")

	// create watch channel
	watch, err := api.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	events := watch.ResultChan()

	dispatch := make(chan *v1.WorkflowTask)

	// start dummy executor
	go dummyExecutor(dispatch, api)

	// handle watch events
	for {
		event := <-events
		workflowtask, ok := event.Object.(*v1.WorkflowTask)
		if !ok {
			panic("Could not cast to WorkflowTask")
		}

		if event.Type == "ADDED" {
			dispatch <- workflowtask
		}
	}
}

func dummyExecutor(dispatch chan *v1.WorkflowTask, api clientv1.WorkflowTaskInterface) {
	executorType := "testType"

	for {
		task := <-dispatch
		taskType := task.Status.State

		if taskType == executorType {
			task.Status.State = v1.StateExecuting
			api.UpdateStatus(context.TODO(), task, metav1.UpdateOptions{})
			fmt.Println("Performing test task")

			task.Status.State = v1.StateCompleted
			api.UpdateStatus(context.TODO(), task, metav1.UpdateOptions{})

		} else {
			dispatch <- task
		}
	}
}
