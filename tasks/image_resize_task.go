package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os"
	"task-queue-asynq/model/kraken"
	"task-queue-asynq/repository"
	"task-queue-asynq/tasks/prority"
	taskType "task-queue-asynq/tasks/type"
	"time"
)

type ImageResizePayload struct {
	SourceURL string
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, err
	}
	// tasks options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(taskType.TypeImageResize, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	AsynqClient           *asynq.Client
	ImageResizeRepository repository.ImageResizeRepository
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: src=%s", p.SourceURL)

	data, err := processor.ImageResizeRepository.Resize(kraken.Request{
		Auth: kraken.Auth{
			ApiKey:    os.Getenv("KRAKEN_API_KEY"),
			ApiSecret: os.Getenv("KRAKEN_API_SECRET"),
		},
		Url:  p.SourceURL,
		Wait: true,
		Resize: kraken.Resize{
			Width:    1080,
			Height:   608,
			Strategy: kraken.Auto,
		},
	})

	if err != nil {
		return fmt.Errorf("resize image failed: %v: %w", err, asynq.SkipRetry)
	}

	if !data.Success {
		return fmt.Errorf("Failed, error message: %v, %v ", err, data)
	} else {
		log.Println("Success, Optimized image URL: ", data.KrakedUrl)
	}

	go func() {
		task, errTask := NewFirebaseMessageTask("Success, Optimized image URL", data.KrakedUrl)
		if errTask != nil {
			return
		}

		info, errTask := processor.AsynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(3*time.Minute), asynq.Queue(prority.PriorityDefault))
		if errTask != nil {
			return
		}
		log.Printf("enqueued tasks: id=%s queue=%s", info.ID, info.Queue)
	}()

	go func() {
		task, errTask := NewFirebaseDatabaseTask(p.SourceURL, data.KrakedUrl)
		if errTask != nil {
			return
		}

		info, errTask := processor.AsynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(3*time.Minute), asynq.Queue(prority.PriorityCritical))
		if errTask != nil {
			return
		}

		log.Printf("enqueued tasks: id=%s queue=%s", info.ID, info.Queue)
	}()

	return nil
}

func NewImageProcessor(asynqClient *asynq.Client, resizeRepository repository.ImageResizeRepository) *ImageProcessor {
	return &ImageProcessor{
		AsynqClient:           asynqClient,
		ImageResizeRepository: resizeRepository,
	}
}
