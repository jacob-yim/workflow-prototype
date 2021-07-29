package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	v1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	cs "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	clientv1 "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned/typed/workflow/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const TASKS = 10
const EXECUTORS = 3
const EXEC_TIME = 5 // execution time in seconds

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

	// start executors
	for i := 0; i < EXECUTORS; i++ {
		go taskExecutor(api, dispatch, "testType")
	}

	// start scheduler
	go taskScheduler(api, "testType", TASKS)

	// start watcher
	taskWatcher(api, dispatch)
}

func taskScheduler(api clientv1.WorkflowTaskInterface, taskType string, replicas int) {
	for i := 0; i < replicas; i++ {
		id := uuid.New()
		taskName := "test-" + strconv.Itoa(i) + "-" + id.String()

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

func taskWatcher(api clientv1.WorkflowTaskInterface, dispatch chan<- *v1.WorkflowTask) {
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

func taskExecutor(api clientv1.WorkflowTaskInterface, dispatch chan *v1.WorkflowTask, executorType string) {
	var taskCount = 0
	for task := range dispatch {

		if task.Status.Executor == "" {
			taskName := task.Name
			taskType := task.Spec.Type

			if taskType == executorType {
				var err error
				executorID := uuid.New()
				var executorHostname = ""
				executorHostname, err = os.Hostname()
				if err != nil {
					panic(err.Error())
				}
				task.Status.Executor = executorHostname + executorID.String()

				task.Status.State = v1.StateExecuting
				task, err := api.UpdateStatus(context.TODO(), task, metav1.UpdateOptions{})
				if errors.IsConflict(err) {
					continue
				} else if err != nil {
					panic(err.Error())
				}

				log.Printf("Task %v executing...\n", taskName)

				//simulate execution
				time.Sleep(EXEC_TIME * time.Second)

				taskCount += 1
				log.Printf("Task %v completed. Executor total: %v\n", taskName, taskCount)

				task.Status.State = v1.StateCompleted
				_, err = api.UpdateStatus(context.TODO(), task, metav1.UpdateOptions{})
				if errors.IsConflict(err) {
					continue
				} else if err != nil {
					panic(err.Error())
				}

			} else {
				// send task to a different executor
				dispatch <- task
			}
		}
	}
}
