package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"saga/model"
	"saga/repository"
	"saga/service"

	"github.com/google/uuid"
)

func (r *RabbitMQ) ConsumeOrderCreated(eventRepository repository.ProcessEventRepository) error {

	messages, err := r.Ch.Consume(
		OrderSagaQueue, // IMPORTANT
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to start order consumer: %w",
			err,
		)
	}

	log.Println("Saga order.created consumer started")

	for message := range messages {

		var event OrderCreatedEvent

		if err := json.Unmarshal(
			message.Body,
			&event,
		); err != nil {

			log.Println(
				"failed to unmarshal order.created:",
				err,
			)

			message.Nack(false, false)
			continue
		}

		log.Printf(
			"OrderCreated received: orderID=%s userID=%s",
			event.OrderID,
			event.UserID,
		)
		isProcessed, err := eventRepository.IsProcessed(context.Background(), event.EventID)
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
				event.EventID,
			)

			message.Ack(false)
			continue
		}

		command := ReserveStockCommand{
			Event:       ReserveStockRoutingKey,
			EventID:     uuid.New(),
			OrderID:     event.OrderID,
			UserID:      event.UserID,
			Items:       event.Items,
			TotalAmount: event.TotalAmount,
		}

		err = r.PublishReserveStock(
			context.Background(),
			command,
		)

		if err != nil {

			log.Printf(
				"failed to publish stock.reserve: %v",
				err,
			)

			message.Nack(false, true)
			continue
		}

		log.Printf(
			"stock.reserve published: orderID=%s",
			event.OrderID,
		)
		// Mark the RECEIVED event as processed
		err = eventRepository.MarkProcessed(context.Background(), event.EventID)
		if err != nil {

			log.Printf(
				"event published but failed to mark processed: eventID=%s error=%v",
				event.EventID,
				err,
			)

			// Important:
			// RabbitMQ may deliver this event again.
			// Therefore downstream stock.reserve must also be idempotent.
			message.Nack(false, true)
			continue
		}

		if err := message.Ack(false); err != nil {
			log.Println("failed to ACK:", err)
		}
	}

	return nil
}
func (r *RabbitMQ) ConsumeStockEvent(
	orderService *service.OrderSagaService,
	eventRepository *repository.ProcessEventRepository,
) error {

	messages, err := r.Ch.Consume(
		StockSagaQueue, // IMPORTANT
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

	log.Println("Saga stock event consumer started")

	for message := range messages {

		switch message.RoutingKey {

		// ========================================
		// STOCK RESERVED
		// ========================================

		case StockReservedRoutingKey:

			var event ReserveStockEvent

			if err := json.Unmarshal(
				message.Body,
				&event,
			); err != nil {

				log.Println(
					"invalid stock.reserved event:",
					err,
				)

				message.Nack(false, false)
				continue
			}

			log.Printf(
				"Stock reserved: orderID=%s",
				event.OrderID,
			)
			isProcced, err := eventRepository.IsProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf("failed to check whether event is already proccessed or not %v", err)
				message.Nack(false, true)
				continue
			}
			if isProcced {
				log.Printf("event is already proccessed eventID %v", event.EventID)
				message.Ack(false)
				continue
			}
			// NEXT STEP:
			paymentCommand := PaymentProcessCommand{
				EventID:     uuid.New(),
				Event:       PaymentProcessRoutingKey,
				OrderID:     event.OrderID,
				UserID:      event.UserID,
				TotalAmount: event.TotalAmount,
			}
			// publish payment.process
			err = r.PublishPament(context.Background(), paymentCommand)
			if err != nil {

				log.Printf(
					"failed to publish payment.process: %v",
					err,
				)

				message.Nack(false, true)
				continue

			}
			err = eventRepository.MarkProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf(
					"processed business operation but failed to mark event processed: eventID=%s error=%v",
					event.EventID,
					err,
				)

				message.Nack(false, true)
				continue
			}
			if err := message.Ack(false); err != nil {
				log.Println(
					"failed to ACK:",
					err,
				)
			}
		// ========================================
		// STOCK RESERVATION FAILED
		// ========================================

		case ReservationStockFailedRoutingKey:

			var event ReserveStockFailedEvent

			if err := json.Unmarshal(
				message.Body,
				&event,
			); err != nil {

				log.Println(
					"invalid stock.reservation.failed event:",
					err,
				)

				message.Nack(false, false)
				continue
			}

			log.Printf(
				"Stock reservation failed: orderID=%s reason=%s",
				event.OrderID,
				event.Reason,
			)

			isProcced, err := eventRepository.IsProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf("failed to check whether event is already proccessed or not %v", err)
				message.Nack(false, true)
				continue
			}
			if isProcced {
				log.Printf("event is already proccessed eventID %v", event.EventID)
				message.Ack(false)
				continue
			}

			// ====================================
			// COMPENSATING ACTION
			// ====================================

			err = orderService.CancelOrder(
				context.Background(),
				event.OrderID,
			)

			if err != nil {

				log.Printf(
					"failed to cancel order: %v",
					err,
				)

				// Retry
				message.Nack(false, true)
				continue
			}

			log.Printf(
				"Order cancelled successfully: orderID=%s",
				event.OrderID,
			)
			err = eventRepository.MarkProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf(
					"processed business operation but failed to mark event processed: eventID=%s error=%v",
					event.EventID,
					err,
				)

				message.Nack(false, true)
				continue
			}

			message.Ack(false)
		}
	}

	return nil
}
func (r *RabbitMQ) ConsumePaymentEvent(
	orderService *service.OrderSagaService, productService *service.SagaProductService, eventRepository *repository.ProcessEventRepository,
) error {

	messages, err := r.Ch.Consume(
		PaymentSagaQueue, // IMPORTANT
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

	log.Println("Saga payment event consumer started")

	for message := range messages {

		switch message.RoutingKey {

		// ========================================
		// PAYMENT RESERVED
		// ========================================

		case PaymentSuccessRoutingKey:

			var event PaymentSuccessEvent

			if err := json.Unmarshal(
				message.Body,
				&event,
			); err != nil {

				log.Println(
					"invalid payment.process event:",
					err,
				)

				message.Nack(false, false)
				continue
			}

			log.Printf(
				"Payment Processed: orderID=%s",
				event.OrderID,
			)
			isProcced, err := eventRepository.IsProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf("failed to check whether event is already proccessed or not %v", err)
				message.Nack(false, true)
				continue
			}
			if isProcced {
				log.Printf("event is already proccessed eventID %v", event.EventID)
				message.Ack(false)
				continue
			}
			// NEXT STEP:

			err = orderService.UpdateOrderStatus(
				context.Background(),
				event.OrderID,
				"CONFIRMED",
			)

			if err != nil {

				log.Printf(
					"failed to confirm order: %v",
					err,
				)

				// Retry
				message.Nack(false, true)
				continue
			}

			log.Printf(
				"Order successfully Confirmed: orderID=%s",
				event.OrderID,
			)

			err = eventRepository.MarkProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf(
					"processed business operation but failed to mark event processed: eventID=%s error=%v",
					event.EventID,
					err,
				)

				message.Nack(false, true)
				continue
			}

			if err := message.Ack(false); err != nil {
				log.Println(
					"failed to ACK:",
					err,
				)
			}

		// ========================================
		// PAYMENT FAILED
		// ========================================

		case PaymentFailedRoutingKey:

			var event PaymentFailedEvent

			if err := json.Unmarshal(
				message.Body,
				&event,
			); err != nil {

				log.Println(
					"invalid payment.failed event:",
					err,
				)

				message.Nack(false, false)
				continue
			}

			log.Printf(
				"Payment failed: orderID=%s reason=%s",
				event.OrderID,
				event.Reason,
			)
			isProcced, err := eventRepository.IsProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf("failed to check whether event is already proccessed or not %v", err)
				message.Nack(false, true)
				continue
			}
			if isProcced {
				log.Printf("event is already proccessed eventID %v", event.EventID)
				message.Ack(false)
				continue
			}

			// ========================================
			// 1. Get Order
			// ========================================

			order, err := orderService.GetOrderByID(
				context.Background(),
				event.OrderID,
			)

			if err != nil {

				log.Printf(
					"failed to get order: %v",
					err,
				)

				message.Nack(false, true)
				continue
			}

			// ========================================
			// 2. Build items
			// ========================================

			orderItems := make(
				[]model.OrderEventItem,
				0,
				len(order.Items),
			)

			for _, item := range order.Items {

				orderItems = append(
					orderItems,
					model.OrderEventItem{
						ProductID: item.ProductID,
						Quantity:  item.Quantity,
					},
				)
			}

			// ========================================
			// 3. Release Stock
			// ========================================

			err = productService.ReleaseStock(
				context.Background(),
				orderItems,
			)

			if err != nil {

				log.Printf(
					"failed to release stock: %v",
					err,
				)

				message.Nack(false, true)
				continue
			}

			log.Printf(
				"Stock released: orderID=%s",
				event.OrderID,
			)

			// ========================================
			// 4. Cancel Order
			// ========================================

			err = orderService.CancelOrder(
				context.Background(),
				event.OrderID,
			)

			if err != nil {

				log.Printf(
					"failed to cancel order: %v",
					err,
				)

				message.Nack(false, true)
				continue
			}
			err = eventRepository.MarkProcessed(context.Background(), event.EventID)
			if err != nil {
				log.Printf(
					"processed business operation but failed to mark event processed: eventID=%s error=%v",
					event.EventID,
					err,
				)

				message.Nack(false, true)
				continue
			}

			log.Printf(
				"Order cancelled successfully: orderID=%s",
				event.OrderID,
			)

			// ========================================
			// 5. ACK
			// ========================================

			if err := message.Ack(false); err != nil {
				log.Println(
					"failed to ACK:",
					err,
				)
			}
		}
	}

	return nil
}
