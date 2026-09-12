package service_test

import (
	"testing"
	"time"

	"eda-demo/api-go/internal/models"
	"eda-demo/api-go/internal/service"
)

func TestProcessOrder_Success(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()

	inEvent := models.OrderCreatedEvent{
		EventID:       "evt-test-001",
		CorrelationID: "corr-test-12345",
		Timestamp:     time.Now().UTC(),
		Type:          "legacy.pedidos.creado",
		Data: models.OrderData{
			OrderID:    "ORD-TEST-01",
			CustomerID: "cli-999",
			Total:      200.00,
			Items: []models.OrderItem{
				{SKU: "APPLE-PROD", Quantity: 2, Price: 50.00},
				{SKU: "BANANA-PROD", Quantity: 1, Price: 100.00},
			},
		},
	}

	// ==========================================
	// Act
	// ==========================================
	outEvent, err := svc.ProcessOrder(inEvent)

	// ==========================================
	// Assert
	// ==========================================
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if outEvent == nil {
		t.Fatal("expected outEvent to be non-nil")
	}
	if outEvent.CorrelationID != inEvent.CorrelationID {
		t.Errorf("expected CorrelationID %s, got %s", inEvent.CorrelationID, outEvent.CorrelationID)
	}
	if outEvent.Type != "facturacion.facturas.generada" {
		t.Errorf("expected Type facturacion.facturas.generada, got %s", outEvent.Type)
	}
	if outEvent.Data.OrderID != inEvent.Data.OrderID {
		t.Errorf("expected OrderID %s, got %s", inEvent.Data.OrderID, outEvent.Data.OrderID)
	}
	if outEvent.Data.Subtotal != 200.00 {
		t.Errorf("expected Subtotal 200.00, got %.2f", outEvent.Data.Subtotal)
	}
	// 19% of 200 is 38.00
	if outEvent.Data.Tax != 38.00 {
		t.Errorf("expected Tax 38.00, got %.2f", outEvent.Data.Tax)
	}
	// Total 238.00
	if outEvent.Data.Total != 238.00 {
		t.Errorf("expected Total 238.00, got %.2f", outEvent.Data.Total)
	}
	if outEvent.Data.Status != "EMITIDA" {
		t.Errorf("expected Status EMITIDA, got %s", outEvent.Data.Status)
	}
	if outEvent.Data.WarehouseAssigned != "WH-NORTH-01" {
		t.Errorf("expected Warehouse WH-NORTH-01 for SKU starting with 'A', got %s", outEvent.Data.WarehouseAssigned)
	}

	// Check audit log
	invoices := svc.GetRecentInvoices()
	if len(invoices) != 1 {
		t.Errorf("expected 1 stored invoice, got %d", len(invoices))
	}
}

func TestProcessOrder_MissingOrderID(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()
	inEvent := models.OrderCreatedEvent{
		Data: models.OrderData{
			OrderID: "",
			Items:   []models.OrderItem{{SKU: "TEST", Quantity: 1, Price: 10}},
		},
	}

	// ==========================================
	// Act
	// ==========================================
	outEvent, err := svc.ProcessOrder(inEvent)

	// ==========================================
	// Assert
	// ==========================================
	if err == nil {
		t.Fatal("expected error for missing order_id, got nil")
	}
	if outEvent != nil {
		t.Errorf("expected nil outEvent on error, got %+v", outEvent)
	}
}

func TestProcessOrder_EmptyItems(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()
	inEvent := models.OrderCreatedEvent{
		Data: models.OrderData{
			OrderID: "ORD-EMPTY",
			Items:   []models.OrderItem{},
		},
	}

	// ==========================================
	// Act
	// ==========================================
	_, err := svc.ProcessOrder(inEvent)

	// ==========================================
	// Assert
	// ==========================================
	if err == nil {
		t.Fatal("expected error for empty items, got nil")
	}
}

func TestProcessOrder_NegativePrice(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()
	inEvent := models.OrderCreatedEvent{
		Data: models.OrderData{
			OrderID: "ORD-NEG",
			Items:   []models.OrderItem{{SKU: "TEST", Quantity: 1, Price: -5.0}},
		},
	}

	// ==========================================
	// Act
	// ==========================================
	_, err := svc.ProcessOrder(inEvent)

	// ==========================================
	// Assert
	// ==========================================
	if err == nil {
		t.Fatal("expected error for negative price, got nil")
	}
}
