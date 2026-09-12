package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"eda-demo/api-go/internal/config"
	"eda-demo/api-go/internal/models"
	"eda-demo/api-go/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

// OrderConsumer coordinates consuming orders, processing them, and publishing invoices.
type OrderConsumer struct {
	cfg     *config.Config
	service *service.InventoryBillingService
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewOrderConsumer constructs a new consumer instance.
func NewOrderConsumer(cfg *config.Config, svc *service.InventoryBillingService) *OrderConsumer {
	return &OrderConsumer{
		cfg:     cfg,
		service: svc,
	}
}

// Start connects to RabbitMQ and begins consuming messages until context is cancelled.
func (c *OrderConsumer) Start(ctx context.Context) error {
	var err error
	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		log.Printf("[OrderConsumer] Connecting to RabbitMQ at %s (attempt %d/%d)...", c.cfg.RabbitMQURL, i, maxRetries)
		c.conn, err = amqp.Dial(c.cfg.RabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("[OrderConsumer] RabbitMQ connection failed: %v. Retrying in 3s...", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("could not connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open AMQP channel: %w", err)
	}

	// Fair dispatch
	if err := c.channel.Qos(10, 0, false); err != nil {
		return fmt.Errorf("failed to set channel QoS: %w", err)
	}

	deliveries, err := c.channel.Consume(
		c.cfg.QueueName, // queue name: api-go.procesamiento_inventario
		"api-go-worker", // consumer tag
		false,           // autoAck = false (manual ACK mandatory)
		false,           // exclusive
		false,           // noLocal
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer on queue %s: %w", c.cfg.QueueName, err)
	}

	log.Printf("[OrderConsumer] Successfully subscribed to queue '%s'. Awaiting orders...", c.cfg.QueueName)

	for {
		select {
		case <-ctx.Done():
			log.Println("[OrderConsumer] Stopping consumer (context cancelled)...")
			return nil
		case d, ok := <-deliveries:
			if !ok {
				log.Println("[OrderConsumer] Deliveries channel closed by broker.")
				return nil
			}
			c.handleDelivery(d)
		}
	}
}

func (c *OrderConsumer) handleDelivery(d amqp.Delivery) {
	var orderEvent models.OrderCreatedEvent
	if err := json.Unmarshal(d.Body, &orderEvent); err != nil {
		log.Printf("[OrderConsumer] ERROR: Malformed message payload: %v. Sending to DLQ.", err)
		// Nack with requeue=false routes message to Dead Letter Exchange
		_ = d.Nack(false, false)
		return
	}

	// Preserve or extract correlation ID
	if orderEvent.CorrelationID == "" && d.CorrelationId != "" {
		orderEvent.CorrelationID = d.CorrelationId
	}

	log.Printf("[OrderConsumer] Processing order %s for customer %s [CorrelationID: %s]",
		orderEvent.Data.OrderID, orderEvent.Data.CustomerID, orderEvent.CorrelationID)

	invoiceEvent, err := c.service.ProcessOrder(orderEvent)
	if err != nil {
		log.Printf("[OrderConsumer] ERROR: Business failure processing order %s: %v. Sending to DLQ.",
			orderEvent.Data.OrderID, err)
		_ = d.Nack(false, false)
		return
	}

	// Publish generated invoice event
	body, err := json.Marshal(invoiceEvent)
	if err != nil {
		log.Printf("[OrderConsumer] ERROR: Failed to marshal invoice event: %v", err)
		_ = d.Nack(false, false)
		return
	}

	routingKey := "facturacion.facturas.generada"
	err = c.channel.Publish(
		c.cfg.ExchangeName, // sistema.eventos.bus
		routingKey,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			DeliveryMode:  amqp.Persistent,
			CorrelationId: invoiceEvent.CorrelationID,
			Timestamp:     invoiceEvent.Timestamp,
			Body:          body,
		},
	)
	if err != nil {
		log.Printf("[OrderConsumer] ERROR: Failed to publish invoice event to exchange: %v. Requeuing order.", err)
		_ = d.Nack(false, true) // Requeue since broker publication failed
		return
	}

	// Acknowledge original message
	if err := d.Ack(false); err != nil {
		log.Printf("[OrderConsumer] Warning: Ack failed: %v", err)
	} else {
		log.Printf("[OrderConsumer] SUCCESS: Invoice %s generated for order %s (Total: $%.2f, Warehouse: %s) [CorrelationID: %s]. Event published -> %s",
			invoiceEvent.Data.InvoiceID,
			orderEvent.Data.OrderID,
			invoiceEvent.Data.Total,
			invoiceEvent.Data.WarehouseAssigned,
			invoiceEvent.CorrelationID,
			routingKey,
		)
	}
}

// Close gracefully closes channel and connection.
func (c *OrderConsumer) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
