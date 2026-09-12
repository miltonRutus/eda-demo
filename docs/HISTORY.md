# HISTORY.md — Bitácora Cronológica de Sesiones y Cambios

Este archivo registra de forma secuencial los cambios, decisiones y avances realizados en cada sesión de trabajo.

---

## [2026-09-12] — Definición de Arquitectura EDA, Dockerización Centralizada y Estándares Kong SOTA

### Acciones Realizadas:
1. **Estandarización de Contenerización:**
   - Se estableció que el 100% de los servicios deben estar dockerizados y orquestados mediante un único `docker-compose.yml` en la raíz del repositorio.
   - Definición de la red `eda-network`, resolución DNS interna por servicio y mapeo de puertos estandarizado hacia el host.
2. **Integración de Kong API Gateway (SOTA):**
   - Incorporación de Kong en modo declarativo (*DB-less*) con `kong.yml`.
   - Adopción formal del principio de **Auth Offloading** y **Token Stripping**: microservicios downstream operan sin auth local, confiando en cabeceras inyectadas por Kong (`X-User-Id`, `X-User-Roles`).
   - Estandarización del pipeline de plugins en 4 fases canónicas y trazabilidad distribuida con `X-Correlation-ID`.
3. **Refactorización de `README.md`:**
   - Sustitución de diagramas ASCII por diagramas interactivos en **Mermaid** (Flowchart de Arquitectura y Sequence Diagram del caso de uso real).
   - Estandarización de *routing keys* en formato canónico de tres niveles: `[dominio].[entidad].[accion]`.
4. **Adaptación del Modelo de Gobernanza de `AGENTS.md`:**
   - Incorporación de la **Regla de Oro** (flujo de ejecución obligatorio antes de actuar).
   - Estandarización de roles Senior, Conventional Commits, mentalidad TDD con patrón AAA y separación código en inglés / documentación en español.
   - Implementación del sistema formal de documentación (`TODO.md`, `LESSONS.md`, `HISTORY.md`) y gestión de **ADRs** con tres decisiones fundacionales (`ADR-0001`, `ADR-0002`, `ADR-0003`).
5. **Ejecución y Validación de la Fase 1 (Infraestructura Central):**
   - Configuración de `.env.example`, `.env` y `docker-compose.yml`.
   - Implementación de topología declarativa en `services/rabbitmq/definitions.json` (`sistema.eventos.bus`, `sistema.dlx`, colas y bindings).
   - Implementación de configuración declarativa DB-less en `services/kong/kong.yml` (upstreams, servicios, rutas `/api/v1/fastapi`, `/api/v1/go`, `/ws`, plugins `cors`, `correlation-id`, `rate-limiting`).
   - Verificación de salud: contenedores `eda-rabbitmq` y `eda-kong` levantados y en estado `running (healthy)`.
6. **Implementación de la Fase 2 (Servicio Legacy de Compras):**
   - Construcción de `services/legacy-service/Dockerfile` Multi-Stage con targets explícitos `dev` y `prod` sobre Python 3.12-slim.
   - Implementación de `POST /orders` con patrón Fire & Forget asíncrono hacia `sistema.eventos.bus` con routing key `legacy.pedidos.creado`.
   - Preservación y generación de `X-Correlation-ID` en payload y metadatos de AMQP.
   - Suite de pruebas unitarias (`pytest`) con patrón AAA validando contratos públicos y resiliencia ante caídas del broker (3/3 passing).
   - Orquestación en `docker-compose.yml` (target `dev`), arranque en caliente y verificación de encolamiento exitoso en RabbitMQ.
 7. **Formalización de ADR-0004 (Dockerfiles Multi-Stage `dev` y `prod`):**
    - Adopción mandatoria de Dockerfiles multi-stage con targets canónicos `dev` (hot-reload, testing, tooling) y `prod` (mínimo, seguro, non-root) en todos los microservicios.
 8. **Implementación de la Fase 3 (API 2 — Golang Service: Inventario y Facturación):**
    - Creación de `services/api-go/Dockerfile` multi-stage con targets `dev` (Go 1.23-alpine) y `prod` (binario estático en Alpine 3.20 con usuario `appuser`).
    - Consumidor AMQP con manual ACK (`api-go.procesamiento_inventario`), lógica de reserva de stock atómica y cálculo fiscal (19% IVA).
    - Publicación del evento `facturacion.facturas.generada` preservando el `correlation_id`.
    - Endpoints HTTP `GET /health` y `GET /invoices`.
    - Pruebas unitarias nativas (`go test -v ./...`) pasando al 100% (4/4).
    - Integración en `docker-compose.yml` y validación end-to-end a través de Kong Gateway (`http://localhost:8000/api/v1/go/*`), confirmando el consumo inmediato del evento emitido por Legacy y enrutamiento perimetral.
 9. **Adopción de Google Cloud REST Standards, HATEOAS y OpenAPI Discovery:**
    - Formalización en `AGENTS.md` y creación de `ADR-0005: Estandarización de APIs RESTful, HATEOAS y OpenAPI Discovery`.
    - Implementación de bloque canónico `_links` en todas las respuestas JSON de `legacy-service` y `api-go` (enlaces mínimos `self`, `docs`, `collection`).
    - Exposición de documentación OpenAPI 3.0 (`/docs` con Swagger UI y `/openapi.json` con especificación JSON) y catálogo raíz descubrible (`GET /`) en cada microservicio.
    - Manejo estructurado de errores (`404` con código legible, mensaje descriptivo y bloque `_links`).
    - Suites de pruebas actualizadas y pasando al 100% en contenedores (`legacy-service`: 4/4 en `pytest`, `api-go`: 8/8 en `go test`).
    - Verificación a través de Kong Gateway en los prefijos `/api/v1/legacy/*` y `/api/v1/go/*`.
 10. **Mejora Integral de AGENTS.md con Contratos EDA + REST SOTA:**
     - Estandarización de eventos en RabbitMQ bajo la especificación **CloudEvents 1.0** y Envelope Pattern (`specversion`, `id`, `source`, `type`, `time`, `correlationid`, `data`).
     - Creación de `asyncapi.yaml` (AsyncAPI 2.6) documentando la topología de canales (`legacy.pedidos.creado`, `facturacion.facturas.generada`) y esquemas de payload.
     - Estandarización de respuestas de error `4xx` y `5xx` bajo **RFC 7807 (Problem Details)** con cabecera `Content-Type: application/problem+json`.
     - Formalización de la estrategia de adaptación para el sistema Legacy (sin cambios estructurales en BD, envoltura modular con captura de errores ante indisponibilidad del broker).
     - Formalización del estándar de integración con Frontend (aislamiento total del puerto AMQP, navegación dinámica vía HATEOAS `_links` y streaming en tiempo real vía WebSockets a través de Kong).
     - Validación en vivo de extremo a extremo: orden de compra enviada con CloudEvents, consumida por Go emitiendo factura CloudEvents, y endpoint de error 404 validado con `application/problem+json` y HATEOAS.
 11. **Generación Automática de OpenAPI a partir de Modelos y Corrección de Swagger UI tras Kong Gateway:**
     - **Automatización en Go:** Implementación de **Swaggo (`swag`)** para inspeccionar automáticamente los modelos y structs fuertemente tipados (`models.InvoiceData`, `models.HealthResponse`, `models.ProblemDetails`, `models.RootDiscoveryResponse`), generando `swagger.json` y `docs.go` sin definiciones manuales.
     - **Resolución de Error de Ruta en Gateway:** Corrección del error `Failed to load API definition (Fetch error Not Found /openapi.json)`. Swagger UI ahora computa dinámicamente el prefijo del Gateway (`/api/v1/go/openapi.json` y `/api/v1/legacy/openapi.json`) a partir de `window.location.pathname`, funcionando transparentemente tanto en acceso directo como detrás del Gateway.
     - **Validación Automatizada y Visual:** Pruebas unitarias de Go (8/8) y Python (4/4) superadas; verificación en vivo mediante navegador headless de `http://localhost:8000/api/v1/go/docs` confirmando la carga limpia e interactiva de todos los endpoints y esquemas.
 12. **Implementación de la Fase 4 (API 1 — Python FastAPI: Orquestación, Fidelidad y WebSockets):**
     - **Microservicio `services/api-fastapi`:** Multi-stage `Dockerfile` (`dev`/`prod`), usuario `appuser`, FastAPI con Pydantic V2.
     - **Consumidor Asíncrono AMQP (`aio-pika`):** Escucha `api-fastapi.notificaciones_pedido` con manual ACK, procesando `legacy.pedidos.creado` para acumulación de puntos de fidelidad y `facturacion.facturas.generada` para notificación fiscal.
     - **Servidor WebSockets:** Endpoint `/ws/{client_id}` con `ConnectionManager`, administrando suscripciones de clientes y retransmisión push bidireccional.
     - **Endpoints REST HATEOAS & RFC 7807:** Catálogo descubrible (`/`), health check (`/health`), consulta de saldo (`/customers/{client_id}/points`) con código 404 estandarizado y Swagger UI interactivo adaptado a Kong en `/docs`.
     - **Testing y Orquestación:** 9/9 pruebas unitarias pasando (`pytest`), orquestación en `docker-compose.yml` (`api-fastapi`), conexión en red `eda-network` y validación end-to-end de recepción en vivo vía WebSocket de ambos eventos tras compra en Legacy.
 13. **Implementación de la Fase 5 (Frontend Reactivo con Vue.js 3 + Vite + Vuetify):**
     - **Arquitectura Smart / Dumb Components en `services/frontend-vue`:**
       - Smart Container [`DashboardView.vue`](file:///home/eabol/coding/eda-demo/services/frontend-vue/src/views/DashboardView.vue) integrando stores Pinia y conexión WebSocket a través de Kong Gateway.
       - Dumb Components: [`OrderForm.vue`](file:///home/eabol/coding/eda-demo/services/frontend-vue/src/components/OrderForm.vue) (catálogo interactivo de productos), [`LoyaltyCard.vue`](file:///home/eabol/coding/eda-demo/services/frontend-vue/src/components/LoyaltyCard.vue) (puntos, tiers y barra de progreso), [`InvoiceList.vue`](file:///home/eabol/coding/eda-demo/services/frontend-vue/src/components/InvoiceList.vue) (facturas y almacén), [`EventStream.vue`](file:///home/eabol/coding/eda-demo/services/frontend-vue/src/components/EventStream.vue) (consola de eventos en tiempo real con chips de color), [`ArchitecturePipeline.vue`](file:///home/eabol/coding/eda-demo/services/frontend-vue/src/components/ArchitecturePipeline.vue) (diagrama de flujo interactivo con pulsos luminosos en nodos activos).
     - **Stores en Pinia:** `orderStore` (compras REST), `webSocketStore` (streaming en vivo con reconexión), `loyaltyStore` (puntos reactivos) e `invoiceStore` (facturación fiscal).
     - **Multi-stage Dockerfile:** `target: dev` (Vite dev server en puerto 5173 con hot-reload) y `target: prod` (Nginx Alpine inmutable).
     - **Testing Automatizado:** 5/5 pruebas unitarias con Vitest y `@vue/test-utils` pasando en contenedor Docker (`docker compose run --rm frontend-vue npm test`).
     - **Verificación End-to-End con Subagente de Navegador:** Navegación en vivo a `http://localhost:5173/`, confirmación de socket conectado a Kong (`Kong Proxy Conectado`), simulación de compra por $299.98 y actualización automática en tiempo real: +29 puntos de fidelidad (96 pts totales, Silver) y nueva factura emitida (`INV-6e7db467` en `WH-SOUTH-02`) con captura de pantalla registrada.
 14. **Elaboración de la Guía de Demostración Ejecutiva y Técnica ([`docs/DEMO_GUIDE.md`](./docs/DEMO_GUIDE.md)):**
     - Creación de un guión de presentación en 7 actos diseñado para reuniones de equipo y stakeholders técnicos.
     - Documentación de pruebas en vivo del ciclo EDA completo, inspección en RabbitMQ Management Console (`15672`) y prueba de resiliencia ante parada intencional del broker (`docker compose stop rabbitmq`), demostrando la continuidad de transacciones en el sistema Legacy.
     - Actualización de `README.md` con enlaces directos a las 3 instancias de Swagger UI, la guía de demo y comandos de ejecución de las 26 pruebas unitarias contenerizadas.
