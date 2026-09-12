package service

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"eda-demo/api-go/internal/models"

	"github.com/google/uuid"
)

// InventoryBillingService handles inventory reservations and fiscal invoice generation.
type InventoryBillingService struct {
	mu       sync.RWMutex
	invoices []models.InvoiceData
}

// NewInventoryBillingService initializes a new service instance.
func NewInventoryBillingService() *InventoryBillingService {
	return &InventoryBillingService{
		invoices: make([]models.InvoiceData, 0),
	}
}

// ProcessOrder reserves warehouse stock and produces a fiscal invoice.
func (s *InventoryBillingService) ProcessOrder(event models.OrderCreatedEvent) (*models.InvoiceGeneratedEvent, error) {
	if event.Data.OrderID == "" {
		return nil, errors.New("missing order_id in event payload")
	}
	if len(event.Data.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// Calculate subtotal from items to ensure financial integrity
	var subtotal float64
	for _, item := range event.Data.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity %d for SKU %s", item.Quantity, item.SKU)
		}
		if item.Price <= 0 {
			return nil, fmt.Errorf("invalid price %.2f for SKU %s", item.Price, item.SKU)
		}
		subtotal += float64(item.Quantity) * item.Price
	}

	subtotal = math.Round(subtotal*100) / 100
	taxRate := 0.19 // 19% standard VAT
	tax := math.Round(subtotal*taxRate*100) / 100
	total := math.Round((subtotal+tax)*100) / 100

	// Select warehouse deterministically based on first SKU
	warehouse := "WH-CENTRAL-01"
	if len(event.Data.Items) > 0 {
		switch {
		case len(event.Data.Items[0].SKU) > 0 && event.Data.Items[0].SKU[0] < 'M':
			warehouse = "WH-NORTH-01"
		case len(event.Data.Items[0].SKU) > 0 && event.Data.Items[0].SKU[0] >= 'M':
			warehouse = "WH-SOUTH-02"
		}
	}

	invoiceID := fmt.Sprintf("INV-%s", uuid.New().String()[:8])

	invoiceData := models.InvoiceData{
		InvoiceID:         invoiceID,
		OrderID:           event.Data.OrderID,
		CustomerID:        event.Data.CustomerID,
		Subtotal:          subtotal,
		Tax:               tax,
		Total:             total,
		WarehouseAssigned: warehouse,
		Status:            "EMITIDA",
		Items:             event.Data.Items,
		Links: map[string]models.Link{
			"self":       {Href: fmt.Sprintf("/invoices/%s", invoiceID), Method: "GET"},
			"docs":       {Href: "/docs", Method: "GET"},
			"collection": {Href: "/invoices", Method: "GET"},
		},
	}

	s.mu.Lock()
	s.invoices = append(s.invoices, invoiceData)
	s.mu.Unlock()

	correlationID := event.CorrelationID
	if correlationID == "" && event.ID != "" {
		correlationID = event.ID
	}
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	eventID := uuid.New().String()
	now := time.Now().UTC()

	outEvent := &models.InvoiceGeneratedEvent{
		SpecVersion:     "1.0",
		ID:              eventID,
		Source:          "/sistema-facturacion/facturas",
		Type:            "facturacion.facturas.generada",
		DataContentType: "application/json",
		Time:            now,
		CorrelationID:   correlationID,
		EventID:         eventID,
		Timestamp:       now,
		Data:            invoiceData,
	}

	return outEvent, nil
}

// GetRecentInvoices returns a copy of all stored invoices for audit and API inspection.
func (s *InventoryBillingService) GetRecentInvoices() []models.InvoiceData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.InvoiceData, len(s.invoices))
	copy(result, s.invoices)
	return result
}

// GetInvoiceByID returns a single invoice by its ID if found.
func (s *InventoryBillingService) GetInvoiceByID(id string) (*models.InvoiceData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, inv := range s.invoices {
		if inv.InvoiceID == id {
			copyInv := inv
			return &copyInv, true
		}
	}
	return nil, false
}
