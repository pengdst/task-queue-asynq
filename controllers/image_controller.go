package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"log"
	"task-queue-asynq/model/web"
	"task-queue-asynq/tasks"
	"task-queue-asynq/tasks/prority"
	"time"
)

type ImageController interface {
	Resize(ctx *gin.Context)
}

type ImageControllerImpl struct {
	asynqClient *asynq.Client
}

func (i *ImageControllerImpl) Resize(ctx *gin.Context) {
	var payload web.ImageResizePayload
	ctx.Bind(&payload)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune tasks processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	task, err := tasks.NewImageResizeTask(payload.ImageUrl)
	if err != nil {
		log.Fatalf("could not create tasks: %v", err)
	}
	info, err := i.asynqClient.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute), asynq.Queue(prority.PriorityLow))
	if err != nil {
		log.Fatalf("could not enqueue tasks: %v", err)
	}
	log.Printf("enqueued tasks: id=%s queue=%s", info.ID, info.Queue)
}

func NewImageController(asyncClient *asynq.Client) ImageController {
	return &ImageControllerImpl{asynqClient: asyncClient}
}
