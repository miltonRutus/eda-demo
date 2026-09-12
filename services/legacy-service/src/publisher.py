import json
import logging
import os
import pika

logger = logging.getLogger(__name__)

RABBITMQ_HOST = os.getenv("RABBITMQ_HOST", "rabbitmq")
RABBITMQ_PORT = int(os.getenv("RABBITMQ_PORT", "5672"))
RABBITMQ_USER = os.getenv("RABBITMQ_USER", "guest")
RABBITMQ_PASS = os.getenv("RABBITMQ_PASS", "guest")
EXCHANGE_NAME = "sistema.eventos.bus"


def get_connection():
    credentials = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASS)
    parameters = pika.ConnectionParameters(
        host=RABBITMQ_HOST,
        port=RABBITMQ_PORT,
        credentials=credentials,
        connection_attempts=3,
        retry_delay=2,
    )
    return pika.BlockingConnection(parameters)


def publish_order_created_event(event_payload: dict, correlation_id: str) -> bool:
    """
    Publica el evento 'legacy.pedidos.creado' hacia el Topic Exchange de RabbitMQ.
    Implementa el principio 'Fire and Forget': si RabbitMQ falla, captura la excepción,
    registra el error y permite que la respuesta al cliente continúe sin bloquearse.
    """
    connection = None
    try:
        connection = get_connection()
        channel = connection.channel()

        # Asegurar que el exchange principal esté declarado (Topic, Durable)
        channel.exchange_declare(
            exchange=EXCHANGE_NAME,
            exchange_type="topic",
            durable=True,
        )

        properties = pika.BasicProperties(
            delivery_mode=2,  # Mensaje persistente
            content_type="application/json",
            correlation_id=correlation_id,
            message_id=event_payload.get("event_id"),
        )

        channel.basic_publish(
            exchange=EXCHANGE_NAME,
            routing_key="legacy.pedidos.creado",
            body=json.dumps(event_payload),
            properties=properties,
        )

        logger.info(
            "Event 'legacy.pedidos.creado' published successfully. EventID: %s, CorrelationID: %s",
            event_payload.get("event_id"),
            correlation_id,
        )
        return True
    except Exception as exc:
        logger.error(
            "Error publishing to RabbitMQ. Order saved locally but event publication failed: %s",
            exc,
            exc_info=True,
        )
        return False
    finally:
        if connection and connection.is_open:
            try:
                connection.close()
            except Exception:
                pass
