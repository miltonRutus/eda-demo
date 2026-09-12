import json
import logging
from typing import Dict, Set
from fastapi import WebSocket

logger = logging.getLogger("api-fastapi.websocket")


class ConnectionManager:
    """Administrador en memoria de conexiones WebSocket agrupadas por cliente."""

    def __init__(self):
        self._connections: Dict[str, Set[WebSocket]] = {}

    @property
    def total_connections(self) -> int:
        return sum(len(conns) for conns in self._connections.values())

    async def connect(self, websocket: WebSocket, client_id: str) -> None:
        """Acepta y registra una nueva conexión para un cliente."""
        await websocket.accept()
        if client_id not in self._connections:
            self._connections[client_id] = set()
        self._connections[client_id].add(websocket)
        logger.info(
            "WebSocket connected for client '%s'. Active sockets for client: %d. Total: %d",
            client_id,
            len(self._connections[client_id]),
            self.total_connections,
        )

    def disconnect(self, websocket: WebSocket, client_id: str) -> None:
        """Remueve una conexión cerrada."""
        if client_id in self._connections:
            self._connections[client_id].discard(websocket)
            if not self._connections[client_id]:
                del self._connections[client_id]
        logger.info(
            "WebSocket disconnected for client '%s'. Remaining total: %d",
            client_id,
            self.total_connections,
        )

    async def send_personal_message(self, message: dict, client_id: str) -> None:
        """Envía un mensaje JSON a todas las conexiones activas de un cliente específico."""
        if client_id not in self._connections:
            logger.debug("No active WebSocket connections found for client '%s'", client_id)
            return

        payload = json.dumps(message)
        dead_connections = set()

        for conn in self._connections[client_id]:
            try:
                await conn.send_text(payload)
            except Exception as exc:
                logger.warning(
                    "Error sending message to client '%s' connection: %s", client_id, exc
                )
                dead_connections.add(conn)

        for dead_conn in dead_connections:
            self.disconnect(dead_conn, client_id)

    async def broadcast(self, message: dict) -> None:
        """Difunde un mensaje JSON a todos los clientes actualmente conectados."""
        payload = json.dumps(message)
        dead_targets = []

        for client_id, conns in self._connections.items():
            for conn in list(conns):
                try:
                    await conn.send_text(payload)
                except Exception as exc:
                    logger.warning("Error broadcasting to client '%s': %s", client_id, exc)
                    dead_targets.append((conn, client_id))

        for dead_conn, client_id in dead_targets:
            self.disconnect(dead_conn, client_id)
