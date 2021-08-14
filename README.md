# workflow-prototype-go #

## About ##

A Golang prototype for Kubernetes-native Nirmata workflows. Watches [WorkflowTask Kubernetes custom resources](config/crd/bases/nirmata.com_workflowtasks.yaml) and executes user-defined Go tasks upon resource creation. Tasks of different types may be assigned individually to different executors and run in parallel.

A Java implementation can be found [here](https://github.com/aupadhyay3/workflow-prototype-java). 

## Usage ##

First, define a slice of Executors. The Executor struct contains Task, a user-defined function to be executed, TaskType, a string, and an integer ThreadPoolSize. This assigns tasks of type TaskType to an executor with a fixed thread pool of size ThreadPoolSize that runs Task.

    executors := make([]app.Executor, 1)
	executors[0] = app.Executor{Task: executeTestTask, TaskType: "type", ThreadPoolSize: 3}

app.Start() starts the workflow application by taking in a client-go config and the slice of executors.

    app.Start(config, executors)