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
	// ... fields for struct
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: src=%s", p.SourceURL)

	imageResizeRepository := repository.NewImageResizeRepository()

	data, err := imageResizeRepository.Resize(kraken.Request{
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
	return nil
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
