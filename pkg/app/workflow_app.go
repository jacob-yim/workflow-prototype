package app

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	v1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	cs "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	clientv1 "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned/typed/workflow/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Task struct {
	Execute  Execution
	TaskType string
}

type Executor struct {
	TaskToExecute  Task
	ThreadPoolSize int
}

type Execution func(*v1.WorkflowTask)

func Start(executors []Executor) {
	// get config
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// get WorkflowTasks clientset
	clientset := cs.NewForConfigOrDie(config)
	api := clientset.WorkflowV1().WorkflowTasks("default")

	dispatchMap := make(map[string]chan *v1.WorkflowTask)

	// start executors
	for _, exec := range executors {
		taskToExecute := exec.TaskToExecute
		taskType := taskToExecute.TaskType

		dispatch := make(chan *v1.WorkflowTask)
		dispatchMap[taskType] = dispatch

		for i := 0; i < exec.ThreadPoolSize; i++ {
			go taskExecutor(api, dispatch, taskToExecute)
		}
	}

	// start watcher
	go taskWatcher(api, dispatchMap)
}

func taskWatcher(api clientv1.WorkflowTaskInterface, dispatchMap map[string]chan *v1.WorkflowTask) {
	// create watch channel
	watch, err := api.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	events := watch.ResultChan()

	for event := range events {
		taskResource, ok := event.Object.(*v1.WorkflowTask)
		if !ok {
			panic("Could not cast to WorkflowTask")
		}

		// dispatch task
		if event.Type == "ADDED" {
			dispatch := dispatchMap[taskResource.Spec.Type]
			dispatch <- taskResource
		}
	}
}

func taskExecutor(api clientv1.WorkflowTaskInterface, dispatch chan *v1.WorkflowTask, taskToExecute Task) {
	executorHostname, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}

	executorID := executorHostname + uuid.New().String()
	taskCount := 0

	for taskResource := range dispatch {
		if taskResource.Status.Executor == "" {
			taskName := taskResource.Name

			taskResource.Status.Executor = executorID

			taskResource.Status.State = v1.StateExecuting
			taskResource, err := api.UpdateStatus(context.TODO(), taskResource, metav1.UpdateOptions{})
			if errors.IsConflict(err) {
				continue
			} else if err != nil {
				panic(err.Error())
			}

			log.Printf("Task %v executing...\n", taskName)

			taskToExecute.Execute(taskResource)

			taskCount += 1
			log.Printf("Task %v completed. Executor total: %v\n", taskName, taskCount)

			taskResource.Status.State = v1.StateCompleted
			_, err = api.UpdateStatus(context.TODO(), taskResource, metav1.UpdateOptions{})
			if errors.IsConflict(err) {
				continue
			} else if err != nil {
				panic(err.Error())
			}
		}
	}
}
