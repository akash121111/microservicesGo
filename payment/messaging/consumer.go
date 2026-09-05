package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"payment/repository"
	"payment/service"

	"github.com/google/uuid"
)

func (r *RabbitMQ) ConsumePayment(
	paymentService *service.PaymentService, eventRepository *repository.ProcessEventRepository,
) error {

	messages, err := r.Ch.Consume(
		PaymentQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to start payment consumer: %w",
			err,
		)
	}

	log.Println("Payment consumer started")

	for message := range messages {

		var command PaymentProcessCommand

		// ----------------------------------------
		// Deserialize command
		// ----------------------------------------

		if err := json.Unmarshal(
			message.Body,
			&command,
		); err != nil {

			log.Printf(
				"invalid payment.process message: %v",
				err,
			)

			// Bad message -> DLQ
			message.Nack(false, false)
			continue
		}

		log.Printf(
			"Payment process received: orderID=%s userID=%s amount=%.2f",
			command.OrderID,
			command.UserID,
			command.TotalAmount,
		)

		// ----------------------------------------
		// Process payment
		// ----------------------------------------

		err := paymentService.ProcessPayment(
			command.UserID,
			command.TotalAmount,
		)

		// ========================================
		// PAYMENT FAILED
		// ========================================

		if err != nil {

			log.Printf(
				"Payment failed: orderID=%s reason=%v",
				command.OrderID,
				err,
			)
			isProcced, isProccedErr := eventRepository.IsProcessed(context.Background(), command.EventID)
			if isProccedErr != nil {
				log.Printf("failed to check whether event is already proccessed or not %v", isProccedErr)
				message.Nack(false, true)
				continue
			}
			if isProcced {
				log.Printf("event is already proccessed eventID %v", command.EventID)
				message.Ack(false)
				continue
			}
			event := PaymentFailedEvent{
				Event:   PaymentFailedRoutingKey,
				EventID: uuid.New(),
				OrderID: command.OrderID,
				UserID:  command.UserID,
				Reason:  err.Error(),
			}

			err = r.PublishPaymentFailed(
				context.Background(),
				event,
			)

			if err != nil {

				log.Printf(
					"failed to publish payment.failed: %v",
					err,
				)

				// Don't ACK.
				// RabbitMQ will retry.
				message.Nack(false, true)
				continue
			}

			err = eventRepository.MarkProcessed(context.Background(), command.EventID)
			if err != nil {
				log.Printf(
					"processed business operation but failed to mark event processed: eventID=%s error=%v",
					command.EventID,
					err,
				)

				message.Nack(false, true)
				continue
			}
			// Event successfully published.
			message.Ack(false)
			continue
		}

		// ========================================
		// PAYMENT SUCCESS
		// ========================================

		log.Printf(
			"Payment successful: orderID=%s amount=%.2f",
			command.OrderID,
			command.TotalAmount,
		)
		isProcced, err := eventRepository.IsProcessed(context.Background(), command.EventID)
		if err != nil {
			log.Printf("failed to check whether event is already proccessed or not %v", err)
			message.Nack(false, true)
			continue
		}
		if isProcced {
			log.Printf("event is already proccessed eventID %v", command.EventID)
			message.Ack(false)
			continue
		}
		event := PaymentSuccessEvent{
			Event:   PaymentSuccessRoutingKey,
			EventID: uuid.New(),
			OrderID: command.OrderID,
			UserID:  command.UserID,
			Amount:  command.TotalAmount,
		}

		err = r.PublishPaymentSuccess(
			context.Background(),
			event,
		)

		if err != nil {

			log.Printf(
				"failed to publish payment.success: %v",
				err,
			)

			message.Nack(false, true)
			continue
		}

		err = eventRepository.MarkProcessed(context.Background(), command.EventID)
		if err != nil {
			log.Printf(
				"processed business operation but failed to mark event processed: eventID=%s error=%v",
				command.EventID,
				err,
			)

			message.Nack(false, true)
			continue
		}
		// Event successfully published.
		message.Ack(false)
	}

	return nil
}
