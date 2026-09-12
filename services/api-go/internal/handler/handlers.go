package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"eda-demo/api-go/docs"
	"eda-demo/api-go/internal/models"
	"eda-demo/api-go/internal/service"
)

// Handler holds HTTP endpoints for api-go.
type Handler struct {
	service *service.InventoryBillingService
}

// NewHandler constructs a new Handler.
func NewHandler(svc *service.InventoryBillingService) *Handler {
	return &Handler{service: svc}
}

// RegisterRoutes sets up HTTP routes on the mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", h.RootDiscovery)
	mux.HandleFunc("GET /docs", h.SwaggerUI)
	mux.HandleFunc("GET /openapi.json", h.OpenAPISpec)
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /status", h.HealthCheck)
	mux.HandleFunc("GET /invoices", h.ListInvoices)
	mux.HandleFunc("GET /invoices/{id}", h.GetInvoice)
}

// RootDiscovery responds with the service catalogue and HATEOAS links.
// @Summary Catálogo HATEOAS y descubrimiento de la API
// @Description Retorna el catálogo navegable de hipermedios con enlaces a todos los recursos y a la documentación OpenAPI.
// @Tags Discovery
// @Produce json
// @Success 200 {object} models.RootDiscoveryResponse
// @Router / [get]
func (h *Handler) RootDiscovery(w http.ResponseWriter, r *http.Request) {
	h.setHeaders(w, r)
	w.WriteHeader(http.StatusOK)

	resp := models.RootDiscoveryResponse{
		Service:     "Go Inventory & Billing Service",
		Version:     "1.0.0",
		Description: "High-performance asynchronous inventory reservation and fiscal invoicing service in Go.",
		Links: map[string]models.Link{
			"self":     {Href: "/", Method: "GET"},
			"docs":     {Href: "/docs", Method: "GET"},
			"openapi":  {Href: "/openapi.json", Method: "GET"},
			"invoices": {Href: "/invoices", Method: "GET"},
			"health":   {Href: "/health", Method: "GET"},
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// HealthCheck responds with current health and service status.
// @Summary Health check
// @Description Comprueba la disponibilidad del servicio con enlaces HATEOAS.
// @Tags Health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.setHeaders(w, r)
	w.WriteHeader(http.StatusOK)

	resp := models.HealthResponse{
		Status:         "ok",
		Service:        "api-go",
		InventoryReady: true,
		Links: map[string]models.Link{
			"self":     {Href: "/health", Method: "GET"},
			"docs":     {Href: "/docs", Method: "GET"},
			"invoices": {Href: "/invoices", Method: "GET"},
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// ListInvoices responds with the in-memory log of generated invoices.
// @Summary Listar facturas generadas
// @Description Lista todas las facturas procesadas concurrentemente a partir de eventos de compras.
// @Tags Invoices
// @Produce json
// @Success 200 {object} models.InvoiceListResponse
// @Router /invoices [get]
func (h *Handler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	h.setHeaders(w, r)
	w.WriteHeader(http.StatusOK)

	invoices := h.service.GetRecentInvoices()
	resp := models.InvoiceListResponse{
		Status:   "success",
		Count:    len(invoices),
		Invoices: invoices,
		Links: map[string]models.Link{
			"self": {Href: "/invoices", Method: "GET"},
			"docs": {Href: "/docs", Method: "GET"},
			"root": {Href: "/", Method: "GET"},
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// GetInvoice retrieves a single invoice by its ID with HATEOAS links.
// @Summary Obtener factura por ID
// @Description Obtiene los datos fiscales y enlaces de una factura por su ID.
// @Tags Invoices
// @Produce json
// @Param id path string true "ID de la Factura (ej. FAC-12345)"
// @Success 200 {object} models.InvoiceData
// @Failure 404 {object} models.ProblemDetails
// @Router /invoices/{id} [get]
func (h *Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	h.setHeaders(w, r)
	id := r.PathValue("id")
	if id == "" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) > 2 {
			id = parts[2]
		}
	}

	invoice, found := h.service.GetInvoiceByID(id)
	if !found {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusNotFound)
		errResp := models.ProblemDetails{
			Type:     "https://api.empresa.com/errors/invoice-not-found",
			Title:    "Invoice Not Found",
			Status:   http.StatusNotFound,
			Detail:   fmt.Sprintf("The requested invoice '%s' does not exist in the system.", id),
			Instance: fmt.Sprintf("/invoices/%s", id),
			Links: map[string]models.Link{
				"self":       {Href: fmt.Sprintf("/invoices/%s", id), Method: "GET"},
				"collection": {Href: "/invoices", Method: "GET"},
				"docs":       {Href: "/docs", Method: "GET"},
			},
		}
		_ = json.NewEncoder(w).Encode(errResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(invoice)
}

// SwaggerUI renders the interactive Swagger UI HTML.
func (h *Handler) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>API 2 (Go) — OpenAPI Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
<script>
  window.onload = () => {
    // Dynamically calculate openapi.json path relative to the current gateway route prefix
    const basePath = window.location.pathname.replace(/\/docs\/?$/, '');
    const specUrl = (basePath || '') + '/openapi.json';
    window.ui = SwaggerUIBundle({
      url: specUrl,
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis],
      layout: "BaseLayout"
    });
  };
</script>
</body>
</html>`
	_, _ = w.Write([]byte(html))
}

// OpenAPISpec returns the OpenAPI specification in JSON format generated automatically by Swaggo from Go models.
func (h *Handler) OpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	doc := docs.SwaggerInfo.ReadDoc()
	_, _ = w.Write([]byte(doc))
}

func (h *Handler) setHeaders(w http.ResponseWriter, r *http.Request) {
	correlationID := r.Header.Get("X-Correlation-ID")
	if correlationID != "" {
		w.Header().Set("X-Correlation-ID", correlationID)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
}
