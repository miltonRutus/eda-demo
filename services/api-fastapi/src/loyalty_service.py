from datetime import datetime, timezone
import logging
from typing import Dict, List, Optional

logger = logging.getLogger("api-fastapi.loyalty")


class LoyaltyService:
    """Servicio de cálculo y gestión de puntos de fidelidad de clientes."""

    def __init__(self):
        self._balances: Dict[str, int] = {}
        self._history: Dict[str, List[dict]] = {}

    @staticmethod
    def calculate_tier(points: int) -> str:
        """Determina la categoría de fidelidad del cliente según sus puntos acumulados."""
        if points > 150:
            return "Gold"
        if points >= 50:
            return "Silver"
        return "Standard"

    def add_order_points(
        self, client_id: str, order_id: str, total_amount: float, correlation_id: str
    ) -> dict:
        """Calcula y acredita los puntos obtenidos por una compra (1 punto por cada $10 gastados)."""
        points_earned = int(total_amount // 10)
        current_balance = self._balances.get(client_id, 0)
        new_balance = current_balance + points_earned
        self._balances[client_id] = new_balance

        history_item = {
            "pedido_id": order_id,
            "puntos_obtenidos": points_earned,
            "monto_compra": round(total_amount, 2),
            "fecha": datetime.now(timezone.utc).isoformat(),
            "correlation_id": correlation_id,
        }

        if client_id not in self._history:
            self._history[client_id] = []
        self._history[client_id].insert(0, history_item)

        tier = self.calculate_tier(new_balance)

        logger.info(
            "Loyalty points updated for client '%s': +%d pts from order %s (Total: %d, Tier: %s)",
            client_id,
            points_earned,
            order_id,
            new_balance,
            tier,
        )

        return {
            "cliente_id": client_id,
            "pedido_id": order_id,
            "puntos_obtenidos": points_earned,
            "total_puntos": new_balance,
            "nivel": tier,
            "correlation_id": correlation_id,
        }

    def get_customer_info(self, client_id: str) -> Optional[dict]:
        """Obtiene la información de fidelidad de un cliente si existe."""
        if client_id not in self._balances and client_id not in self._history:
            return None

        balance = self._balances.get(client_id, 0)
        return {
            "cliente_id": client_id,
            "total_puntos": balance,
            "nivel": self.calculate_tier(balance),
            "historial_reciente": self._history.get(client_id, [])[:10],
        }
