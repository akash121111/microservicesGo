package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"product/repository"
	"product/service"

	"github.com/google/uuid"
)

func (r *RabbitMQ) ConsumeReserveStock(
	productService *service.ProductService,
	eventRepository *repository.ProcessEventRepository,
) error {

	messages, err := r.Ch.Consume(
		StockQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to start stock consumer: %w",
			err,
		)
	}

	log.Println("Product stock consumer started")

	for message := range messages {

		var command ReserveStockCommand

		err := json.Unmarshal(
			message.Body,
			&command,
		)

		if err != nil {
			log.Printf(
				"invalid stock.reserve message: %v",
				err,
			)

			if err := r.PublishToDLQ(message); err != nil {
				log.Printf(
					"failed to publish invalid message to DLQ: %v",
					err,
				)

				message.Nack(false, true)
				continue
			}
			message.Ack(false)
			continue
		}

		log.Printf(
			"ReserveStock received: orderID=%s",
			command.OrderID,
		)
		//idempotency
		isProcessed, err := eventRepository.IsProcessed(context.Background(), command.EventID)
		if err != nil {
			log.Println(
				"error checking processed event:",
				err,
			)

			message.Nack(false, true)
			continue
		}
		if isProcessed {
			log.Printf(
				"event already processed: eventID=%s",
				command.EventID,
			)

			message.Ack(false)
			continue
		}

		// --------------------------------
		// Reserve ALL stock
		// --------------------------------

		err = productService.ReserveProductStock(
			context.Background(),
			command.Items,
		)

		if err != nil {

			log.Printf(
				"failed to reserve stock: orderID=%s error=%v",
				command.OrderID,
				err,
			)

			// Nothing was reserved because
			// ProductService uses a transaction.
			// Now tell Saga that reservation failed.

			if err.Error() == "Insufficient stock" {
				event := ReserveStockFailedEvent{
					Event:   ReservationStockFailedRoutingKey,
					EventID: uuid.New(),
					OrderID: command.OrderID,
					Items:   command.Items,
					Reason:  err.Error(),
				}

				if publishErr := r.PublishStockReservationFailed(
					context.Background(),
					event,
				); publishErr != nil {

					log.Printf(
						"failed to publish stock reservation failed event: %v",
						publishErr,
					)

					r.RetryMessage(message)
					continue
				}

				err = eventRepository.MarkProcessed(context.Background(), command.EventID)
				if err != nil {

					log.Printf(
						"event published but failed to mark processed: eventID=%s error=%v",
						event.EventID,
						err,
					)

					r.RetryMessage(message)
					continue
				}

				// Event successfully published.
				message.Ack(false)
				continue
			}
			// ------------------------------------
			// Temporary/system failure
			// ------------------------------------

			r.RetryMessage(message)
			continue

		}

		// --------------------------------
		// ALL stock reserved successfully
		// --------------------------------

		log.Printf(
			"stock reserved successfully: orderID=%s",
			command.OrderID,
		)

		event := ReserveStockCommand{
			Event:       StockReservedRoutingKey,
			EventID:     uuid.New(),
			OrderID:     command.OrderID,
			UserID:      command.UserID,
			Items:       command.Items,
			TotalAmount: command.TotalAmount,
		}

		isProcced, err := eventRepository.IsProcessed(context.Background(), command.EventID)
		if err != nil {
			log.Println(
				"error checking processed event:",
				err,
			)

			r.RetryMessage(message)
			continue
		}
		if isProcced {
			log.Printf(
				"event already processed: eventID=%s",
				command.EventID,
			)

			message.Ack(false)
			continue
		}
		if err := r.PublishStockReservation(
			context.Background(),
			event,
		); err != nil {

			log.Printf(
				"failed to publish stock.reserved: %v",
				err,
			)

			r.RetryMessage(message)
			continue
		}
		err = eventRepository.MarkProcessed(context.Background(), command.EventID)
		if err != nil {

			log.Printf(
				"event published but failed to mark processed: eventID=%s error=%v",
				event.EventID,
				err,
			)

			r.RetryMessage(message)
			continue
		}
		message.Ack(false)
	}

	return nil
}
