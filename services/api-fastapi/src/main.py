from contextlib import asynccontextmanager
import logging
import os
from typing import Optional

from fastapi import FastAPI, HTTPException, Response, WebSocket, WebSocketDisconnect, status
from fastapi.responses import HTMLResponse, JSONResponse

from src.consumer import NotificationConsumer
from src.loyalty_service import LoyaltyService
from src.models import (
    CustomerPointsResponse,
    HealthResponse,
    LinkItem,
    ProblemDetails,
    RootDiscoveryResponse,
)
from src.websocket_manager import ConnectionManager

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("api-fastapi")

# Instancias compartidas del microservicio
loyalty_service = LoyaltyService()
ws_manager = ConnectionManager()
consumer = NotificationConsumer(loyalty_service=loyalty_service, ws_manager=ws_manager)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Ciclo de vida de FastAPI: inicia el consumidor AMQP y lo detiene al apagar el servicio."""
    logger.info("Starting up api-fastapi and AMQP notification consumer...")
    # Solo iniciar conexión si no estamos en entorno de pruebas unitarias aisladas
    if os.getenv("TESTING") != "true":
        await consumer.start()
    yield
    logger.info("Shutting down api-fastapi and AMQP notification consumer...")
    await consumer.stop()


app = FastAPI(
    title="API 1 — FastAPI Orchestrator & Real-Time WebSockets",
    description="Microservicio de orquestación, fidelización de clientes y retransmisión push vía WebSockets en EDA.",
    version="1.0.0",
    docs_url=None,  # Manejado dinámicamente para compatibilidad con Kong Gateway
    lifespan=lifespan,
)


@app.get("/docs", include_in_schema=False)
def swagger_ui():
    """Swagger UI interactivo con resolución dinámica de URL para Reverse Proxies (Kong Gateway)."""
    html_content = """<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>API 1 (FastAPI) — Swagger UI</title>
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


@app.get(
    "/",
    response_model=RootDiscoveryResponse,
    summary="Catálogo HATEOAS de descubrimiento",
    description="Retorna el catálogo raíz con hipervínculos hacia todos los recursos y la documentación OpenAPI.",
)
def root_discovery():
    return RootDiscoveryResponse(
        service="API 1 — FastAPI Orchestrator & WebSockets",
        version="1.0.0",
        description="Microservicio de fidelización, orquestación y notificaciones push en tiempo real mediante WebSockets.",
        _links={
            "self": LinkItem(href="/", method="GET"),
            "docs": LinkItem(href="/docs", method="GET"),
            "openapi": LinkItem(href="/openapi.json", method="GET"),
            "health": LinkItem(href="/health", method="GET"),
            "customers_points": LinkItem(href="/customers/{client_id}/points", method="GET"),
            "websocket": LinkItem(href="/ws/{client_id}", method="GET"),
        },
    )


@app.get(
    "/health",
    response_model=HealthResponse,
    summary="Health check y estado del microservicio",
    description="Comprueba la salud del microservicio, el estado del broker RabbitMQ y los WebSockets activos.",
)
def health_check():
    return HealthResponse(
        status="ok",
        service="api-fastapi",
        broker_connected=consumer.is_connected,
        active_websockets=ws_manager.total_connections,
        _links={
            "self": LinkItem(href="/health", method="GET"),
            "docs": LinkItem(href="/docs", method="GET"),
            "root": LinkItem(href="/", method="GET"),
        },
    )


@app.get(
    "/customers/{client_id}/points",
    response_model=CustomerPointsResponse,
    responses={
        404: {
            "model": ProblemDetails,
            "description": "Cliente no encontrado en el sistema de fidelidad (RFC 7807)",
        }
    },
    summary="Consultar puntos de fidelidad de un cliente",
    description="Obtiene el saldo acumulado, categoría de fidelidad e historial reciente de un cliente.",
)
def get_customer_points(client_id: str):
    info = loyalty_service.get_customer_info(client_id)
    if not info:
        problem = ProblemDetails(
            type="https://api.empresa.com/errors/customer-not-found",
            title="Customer Not Found",
            status=status.HTTP_404_NOT_FOUND,
            detail=f"El cliente con identificador '{client_id}' no tiene registro de compras ni puntos de fidelidad acumulados.",
            instance=f"/customers/{client_id}/points",
            _links={
                "self": LinkItem(href=f"/customers/{client_id}/points", method="GET"),
                "docs": LinkItem(href="/docs", method="GET"),
                "root": LinkItem(href="/", method="GET"),
            },
        )
        return JSONResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            content=problem.model_dump(by_alias=True),
            media_type="application/problem+json",
        )

    return CustomerPointsResponse(
        cliente_id=info["cliente_id"],
        total_puntos=info["total_puntos"],
        nivel=info["nivel"],
        historial_reciente=info["historial_reciente"],
        _links={
            "self": LinkItem(href=f"/customers/{client_id}/points", method="GET"),
            "docs": LinkItem(href="/docs", method="GET"),
            "health": LinkItem(href="/health", method="GET"),
        },
    )


@app.websocket("/ws/{client_id}")
async def websocket_endpoint(websocket: WebSocket, client_id: str):
    """Canal WebSocket bidireccional para recepción push de notificaciones EDA en tiempo real."""
    await ws_manager.connect(websocket, client_id)
    try:
        # Enviar confirmación de suscripción al cliente
        welcome_packet = {
            "event": "connection_established",
            "client_id": client_id,
            "message": f"Conexión en tiempo real activa para el cliente '{client_id}'.",
        }
        await websocket.send_json(welcome_packet)

        # Bucle de escucha pasiva para mantener viva la conexión (heartbeats/pings)
        while True:
            data = await websocket.receive_text()
            logger.debug("Received message from client '%s': %s", client_id, data)
    except WebSocketDisconnect:
        ws_manager.disconnect(websocket, client_id)
    except Exception as exc:
        logger.warning("WebSocket exception for client '%s': %s", client_id, exc)
        ws_manager.disconnect(websocket, client_id)
