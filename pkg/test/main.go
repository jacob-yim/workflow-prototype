package main

import (
	"context"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	v1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	app "github.com/jacob-yim/workflow-prototype/pkg/app"
	cs "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned"
	clientv1 "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned/typed/workflow/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const TASKS = 20
const THREAD_POOL_SIZE = 3
const EXEC_TIME_SECONDS = 5 // execution time in seconds
const TIMEOUT_MINUTES = 10  // test timeout in minutes

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

	executors := make([]app.Executor, 1)
	executors[0] = app.Executor{Task: executeTestTask, TaskType: "testType", ThreadPoolSize: THREAD_POOL_SIZE}

	app.Start(config, executors)

	// start scheduler
	go taskScheduler(api, "testType", TASKS)

	time.Sleep(TIMEOUT_MINUTES * time.Minute)
}

func executeTestTask(taskResource *v1.WorkflowTask) {
	time.Sleep(EXEC_TIME_SECONDS * time.Second)
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

		log.Printf("Scheduled task %v\n", taskName)
	}
}
