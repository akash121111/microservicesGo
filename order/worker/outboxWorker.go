package worker

import (
	"context"
	"log"
	"order/messaging"
	"order/repository"
	"time"
)

type OutboxWorker struct {
	repo     *repository.OutboxEventRepository
	rabbitmq *messaging.RabbitMQ
}

func NewOutboxWorker(repo *repository.OutboxEventRepository, rabbitmq *messaging.RabbitMQ) *OutboxWorker {
	return &OutboxWorker{
		repo:     repo,
		rabbitmq: rabbitmq,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("worker stop")
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}
func (w *OutboxWorker) process(ctx context.Context) {

	events, err := w.repo.GetAllNotPublishOutBoxEvent(ctx)

	if err != nil {
		log.Printf(
			"failed to get outbox events: %v",
			err,
		)
		return
	}

	for _, event := range events {

		err := w.rabbitmq.PublishOutboxEvent(
			ctx,
			event.RoutingKey,
			event.Payload,
		)

		if err != nil {
			log.Printf(
				"failed to publish outbox event %s: %v",
				event.ID,
				err,
			)

			// IMPORTANT:
			// Don't mark it as published.
			// It will be retried next cycle.
			continue
		}

		err = w.repo.UpdateBoxEventPublishById(
			ctx,
			event.ID,
		)

		if err != nil {
			log.Printf(
				"failed to mark outbox event %s as published: %v",
				event.ID,
				err,
			)

			continue
		}

		log.Printf(
			"outbox event published: id=%s routingKey=%s",
			event.ID,
			event.RoutingKey,
		)
	}
}
