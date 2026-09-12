package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"eda-demo/api-go/internal/handler"
	"eda-demo/api-go/internal/models"
	"eda-demo/api-go/internal/service"
)

func TestRootDiscovery_HATEOAS(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// ==========================================
	// Act
	// ==========================================
	mux.ServeHTTP(w, req)

	// ==========================================
	// Assert
	// ==========================================
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	links, ok := resp["_links"].(map[string]interface{})
	if !ok {
		t.Fatal("expected _links map in root discovery response")
	}

	docsLink, ok := links["docs"].(map[string]interface{})
	if !ok || docsLink["href"] != "/docs" {
		t.Errorf("expected docs link to be /docs, got %v", docsLink)
	}

	openapiLink, ok := links["openapi"].(map[string]interface{})
	if !ok || openapiLink["href"] != "/openapi.json" {
		t.Errorf("expected openapi link to be /openapi.json, got %v", openapiLink)
	}
}

func TestHealthCheck_HATEOAS(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// ==========================================
	// Act
	// ==========================================
	mux.ServeHTTP(w, req)

	// ==========================================
	// Assert
	// ==========================================
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	links, ok := resp["_links"].(map[string]interface{})
	if !ok {
		t.Fatal("expected _links map in health response")
	}
	if links["docs"] == nil {
		t.Error("expected docs link in _links")
	}
}

func TestOpenAPISpec(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()

	// ==========================================
	// Act
	// ==========================================
	mux.ServeHTTP(w, req)

	// ==========================================
	// Assert
	// ==========================================
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var spec map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&spec); err != nil {
		t.Fatalf("failed to decode OpenAPI JSON: %v", err)
	}
	if spec["swagger"] == nil && spec["openapi"] == nil {
		t.Errorf("expected swagger or openapi field, got %v", spec)
	}
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok || len(paths) == 0 {
		t.Errorf("expected non-empty paths in specification")
	}
}

func TestGetInvoice_FoundAndNotFound(t *testing.T) {
	// ==========================================
	// Arrange
	// ==========================================
	svc := service.NewInventoryBillingService()
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Add an invoice to the service
	evt, _ := svc.ProcessOrder(models.OrderCreatedEvent{
		CorrelationID: "corr-123",
		Timestamp:     time.Now().UTC(),
		Data: models.OrderData{
			OrderID:    "ORD-1",
			CustomerID: "CUST-1",
			Items:      []models.OrderItem{{SKU: "TEST", Quantity: 1, Price: 100}},
		},
	})
	invoiceID := evt.Data.InvoiceID

	// ==========================================
	// Act 1: Get existing invoice
	// ==========================================
	req1 := httptest.NewRequest(http.MethodGet, "/invoices/"+invoiceID, nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)

	// ==========================================
	// Assert 1
	// ==========================================
	if w1.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w1.Code)
	}
	var inv models.InvoiceData
	_ = json.NewDecoder(w1.Body).Decode(&inv)
	if inv.InvoiceID != invoiceID {
		t.Errorf("expected invoiceID %s, got %s", invoiceID, inv.InvoiceID)
	}
	if inv.Links["self"].Href != "/invoices/"+invoiceID {
		t.Errorf("expected self link /invoices/%s, got %s", invoiceID, inv.Links["self"].Href)
	}

	// ==========================================
	// Act 2: Get non-existing invoice
	// ==========================================
	req2 := httptest.NewRequest(http.MethodGet, "/invoices/NON-EXISTENT", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	// ==========================================
	// Assert 2
	// ==========================================
	if w2.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w2.Code)
	}
	var errResp map[string]interface{}
	_ = json.NewDecoder(w2.Body).Decode(&errResp)
	links, ok := errResp["_links"].(map[string]interface{})
	if !ok {
		t.Fatal("expected _links in 404 response")
	}
	if links["docs"] == nil {
		t.Error("expected docs link in 404 _links")
	}
}
