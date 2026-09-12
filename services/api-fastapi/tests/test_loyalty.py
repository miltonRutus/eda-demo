import pytest
from src.loyalty_service import LoyaltyService


def test_calculate_tier():
    # Arrange & Act & Assert
    assert LoyaltyService.calculate_tier(0) == "Standard"
    assert LoyaltyService.calculate_tier(49) == "Standard"
    assert LoyaltyService.calculate_tier(50) == "Silver"
    assert LoyaltyService.calculate_tier(150) == "Silver"
    assert LoyaltyService.calculate_tier(151) == "Gold"
    assert LoyaltyService.calculate_tier(500) == "Gold"


def test_add_order_points_and_progression():
    # Arrange
    service = LoyaltyService()

    # Act: Primera compra de $299.99 -> 29 puntos
    res1 = service.add_order_points(
        client_id="cli-100",
        order_id="ORD-001",
        total_amount=299.99,
        correlation_id="corr-1",
    )

    # Assert
    assert res1["puntos_obtenidos"] == 29
    assert res1["total_puntos"] == 29
    assert res1["nivel"] == "Standard"

    # Act: Segunda compra de $300.00 -> 30 puntos (Total: 59 puntos -> Silver)
    res2 = service.add_order_points(
        client_id="cli-100",
        order_id="ORD-002",
        total_amount=300.00,
        correlation_id="corr-2",
    )

    # Assert
    assert res2["puntos_obtenidos"] == 30
    assert res2["total_puntos"] == 59
    assert res2["nivel"] == "Silver"


def test_get_customer_info_not_found():
    # Arrange
    service = LoyaltyService()

    # Act
    info = service.get_customer_info("non-existent-client")

    # Assert
    assert info is None


def test_get_customer_info_found_with_history():
    # Arrange
    service = LoyaltyService()
    service.add_order_points(
        client_id="cli-200",
        order_id="ORD-A",
        total_amount=150.00,
        correlation_id="corr-a",
    )

    # Act
    info = service.get_customer_info("cli-200")

    # Assert
    assert info is not None
    assert info["cliente_id"] == "cli-200"
    assert info["total_puntos"] == 15
    assert info["nivel"] == "Standard"
    assert len(info["historial_reciente"]) == 1
    assert info["historial_reciente"][0]["pedido_id"] == "ORD-A"
