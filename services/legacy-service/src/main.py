from datetime import datetime, timezone
import logging
from typing import List, Optional
import uuid

from fastapi import FastAPI, Header, Response, status
from fastapi.responses import HTMLResponse
from pydantic import BaseModel, Field

from src.publisher import publish_order_created_event

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("legacy-service")

app = FastAPI(
    title="Legacy Order Service",
    description="Simulación del sistema monolítico Legacy que emite eventos Fire and Forget hacia RabbitMQ.",
    version="1.0.0",
    docs_url=None,
)


@app.get("/docs", include_in_schema=False)
def swagger_ui():
    html_content = """<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Legacy Order Service — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
<script>
  window.onload = () => {
    const basePath = window.location.pathname.replace(/\\/docs\\/?$/, '');
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
</html>"""
    return HTMLResponse(content=html_content)


class OrderItem(BaseModel):
    sku: str = Field(..., example="PROD-A")
    cantidad: int = Field(..., gt=0, example=2)
    precio: float = Field(..., gt=0, example=149.99)


class CreateOrderRequest(BaseModel):
    cliente_id: str = Field(..., example="cli-442")
    items: List[OrderItem]


class LinkItem(BaseModel):
    href: str
    method: str = "GET"
    rel: Optional[str] = None


class OrderResponse(BaseModel):
    status: str
    message: str
    pedido_id: str
    cliente_id: str
    total: float
    event_id: str
    correlation_id: str
    event_published: bool
    links: dict[str, LinkItem] = Field(default_factory=dict, alias="_links")


@app.get("/", summary="Catálogo HATEOAS y descubrimiento de la API")
def root_discovery():
    return {
        "service": "Legacy Order Service",
        "version": "1.0.0",
        "description": "Simulación del sistema monolítico Legacy que emite eventos Fire and Forget hacia RabbitMQ.",
        "_links": {
            "self": {"href": "/", "method": "GET"},
            "docs": {"href": "/docs", "method": "GET"},
            "openapi": {"href": "/openapi.json", "method": "GET"},
            "orders": {"href": "/orders", "method": "POST"},
            "health": {"href": "/health", "method": "GET"},
        },
    }


@app.get("/health", summary="Health check endpoint")
def health_check():
    return {
        "status": "ok",
        "service": "legacy-service",
        "_links": {
            "self": {"href": "/health", "method": "GET"},
            "docs": {"href": "/docs", "method": "GET"},
        },
    }


@app.post(
    "/orders",
    response_model=OrderResponse,
    status_code=status.HTTP_200_OK,
    summary="Registrar una compra en el sistema Legacy",
)
def create_order(
    payload: CreateOrderRequest,
    response: Response,
    x_correlation_id: Optional[str] = Header(None, alias="X-Correlation-ID"),
):
    # Generar o respetar Correlation ID
    correlation_id = x_correlation_id or str(uuid.uuid4())
    event_id = str(uuid.uuid4())
    pedido_id = f"ORD-{uuid.uuid4().hex[:6].upper()}"

    # Calcular total de la compra
    total = round(sum(item.cantidad * item.precio for item in payload.items), 2)

    # Simular persistencia en base de datos monolítica existente
    logger.info(
        "Order %s persisted in Monolithic Database for client %s. Total: $%.2f",
        pedido_id,
        payload.cliente_id,
        total,
    )

    # Estructura del evento estándar CloudEvents 1.0 conforme a AGENTS.md
    now_iso = datetime.now(timezone.utc).isoformat()
    event_payload = {
        "specversion": "1.0",
        "id": event_id,
        "source": "/sistema-legacy/pedidos",
        "type": "legacy.pedidos.creado",
        "datacontenttype": "application/json",
        "time": now_iso,
        "correlationid": correlation_id,
        "event_id": event_id,
        "correlation_id": correlation_id,
        "timestamp": now_iso,
        "data": {
            "pedido_id": pedido_id,
            "cliente_id": payload.cliente_id,
            "total": total,
            "items": [item.model_dump() for item in payload.items],
        },
    }

    # Publicación asíncrona no bloqueante (Fire & Forget)
    published = publish_order_created_event(event_payload, correlation_id)

    # Propagar correlation ID en cabecera de respuesta
    response.headers["X-Correlation-ID"] = correlation_id

    # HATEOAS: enlaces de navegación y recursos relacionados
    links = {
        "self": LinkItem(href="/orders", method="POST"),
        "docs": LinkItem(href="/docs", method="GET"),
        "order": LinkItem(href=f"/orders/{pedido_id}", method="GET"),
        "health": LinkItem(href="/health", method="GET"),
    }

    return OrderResponse(
        status="success",
        message="Order created successfully in Legacy DB"
        + (" and event dispatched to broker." if published else " (broker publication deferred)."),
        pedido_id=pedido_id,
        cliente_id=payload.cliente_id,
        total=total,
        event_id=event_id,
        correlation_id=correlation_id,
        event_published=published,
        _links=links,
    )
