package tasks

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	taskType "task-queue-asynq/tasks/type"
	"time"
)

type FirebaseMessagePayload struct {
	MessageBody string
	SourceURL   string
}

func NewFirebaseMessageTask(messageBody string, imageUrl string) (*asynq.Task, error) {
	payload, err := json.Marshal(FirebaseMessagePayload{
		MessageBody: messageBody,
		SourceURL:   imageUrl,
	})
	if err != nil {
		return nil, err
	}
	// tasks options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(taskType.TypeFirebaseMessage, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

// FirebaseMessageProcessor implements asynq.Handler interface.
type FirebaseMessageProcessor struct {
	FirebaseApp *firebase.App
}

func (processor *FirebaseMessageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p FirebaseMessagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Message: src=%s", p.SourceURL)

	firebaseMessage, err := processor.FirebaseApp.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("init firebase firebaseMessage failed: %v: %w", err, asynq.SkipRetry)
	}

	send, err := firebaseMessage.Send(ctx, &messaging.Message{
		Data: nil,
		Notification: &messaging.Notification{
			Title:    "",
			Body:     p.MessageBody,
			ImageURL: p.SourceURL,
		},
		Topic: "image-resize",
	})
	if err != nil {
		return fmt.Errorf("send message failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Println("Success, send message: ", send)
	return nil
}

func NewFirebaseMessageProcessor(firebaseApp *firebase.App) *FirebaseMessageProcessor {
	return &FirebaseMessageProcessor{
		FirebaseApp: firebaseApp,
	}
}
