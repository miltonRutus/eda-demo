package models

import "time"

// OrderItem represents a line item in an order.
type OrderItem struct {
	SKU      string  `json:"sku"`
	Quantity int     `json:"cantidad"`
	Price    float64 `json:"precio"`
}

// OrderData represents the core order payload from legacy service.
type OrderData struct {
	OrderID    string      `json:"pedido_id"`
	CustomerID string      `json:"cliente_id"`
	Total      float64     `json:"total"`
	Items      []OrderItem `json:"items"`
}

// OrderCreatedEvent is the inbound CloudEvents 1.0 event emitted by legacy-service.
type OrderCreatedEvent struct {
	SpecVersion     string    `json:"specversion,omitempty"`
	ID              string    `json:"id,omitempty"`
	Source          string    `json:"source,omitempty"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype,omitempty"`
	Time            time.Time `json:"time,omitempty"`
	CorrelationID   string    `json:"correlationid,omitempty"`
	// Backwards compatibility fields
	EventID       string    `json:"event_id,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
	Data          OrderData `json:"data"`
}

// Link represents a HATEOAS navigable hypermedia link.
type Link struct {
	Href   string `json:"href" example:"/invoices/FAC-001"`
	Method string `json:"method" example:"GET"`
	Rel    string `json:"rel,omitempty" example:"self"`
}

// InvoiceData represents the generated fiscal invoice details with HATEOAS links.
type InvoiceData struct {
	InvoiceID         string          `json:"factura_id" example:"FAC-A1B2C3"`
	OrderID           string          `json:"pedido_id" example:"ORD-XYZ123"`
	CustomerID        string          `json:"cliente_id" example:"USR-442"`
	Subtotal          float64         `json:"subtotal" example:"299.98"`
	Tax               float64         `json:"impuesto" example:"63.00"`
	Total             float64         `json:"total" example:"362.98"`
	WarehouseAssigned string          `json:"almacen_asignado" example:"MAD-01"`
	Status            string          `json:"estado" example:"facturado"`
	Items             []OrderItem     `json:"items"`
	Links             map[string]Link `json:"_links,omitempty"`
}

// InvoiceGeneratedEvent is the outbound CloudEvents 1.0 event emitted by api-go.
type InvoiceGeneratedEvent struct {
	SpecVersion     string      `json:"specversion" example:"1.0"`
	ID              string      `json:"id" example:"c4b8e920-7f24-49c1-8411-9e2c608cb174"`
	Source          string      `json:"source" example:"/api-go/billing"`
	Type            string      `json:"type" example:"facturacion.facturas.generada"`
	DataContentType string      `json:"datacontenttype" example:"application/json"`
	Time            time.Time   `json:"time"`
	CorrelationID   string      `json:"correlationid" example:"8f3a5e12-32b1-4c10-8b1e-0123456789ab"`
	// Backwards compatibility fields
	EventID   string      `json:"event_id" example:"c4b8e920-7f24-49c1-8411-9e2c608cb174"`
	Timestamp time.Time   `json:"timestamp"`
	Data      InvoiceData `json:"data"`
}

// RootDiscoveryResponse represents the discovery catalog response.
type RootDiscoveryResponse struct {
	Service     string          `json:"service" example:"Go Inventory & Billing Service"`
	Version     string          `json:"version" example:"1.0.0"`
	Description string          `json:"description" example:"High-performance asynchronous inventory reservation and fiscal invoicing service in Go."`
	Links       map[string]Link `json:"_links"`
}

// HealthResponse represents the health check status.
type HealthResponse struct {
	Status         string          `json:"status" example:"ok"`
	Service        string          `json:"service" example:"api-go"`
	InventoryReady bool            `json:"inventory_ready" example:"true"`
	Links          map[string]Link `json:"_links"`
}

// InvoiceListResponse represents the list of generated invoices with HATEOAS links.
type InvoiceListResponse struct {
	Status   string          `json:"status" example:"success"`
	Count    int             `json:"count" example:"1"`
	Invoices []InvoiceData   `json:"invoices"`
	Links    map[string]Link `json:"_links"`
}

// ProblemDetails represents an RFC 7807 compliant error response.
type ProblemDetails struct {
	Type     string          `json:"type" example:"https://api.empresa.com/errors/invoice-not-found"`
	Title    string          `json:"title" example:"Invoice Not Found"`
	Status   int             `json:"status" example:"404"`
	Detail   string          `json:"detail" example:"The requested invoice 'FAC-001' does not exist in the system."`
	Instance string          `json:"instance" example:"/invoices/FAC-001"`
	Links    map[string]Link `json:"_links"`
}
