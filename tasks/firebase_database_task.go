package tasks

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	taskType "task-queue-asynq/tasks/type"
	"time"
)

type FirebaseDatabasePayload struct {
	SourceURL string
	KrakedURL string
}

func NewFirebaseDatabaseTask(sourceUrl string, krakedUrl string) (*asynq.Task, error) {
	payload, err := json.Marshal(FirebaseDatabasePayload{
		SourceURL: sourceUrl,
		KrakedURL: krakedUrl,
	})
	if err != nil {
		return nil, err
	}
	// tasks options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(taskType.TypeFirebaseDatabase, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

// FirebaseDatabaseProcessor implements asynq.Handler interface.
type FirebaseDatabaseProcessor struct {
	FirebaseApp *firebase.App
}

func (processor *FirebaseDatabaseProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p FirebaseDatabasePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Push Firebase Database: src=%s", p.SourceURL)

	firebaseDatabase, err := processor.FirebaseApp.Database(ctx)
	if err != nil {
		return fmt.Errorf("init firebase firebaseDatabase failed: %v: %w", err, asynq.SkipRetry)
	}

	imageRef := firebaseDatabase.NewRef("images")
	imageRef, err = imageRef.Push(ctx, map[string]string{
		"source_url": p.SourceURL,
		"kraked_url": p.KrakedURL,
	})
	if err != nil {
		return fmt.Errorf("imageRef message failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Println("Success, imageRef message: ", imageRef.Path)
	return nil
}

func NewFirebaseDatabaseProcessor(firebaseApp *firebase.App) *FirebaseDatabaseProcessor {
	return &FirebaseDatabaseProcessor{
		FirebaseApp: firebaseApp,
	}
}
