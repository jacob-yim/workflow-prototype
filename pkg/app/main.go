package main

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

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
	api := clientset.WorkflowV1().WorkflowTasks("default")

	// create watch channel
	watch, err := api.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	events := watch.ResultChan()

	dispatch := make(chan *v1.WorkflowTask)

	// start dummy executor
	go dummyExecutor(dispatch, api, "testType")

	// handle watch events
	for {
		event := <-events
		task, ok := event.Object.(*v1.WorkflowTask)
		if !ok {
			panic("Could not cast to WorkflowTask")
		}

		// dispatch task
		if event.Type == "ADDED" {
			dispatch <- task
		}
	}
}

func dummyExecutor(dispatch chan *v1.WorkflowTask, api clientv1.WorkflowTaskInterface, executorType string) {
	for {
		task := <-dispatch
		taskType := task.Spec.Type

		// check that task type matches executor
		if taskType == executorType && task.Status.Executor == "" {
			task.Status.Executor = "pod name"
			task.Status.State = v1.StateExecuting
			task, err := api.UpdateStatus(context.TODO(), task, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}

			fmt.Println("Task executing...")

			//simulate execution
			time.Sleep(5 * time.Second)

			fmt.Println("Task completed")
			task.Status.State = v1.StateCompleted
			task, err = api.UpdateStatus(context.TODO(), task, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
		} else {
			// send task to a different executor
			dispatch <- task
		}
	}
}
