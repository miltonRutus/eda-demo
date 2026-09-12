from typing import Dict, List, Optional
from pydantic import BaseModel, Field


class LinkItem(BaseModel):
    """Representación canónica de un hipervínculo HATEOAS."""
    href: str = Field(..., description="URI del recurso vinculado")
    method: str = Field("GET", description="Método HTTP requerido")
    rel: Optional[str] = Field(None, description="Relación semántica con el recurso actual")


class PointsHistoryItem(BaseModel):
    """Registro de acreditación de puntos de fidelidad por compra."""
    pedido_id: str = Field(..., description="Identificador único del pedido")
    puntos_obtenidos: int = Field(..., description="Puntos calculados a partir del monto")
    monto_compra: float = Field(..., description="Total de la compra en USD")
    fecha: str = Field(..., description="Marca de tiempo ISO 8601 de la transacción")
    correlation_id: str = Field(..., description="ID de correlación distribuida")


class CustomerPointsResponse(BaseModel):
    """Respuesta con el saldo y nivel de fidelidad de un cliente."""
    cliente_id: str = Field(..., description="Identificador del cliente")
    total_puntos: int = Field(..., description="Saldo acumulado de puntos de fidelidad")
    nivel: str = Field(..., description="Nivel actual del cliente (Standard, Silver, Gold)")
    historial_reciente: List[PointsHistoryItem] = Field(
        default_factory=list,
        description="Historial reciente de acreditaciones de puntos"
    )
    links: Dict[str, LinkItem] = Field(default_factory=dict, alias="_links")


class HealthResponse(BaseModel):
    """Estado de salud del microservicio y sus dependencias."""
    status: str = Field("ok", description="Estado global del servicio")
    service: str = Field("api-fastapi", description="Nombre del microservicio")
    broker_connected: bool = Field(..., description="Indica si la conexión con RabbitMQ está activa")
    active_websockets: int = Field(0, description="Número de clientes WebSocket conectados")
    links: Dict[str, LinkItem] = Field(default_factory=dict, alias="_links")


class RootDiscoveryResponse(BaseModel):
    """Catálogo raíz descubrible HATEOAS."""
    service: str = Field("API 1 — FastAPI Orchestrator & WebSockets", description="Nombre del servicio")
    version: str = Field("1.0.0", description="Versión de la API")
    description: str = Field(..., description="Descripción funcional del servicio")
    links: Dict[str, LinkItem] = Field(default_factory=dict, alias="_links")


class ProblemDetails(BaseModel):
    """Formato estándar RFC 7807 para detalles de problemas HTTP."""
    type: str = Field(..., description="URI de referencia del tipo de error")
    title: str = Field(..., description="Título corto y legible del error")
    status: int = Field(..., description="Código de estado HTTP")
    detail: str = Field(..., description="Descripción detallada de la ocurrencia")
    instance: str = Field(..., description="URI donde ocurrió el problema")
    links: Dict[str, LinkItem] = Field(default_factory=dict, alias="_links")
