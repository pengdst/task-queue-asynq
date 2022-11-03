package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"log"
	"net/http"
	"task-queue-asynq/exceptions"
	"task-queue-asynq/model/web"
	"task-queue-asynq/tasks"
	"task-queue-asynq/tasks/prority"
	"time"
)

type EmailController interface {
	Delivery(ctx *gin.Context)
}

type EmailControllerImpl struct {
	validate    *validator.Validate
	asynqClient *asynq.Client
}

func (e *EmailControllerImpl) Delivery(ctx *gin.Context) {
	var payload web.EmailDeliveryPayload
	ctx.Bind(&payload)

	err := e.validate.Struct(payload)
	if err != nil {
		panic(exceptions.NewErrorBadRequest(err.Error()))
	}

	// ------------------------------------------------------
	// Example 1: Enqueue tasks to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := tasks.NewEmailDeliveryTask(payload.UserId, payload.TemplateId)
	if err != nil {
		log.Fatalf("could not create tasks: %v", err)
	}
	info, err := e.asynqClient.Enqueue(task, asynq.Queue(prority.PriorityCritical))
	if err != nil {
		log.Fatalf("could not enqueue tasks: %v", err)
	}
	log.Printf("enqueued tasks: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule tasks to be processed in the future.
	//            Use ProcessIn option.
	// ------------------------------------------------------------

	info, err = e.asynqClient.Enqueue(task, asynq.ProcessIn(10*time.Minute))
	if err != nil {
		log.Fatalf("could not schedule tasks: %v", err)
	}
	log.Printf("enqueued tasks: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 3: Schedule tasks to be processed in the future.
	//            Use ProcessAt option.
	// ------------------------------------------------------------

	info, err = e.asynqClient.Enqueue(task, asynq.ProcessAt(time.Now().Add(5*time.Minute)))
	if err != nil {
		log.Fatalf("could not schedule tasks: %v", err)
	}
	log.Printf("enqueued tasks: id=%s queue=%s", info.ID, info.Queue)

	ctx.JSON(http.StatusOK, map[string]any{
		"message": fmt.Sprintf("enqueued tasks: id=%s queue=%s", info.ID, info.Queue),
	})
}

func NewEmailController(validate *validator.Validate, asyncClient *asynq.Client) EmailController {
	return &EmailControllerImpl{validate: validate, asynqClient: asyncClient}
}
