package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	v1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	cs "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	clientv1 "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned/typed/workflow/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const TASKS = 20
const EXECUTORS = 3

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

	dispatch := make(chan *v1.WorkflowTask)
	var wg sync.WaitGroup

	// start executors
	for i := 0; i < EXECUTORS; i++ {
		wg.Add(1)
		go taskExecutor(api, dispatch, "testType", &wg)
	}

	wg.Add(2)

	// start watcher
	go taskWatcher(api, dispatch, &wg)

	// start scheduler
	go taskScheduler(api, "testType", TASKS, &wg)

	wg.Wait()
}

func taskScheduler(api clientv1.WorkflowTaskInterface, taskType string, replicas int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < replicas; i++ {
		taskName := "test-task-" + strconv.Itoa(i)

		task := &v1.WorkflowTask{
			ObjectMeta: metav1.ObjectMeta{Name: taskName},
			Spec: v1.WorkflowTaskSpec{
				Type: taskType,
			},
		}

		_, err := api.Create(context.TODO(), task, metav1.CreateOptions{})
		if err != nil {
			panic(err.Error())
		}
	}
}

func taskWatcher(api clientv1.WorkflowTaskInterface, dispatch chan<- *v1.WorkflowTask, wg *sync.WaitGroup) {
	defer wg.Done()

	// create watch channel
	watch, err := api.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	events := watch.ResultChan()
	for event := range events {
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

func taskExecutor(api clientv1.WorkflowTaskInterface, dispatch chan *v1.WorkflowTask, executorType string, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range dispatch {
		if task.Status.Executor == "" {
			taskType := task.Spec.Type

			if taskType == executorType {
				task.Status.Executor = "pod name"
				task.Status.State = v1.StateExecuting
				task, err := api.UpdateStatus(context.TODO(), task, metav1.UpdateOptions{})
				if err != nil {
					panic(err.Error())
				}

				taskName := task.ObjectMeta.Name

				fmt.Printf("Task %v executing...\n", taskName)

				//simulate execution
				time.Sleep(1 * time.Second)

				fmt.Printf("Task %v completed\n", taskName)

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
}
