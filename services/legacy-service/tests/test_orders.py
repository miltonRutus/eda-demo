from unittest.mock import patch
from fastapi.testclient import TestClient

from src.main import app

client = TestClient(app)


def test_root_discovery():
    response = client.get("/")
    assert response.status_code == 200
    data = response.json()
    assert "_links" in data
    assert data["_links"]["docs"]["href"] == "/docs"
    assert data["_links"]["openapi"]["href"] == "/openapi.json"
    assert data["_links"]["orders"]["href"] == "/orders"


def test_health_check():
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["service"] == "legacy-service"
    assert "_links" in data
    assert data["_links"]["docs"]["href"] == "/docs"


@patch("src.main.publish_order_created_event")
def test_create_order_success(mock_publish):
    mock_publish.return_value = True

    order_payload = {
        "cliente_id": "cli-442",
        "items": [
            {"sku": "PROD-A", "cantidad": 2, "precio": 149.99},
            {"sku": "PROD-B", "cantidad": 1, "precio": 50.00},
        ],
    }

    headers = {"X-Correlation-ID": "test-correlation-123"}
    response = client.post("/orders", json=order_payload, headers=headers)

    assert response.status_code == 200
    data = response.json()

    assert data["status"] == "success"
    assert data["cliente_id"] == "cli-442"
    assert data["total"] == 349.98
    assert data["correlation_id"] == "test-correlation-123"
    assert data["event_published"] is True
    assert "X-Correlation-ID" in response.headers
    assert response.headers["X-Correlation-ID"] == "test-correlation-123"
    assert "_links" in data
    assert data["_links"]["docs"]["href"] == "/docs"
    assert data["_links"]["self"]["href"] == "/orders"

    mock_publish.assert_called_once()
    call_args = mock_publish.call_args[0]
    event_payload = call_args[0]
    assert event_payload["type"] == "legacy.pedidos.creado"
    assert event_payload["data"]["total"] == 349.98
    assert len(event_payload["data"]["items"]) == 2


@patch("src.main.publish_order_created_event")
def test_create_order_broker_failure_resilience(mock_publish):
    # Simula falla en RabbitMQ (Fire and Forget)
    mock_publish.return_value = False

    order_payload = {
        "cliente_id": "cli-99",
        "items": [
            {"sku": "PROD-C", "cantidad": 1, "precio": 10.00},
        ],
    }

    response = client.post("/orders", json=order_payload)

    # La respuesta sigue siendo 200 OK porque el pedido se persistió en la BD del Legacy
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "success"
    assert data["event_published"] is False
