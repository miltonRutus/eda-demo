import asyncio
from datetime import datetime, timezone
import json
import logging
import os
from typing import Optional
import aio_pika

from src.loyalty_service import LoyaltyService
from src.websocket_manager import ConnectionManager

logger = logging.getLogger("api-fastapi.consumer")


class NotificationConsumer:
    """Consumidor asíncrono AMQP que procesa compras y facturas para fidelidad y notificaciones WebSocket."""

    def __init__(
        self,
        loyalty_service: LoyaltyService,
        ws_manager: ConnectionManager,
    ):
        self.loyalty_service = loyalty_service
        self.ws_manager = ws_manager
        self._connection: Optional[aio_pika.abc.AbstractRobustConnection] = None
        self._channel: Optional[aio_pika.abc.AbstractRobustChannel] = None
        self._task: Optional[asyncio.Task] = None
        self.is_connected = False

    async def start(self) -> None:
        """Inicia la conexión resiliente con RabbitMQ y lanza la escucha en segundo plano."""
        host = os.getenv("RABBITMQ_HOST", "rabbitmq")
        port = int(os.getenv("RABBITMQ_PORT", "5672"))
        user = os.getenv("RABBITMQ_USER", "guest")
        password = os.getenv("RABBITMQ_PASS", "guest")
        amqp_url = f"amqp://{user}:{password}@{host}:{port}/"

        max_retries = 10
        delay = 2

        for attempt in range(1, max_retries + 1):
            try:
                logger.info(
                    "Connecting to RabbitMQ (attempt %d/%d) at %s:%d...",
                    attempt,
                    max_retries,
                    host,
                    port,
                )
                self._connection = await aio_pika.connect_robust(amqp_url)
                self._channel = await self._connection.channel()
                await self._channel.set_qos(prefetch_count=20)
                self.is_connected = True
                logger.info("Successfully connected to RabbitMQ.")
                break
            except Exception as exc:
                logger.warning("RabbitMQ connection attempt %d failed: %s", attempt, exc)
                if attempt == max_retries:
                    logger.error("Could not establish connection to RabbitMQ after max retries.")
                    return
                await asyncio.sleep(delay)

        self._task = asyncio.create_task(self._consume_queue())

    async def _consume_queue(self) -> None:
        """Configura la cola y procesa los mensajes entrantes con confirmación manual."""
        try:
            queue_name = "api-fastapi.notificaciones_pedido"
            queue = await self._channel.get_queue(queue_name)
            logger.info("Listening for messages on queue '%s'...", queue_name)

            async with queue.iterator() as queue_iter:
                async for message in queue_iter:
                    await self._process_message(message)
        except asyncio.CancelledError:
            logger.info("AMQP consumer task cancelled.")
        except Exception as exc:
            logger.error("Error in AMQP consumer loop: %s", exc, exc_info=True)
            self.is_connected = False

    async def _process_message(self, message: aio_pika.abc.AbstractIncomingMessage) -> None:
        """Procesa un mensaje individual respetando CloudEvents 1.0 y manual ACK."""
        async with message.process(requeue=False):
            try:
                body = json.loads(message.body.decode("utf-8"))
                routing_key = message.routing_key or body.get("type", "")
                logger.info("Received AMQP message with routing_key '%s'", routing_key)

                # Extraer datos de la envoltura CloudEvents 1.0
                event_type = body.get("type", routing_key)
                correlation_id = (
                    body.get("correlationid")
                    or body.get("correlation_id")
                    or message.correlation_id
                    or "unknown-correlation-id"
                )
                data = body.get("data", {})

                # 1. Evento de Compra Legacy: Calcular puntos de fidelidad y notificar
                if event_type.startswith("legacy.pedidos") or routing_key.startswith("legacy.pedidos"):
                    client_id = data.get("cliente_id", "unknown-client")
                    order_id = data.get("pedido_id", "unknown-order")
                    total = float(data.get("total", 0.0))

                    result = self.loyalty_service.add_order_points(
                        client_id=client_id,
                        order_id=order_id,
                        total_amount=total,
                        correlation_id=correlation_id,
                    )

                    ws_payload = {
                        "event": "loyalty.points.updated",
                        "title": "Puntos de Fidelidad Acreditados",
                        "cliente_id": client_id,
                        "pedido_id": order_id,
                        "puntos_obtenidos": result["puntos_obtenidos"],
                        "total_puntos": result["total_puntos"],
                        "nivel": result["nivel"],
                        "correlation_id": correlation_id,
                        "timestamp": datetime.now(timezone.utc).isoformat(),
                    }
                    await self.ws_manager.send_personal_message(ws_payload, client_id)

                # 2. Evento de Facturación Go: Notificar emisión de factura en tiempo real
                elif event_type.startswith("facturacion.facturas") or routing_key.startswith("facturacion.facturas"):
                    client_id = data.get("cliente_id", "unknown-client")
                    invoice_id = data.get("factura_id", "unknown-invoice")
                    order_id = data.get("pedido_id", "unknown-order")
                    total = float(data.get("total", 0.0))
                    warehouse = data.get("almacen_asignado", "DEFAULT-WH")
                    status = data.get("estado", "facturado")

                    ws_payload = {
                        "event": "facturacion.facturas.generada",
                        "title": "Factura Fiscal Emitida",
                        "factura_id": invoice_id,
                        "pedido_id": order_id,
                        "cliente_id": client_id,
                        "total": total,
                        "almacen_asignado": warehouse,
                        "estado": status,
                        "correlation_id": correlation_id,
                        "timestamp": datetime.now(timezone.utc).isoformat(),
                    }
                    await self.ws_manager.send_personal_message(ws_payload, client_id)
                else:
                    logger.warning("Unrecognized event type: '%s'. Ignored.", event_type)

            except Exception as err:
                logger.error("Failed to process message: %s", err, exc_info=True)
                # Al salir del bloque message.process con excepción, aio-pika envía NACK (requeue=False)
                # derivando el mensaje fallido al DLX 'sistema.dlx' automáticamente.
                raise err

    async def stop(self) -> None:
        """Cierra la conexión AMQP y cancela las tareas en segundo plano limpiamente."""
        if self._task and not self._task.done():
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass

        if self._channel and not self._channel.is_closed:
            await self._channel.close()

        if self._connection and not self._connection.is_closed:
            await self._connection.close()

        self.is_connected = False
        logger.info("RabbitMQ notification consumer closed cleanly.")
