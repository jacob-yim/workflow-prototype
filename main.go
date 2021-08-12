package main

import (
	"time"

	v1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	app "github.com/jacob-yim/workflow-prototype/pkg/app"
	"k8s.io/client-go/rest"
)

const THREAD_POOL_SIZE = 3
const EXEC_TIME_SECONDS = 5 // execution time in seconds
const TIMEOUT_MINUTES = 10  // test timeout in minutes

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	executors := make([]app.Executor, 1)
	executors[0] = app.Executor{Task: executeSampleTask, TaskType: "sample", ThreadPoolSize: THREAD_POOL_SIZE}

	app.Start(config, executors)

	time.Sleep(TIMEOUT_MINUTES * time.Minute)
}

func executeSampleTask(taskResource *v1.WorkflowTask) {
	time.Sleep(EXEC_TIME_SECONDS * time.Second)
}
