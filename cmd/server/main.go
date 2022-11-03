package main

import (
	"github.com/hibiken/asynq"
	"log"
	"task-queue-asynq/configs"
	"task-queue-asynq/tasks"
	"task-queue-asynq/tasks/prority"
	taskType "task-queue-asynq/tasks/type"
)

var envConf configs.EnvConf

func init() {
	envConf = configs.NewEnv()
}

func main() {
	redisAddr := envConf.RedisUrl

	asynqServer := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				prority.PriorityCritical: 6,
				prority.PriorityDefault:  3,
				prority.PriorityLow:      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(taskType.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.Handle(taskType.TypeImageResize, tasks.NewImageProcessor())
	// ...register other handlers...

	if err := asynqServer.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

}
