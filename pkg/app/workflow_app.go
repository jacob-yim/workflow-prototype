package app

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	v1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	cs "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	clientv1 "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned/typed/workflow/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Executor struct {
	Task           Execute
	TaskType       string
	ThreadPoolSize int
}

type Execute func(*v1.WorkflowTask) error

func Start(config *rest.Config, executors []Executor) {
	// get WorkflowTasks clientset
	clientset := cs.NewForConfigOrDie(config)
	api := clientset.WorkflowV1().WorkflowTasks("default")

	dispatchMap := make(map[string]chan *v1.WorkflowTask)

	// start executors
	for _, exec := range executors {
		taskToExecute := exec.Task
		taskType := exec.TaskType

		dispatch := make(chan *v1.WorkflowTask)
		dispatchMap[taskType] = dispatch

		for i := 0; i < exec.ThreadPoolSize; i++ {
			go taskExecutor(api, dispatch, taskToExecute, i)
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
			taskType := taskResource.Spec.Type

			if dispatch, ok := dispatchMap[taskType]; ok {
				dispatch <- taskResource
			} else {
				log.Printf("No executor found for task of type %v\n", taskType)
			}
		}
	}
}

func taskExecutor(api clientv1.WorkflowTaskInterface, dispatch chan *v1.WorkflowTask, execute Execute, execNum int) {
	executorHostname, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}

	executorID := executorHostname + "-" + strconv.Itoa(execNum)

	for taskResource := range dispatch {
		if taskResource.Status.Executor == "" {
			taskName := taskResource.Name

			taskResource.Status.Executor = executorID
			taskResource.Status.State = v1.StateExecuting
			taskResource.Status.StartTimeUTC = time.Now().UTC().String()

			taskResource, err := api.UpdateStatus(context.TODO(), taskResource, metav1.UpdateOptions{})
			if errors.IsConflict(err) {
				continue
			} else if err != nil {
				panic(err.Error())
			}

			log.Printf("%v: Task %v executing...\n", executorID, taskName)

			err = execute(taskResource)
			if err != nil {
				log.Printf("%v: Task %v failed with error: %v\n", executorID, taskName, err.Error())

				taskResource.Status.State = v1.StateFailed
				taskResource.Status.Error = err.Error()
			} else {
				log.Printf("%v: Task %v completed.\n", executorID, taskName)

				taskResource.Status.State = v1.StateCompleted
				taskResource.Status.CompletionTimeUTC = time.Now().UTC().String()
			}

			_, err = api.UpdateStatus(context.TODO(), taskResource, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
		}
	}
}
