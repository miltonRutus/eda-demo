# Architectural Decision Records (ADRs)

Este directorio contiene el registro formal de las decisiones arquitectónicas permanentes adoptadas en el proyecto, siguiendo las directrices estipuladas en [`AGENTS.md`](../../AGENTS.md).

---

## 📑 Índice de Decisiones Arquitectónicas

| Número | Título | Estado | Fecha | Área |
|---|---|---|---|---|
| [ADR-0001](ADR-0001-orquestador-docker-compose-centralizado.md) | Orquestación Centralizada con Docker Compose en la Raíz | Aceptado | 2026-09-12 | DevOps / Docker |
| [ADR-0002](ADR-0002-kong-gateway-dbless-auth-offloading.md) | Adopción de Kong API Gateway en modo DB-less con Auth Offloading | Aceptado | 2026-09-12 | Gateway / Seguridad |
| [ADR-0003](ADR-0003-rabbitmq-topic-exchange-eda.md) | Topología de Mensajería con Topic Exchange y Dead Letter Exchange en RabbitMQ | Aceptado | 2026-09-12 | Broker / EDA |
| [ADR-0004](ADR-0004-dockerfiles-multistage-dev-prod.md) | Adopción de Dockerfiles Multi-Stage con Targets dev y prod | Aceptado | 2026-09-12 | DevOps / Docker |
| [ADR-0005](ADR-0005-rest-api-hateoas-openapi-standards.md) | Estandarización de APIs RESTful, HATEOAS y OpenAPI Discovery | Aceptado | 2026-09-12 | Microservicios (FastAPI/Go) / API Design |
