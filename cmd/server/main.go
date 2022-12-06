package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/hibiken/asynq"
	"google.golang.org/api/option"
	"log"
	"os"
	"path"
	"task-queue-asynq/configs"
	"task-queue-asynq/repository"
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

	getwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get wd")
	}
	opt := option.WithCredentialsFile(path.Join(getwd, "serviceAccountKey.json"))

	firebaseApp, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID:   envConf.FirebaseProject,
		DatabaseURL: envConf.FirebaseDatabaseUrl,
	}, opt)
	if err != nil {
		log.Fatalf("error initializing firebaseApp: %v\n", err)
	}

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

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer asynqClient.Close()

	imageResizeRepository := repository.NewImageResizeRepository()

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(taskType.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.Handle(taskType.TypeImageResize, tasks.NewImageProcessor(asynqClient, imageResizeRepository))
	mux.Handle(taskType.TypeFirebaseMessage, tasks.NewFirebaseMessageProcessor(firebaseApp))
	mux.Handle(taskType.TypeFirebaseDatabase, tasks.NewFirebaseDatabaseProcessor(firebaseApp))
	// ...register other handlers...

	if err := asynqServer.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

}
