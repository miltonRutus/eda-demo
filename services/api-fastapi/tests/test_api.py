import os
import pytest
from fastapi.testclient import TestClient

# Desactivar conexión a broker real durante tests unitarios
os.environ["TESTING"] = "true"

from src.main import app, loyalty_service

client = TestClient(app)


def test_root_discovery_hateoas():
    # Arrange & Act
    response = client.get("/")

    # Assert
    assert response.status_code == 200
    data = response.json()
    assert data["service"] == "API 1 — FastAPI Orchestrator & WebSockets"
    assert "_links" in data
    assert data["_links"]["self"]["href"] == "/"
    assert data["_links"]["docs"]["href"] == "/docs"
    assert data["_links"]["openapi"]["href"] == "/openapi.json"
    assert data["_links"]["health"]["href"] == "/health"


def test_health_check_hateoas():
    # Arrange & Act
    response = client.get("/health")

    # Assert
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "api-fastapi"
    assert "_links" in data
    assert data["_links"]["self"]["href"] == "/health"
    assert data["_links"]["docs"]["href"] == "/docs"


def test_customer_points_not_found_rfc7807():
    # Arrange & Act
    response = client.get("/customers/unknown-client-999/points")

    # Assert
    assert response.status_code == 404
    assert response.headers["content-type"] == "application/problem+json"
    data = response.json()
    assert data["type"] == "https://api.empresa.com/errors/customer-not-found"
    assert data["title"] == "Customer Not Found"
    assert data["status"] == 404
    assert "unknown-client-999" in data["detail"]
    assert data["instance"] == "/customers/unknown-client-999/points"
    assert "_links" in data
    assert data["_links"]["self"]["href"] == "/customers/unknown-client-999/points"
    assert data["_links"]["docs"]["href"] == "/docs"


def test_customer_points_success():
    # Arrange: Acreditar puntos previamente
    loyalty_service.add_order_points(
        client_id="cli-user-1",
        order_id="ORD-TEST-1",
        total_amount=199.90,
        correlation_id="corr-test-1",
    )

    # Act
    response = client.get("/customers/cli-user-1/points")

    # Assert
    assert response.status_code == 200
    data = response.json()
    assert data["cliente_id"] == "cli-user-1"
    assert data["total_puntos"] == 19
    assert data["nivel"] == "Standard"
    assert len(data["historial_reciente"]) >= 1
    assert "_links" in data
    assert data["_links"]["self"]["href"] == "/customers/cli-user-1/points"
    assert data["_links"]["docs"]["href"] == "/docs"


def test_websocket_connect_and_welcome():
    # Arrange & Act
    with client.websocket_connect("/ws/cli-test-ws") as websocket:
        # Assert: Debe recibir el paquete de bienvenida inmediato
        data = websocket.receive_json()
        assert data["event"] == "connection_established"
        assert data["client_id"] == "cli-test-ws"
        assert "Conexión en tiempo real activa" in data["message"]
