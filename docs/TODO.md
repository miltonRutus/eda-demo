# TODO.md — Backlog y Tareas del Proyecto EDA Demo

Este archivo contiene el estado actual de las tareas del proyecto divididas por módulos y microservicios.

---

## 📌 En Progreso

- Ninguna tarea en progreso. Todo el ecosistema de la demo EDA está 100% finalizado y listo para presentación.

---

## 📋 Pendientes (Backlog / Futuras Mejoras)

- [ ] Testing end-to-end automatizado en contenedor (Playwright / Cypress).
- [ ] Pipeline CI/CD de verificación con GitHub Actions.

---

## ✅ Completadas

- [x] **Guía de Demostración Ejecutiva y Técnica ([`docs/DEMO_GUIDE.md`](./DEMO_GUIDE.md)):**
  - [x] Guión cronológico paso a paso (7 actos) para presentar ante equipo de ingeniería y gerencia.
  - [x] Storytelling de negocio y técnico explicando desacoplamiento, Auth Offloading y resiliencia.
  - [x] Demostración en vivo de tolerancia a caídas de RabbitMQ con recuperación sin pérdida.
  - [x] Cheat sheet de preguntas frecuentes y comandos de testing en contenedores.
  - [x] Colección oficial de Postman ([`eda-demo.postman_collection.json`](./eda-demo.postman_collection.json)) y ejemplos cURL para simular peticiones de APIs externas (B2B/ERP) y observar la reacción reactiva del Frontend en tiempo real.

- [x] **Fase 5: Frontend Reactivo (Vue.js 3 + Vite + Vuetify):**
  - [x] Multi-Stage `Dockerfile` (targets `dev` con Vite y `prod` con Nginx Alpine estático).
  - [x] Arquitectura Smart / Dumb Components (`DashboardView`, `OrderForm`, `LoyaltyCard`, `InvoiceList`, `EventStream`, `ArchitecturePipeline`).
  - [x] Stores reactivos en Pinia (`orderStore`, `webSocketStore`, `loyaltyStore`, `invoiceStore`).
  - [x] Conexión en tiempo real vía WebSocket a través de Kong Gateway (`ws://localhost:8000/ws/*`) y cliente REST hacia `/api/v1/*`.
  - [x] Suite de pruebas unitarias (`vitest` y `@vue/test-utils`) pasando al 100% (5/5).
  - [x] Verificación end-to-end con subagente de navegador en `http://localhost:5173/`, validando la reactividad inmediata del pipeline, fidelidad y facturación al enviar pedidos.
  - [x] Rediseño integral bajo **Material Design 3 Expressive Design** (mallas ambientales de color, tarjetas de 24px de radio con glassmorphism, chips en píldora, tipografía *Plus Jakarta Sans*, steppers táctiles y micro-animaciones interactivas).
- [x] **Fase 4: API 1 — FastAPI (Orquestación, Fidelidad y Servidor WebSockets):**
  - [x] Multi-Stage `Dockerfile` con targets `dev` y `prod` (Python 3.12-slim, non-root user `appuser`).
  - [x] Servidor WebSocket en `/ws/{client_id}` gestionando conexiones y retransmisión push.
  - [x] Consumidor asíncrono con `aio-pika` escuchando cola `api-fastapi.notificaciones_pedido`:
    - [x] Cálculo y acumulación de puntos de fidelidad en `legacy.pedidos.creado` con progresión de categorías (*Standard*, *Silver*, *Gold*).
    - [x] Retransmisión push inmediata a WebSocket en `facturacion.facturas.generada` con datos fiscales de Go.
  - [x] Endpoints REST HATEOAS (`GET /`, `GET /health`, `GET /customers/{id}/points`) y Swagger UI interactivo en `/docs`.
  - [x] Manejo de errores 404 bajo RFC 7807 (`application/problem+json`).
  - [x] Suite de pruebas unitarias (`pytest` y `pytest-asyncio`) pasando al 100% (9/9).
  - [x] Registro en `docker-compose.yml` (`target: dev`), configuración de upstream en Kong y validación end-to-end de todo el loop EDA en tiempo real.
- [x] **Fase 1: Infraestructura Central (Docker Compose, Broker y Kong Gateway):**
  - [x] Archivos `.env.example` y `.env` configurados.
  - [x] `docker-compose.yml` en la raíz con red `eda-network` y persistencia de volúmenes.
  - [x] Topología de RabbitMQ declarada (`sistema.eventos.bus`, `sistema.dlx`, colas con DLQ bindings).
  - [x] Configuración declarativa DB-less de Kong Gateway (`kong.yml`) con plugins SOTA (`cors`, `correlation-id`, `rate-limiting`).
  - [x] Ambos contenedores verificados y saludables (`running (healthy)`).
- [x] **Fase 2: Servicio Legacy (Productor de Compras):**
  - [x] Multi-Stage `Dockerfile` (targets `dev` y `prod`) basado en Python 3.12-slim y usuario sin privilegios `appuser`.
  - [x] Endpoint `POST /orders` con cálculo de total y extracción/propagación de `X-Correlation-ID`.
  - [x] Publicador resiliente Fire & Forget a RabbitMQ (`sistema.eventos.bus` con `legacy.pedidos.creado`).
  - [x] Pruebas unitarias completas con `pytest` pasando al 100% (3/3).
  - [x] Contenedor `eda-legacy-service` integrado en `docker-compose.yml` y verificado en vivo.
- [x] **Fase 3: API 2 — Golang Service (Inventario y Facturación):**
  - [x] Multi-Stage `Dockerfile` (targets `dev` y `prod` estático mínimo en Alpine con usuario no-root `appuser`).
  - [x] Consumidor AMQP con manual ACK y conexión resiliente a la cola `api-go.procesamiento_inventario`.
  - [x] Reserva de stock con asignación determinista de almacén y cálculo de impuestos (19% IVA).
  - [x] Publicador del evento `facturacion.facturas.generada` preservando `correlation_id`.
  - [x] Endpoints `/health` y `/invoices` con trazabilidad distribuida e inspección de memoria.
  - [x] Suite de pruebas unitarias (`go test ./...`) bajo patrón AAA pasando al 100% (4/4).
  - [x] Enrutamiento perimetral verificado a través de Kong Gateway (`/api/v1/go/*`).
- [x] Definición de la arquitectura EDA y reglas para productores/consumidores en `AGENTS.md`.
- [x] Estandarización de contenerización obligatoria y Compose centralizado en el root.
- [x] Integración de Kong API Gateway SOTA y Auth Offloading perimetral.
- [x] Actualización de `README.md` con diagramas interactivos en Mermaid (arquitectura y flujo de secuencia).
- [x] Adopción de la Regla de Oro, expectativas Senior, conventional commits, sistema documental y ADRs en `AGENTS.md`.
- [x] **Adopción Mandatoria de Google Cloud REST Standards, HATEOAS y OpenAPI Discovery:**
  - [x] Formalización de reglas en `AGENTS.md` y [ADR-0005](docs/adrs/ADR-0005-rest-api-hateoas-openapi-standards.md).
  - [x] Retrofit de `legacy-service`: Catálogo raíz `GET /`, HATEOAS `_links` en `/orders`, `/health` y enlace interactivo a `/docs`.
  - [x] Retrofit de `api-go`: Catálogo raíz `GET /`, Swagger UI interactivo en `/docs`, especificación en `/openapi.json`, endpoints `/invoices` y `/invoices/{id}` con HATEOAS `_links` y errores estructurados.
  - [x] Pruebas unitarias de HATEOAS y OpenAPI validadas en contenedores Docker (`pytest` 4/4 y `go test` 8/8).
  - [x] **Generación Automática de OpenAPI a partir de Modelos y Soporte Reverse-Proxy:**
    - [x] Integración de `swaggo/swag` en Go: generación 100% automática de `swagger.json` desde structs de Go (`models.InvoiceData`, `models.HealthResponse`, `models.ProblemDetails`, etc.).
    - [x] Solución de rutas relativas dinámicas en Swagger UI (`window.location.pathname`) para eliminar el error `Fetch error Not Found /openapi.json` tras Kong API Gateway tanto en Go (`/api/v1/go/docs`) como en Legacy (`/api/v1/legacy/docs`).
- [x] **Mejora Integral de AGENTS.md con Contratos EDA + REST SOTA:**
  - [x] Especificación de eventos bajo **CloudEvents 1.0** y Envelope Pattern (`specversion`, `id`, `source`, `type`, `time`, `correlationid`, `data`).
  - [x] Estandarización de errores HTTP bajo **RFC 7807 (Problem Details)** con `application/problem+json`.
  - [x] Creación de `asyncapi.yaml` (AsyncAPI 2.6) documentando la topología de canales y esquemas de eventos.
  - [x] Formalización de la estrategia de adaptación Legacy (cero modificación de BD, envoltura SDK, resiliencia ante broker caído).
  - [x] Formalización del estándar de integración Frontend (aislamiento de RabbitMQ, HATEOAS dinámico, WebSockets/SSE vía Kong).
