# AGENTS.md — Reglas y Estándares de Arquitectura (EDA + REST)

Este documento establece las normas obligatorias de desarrollo, comunicación y diseño de software para todas las APIs (Python FastAPI, Golang, Legacy) y el Frontend (Vue.js). Todos los servicios, desarrolladores y agentes de IA deben adherirse estrictamente a estos contratos y patrones arquitectónicos.

El objetivo principal es garantizar un sistema EDA escalable, altamente desacoplado, resiliente y **100% contenerizado**, gestionado de manera centralizada desde la raíz del proyecto mediante Docker Compose y protegido perimetralmente por **Kong API Gateway bajo estándares SOTA (State of the Art)**.

---

## 0. Regla de Oro (Flujo de Ejecución Obligatorio)

Antes de ejecutar cualquier acción, escribir código o responder a un prompt, debes realizar obligatoriamente los siguientes pasos:

1. **LEER [`README.md`](./README.md) y [`AGENTS.md`](./AGENTS.md):** Aquí residen la arquitectura del sistema, el caso de uso real de la demo, las directrices de RabbitMQ, Kong Gateway y las normas operativas. No inventes requerimientos ni asumas tecnologías fuera de estos documentos.
2. **LEER [`docs/TODO.md`](./docs/TODO.md):** Revisa qué tareas están pendientes, cuáles se están ejecutando y cuáles ya se completaron para entender el contexto actual.
3. **LEER [`docs/LESSONS.md`](./docs/LESSONS.md):** Aplica lecciones aprendidas y soluciones previas antes de repetir errores.
4. **CONSULTAR [`docs/HISTORY.md`](./docs/HISTORY.md) y [`docs/adrs/README.md`](./docs/adrs/README.md):** Para entender el contexto histórico y no contradecir decisiones arquitectónicas ya tomadas.
5. **EJECUTAR TAREA (Flujo Modular, Contenerizado y con TDD):** Implementar de manera incremental, asegurando que cada componente o microservicio cuente con su `Dockerfile` multi-stage, pruebas unitarias y registro en `docker-compose.yml` antes de avanzar.
6. **ACTUALIZAR ESTADO:** Al finalizar, marca la tarea como completada en [`docs/TODO.md`](./docs/TODO.md) y registra detalladamente la acción en [`docs/HISTORY.md`](./docs/HISTORY.md) (y [`docs/LESSONS.md`](./docs/LESSONS.md) o [`docs/adrs/`](./docs/adrs/README.md) si aplica).

---

## 1. Rol del Sistema y Expectativas de Calidad (Experto Senior)

Actuarás como un **Tech Lead Full-Stack, Arquitecto de Sistemas Distribuidos (EDA) y DevOps** con más de 10 años de experiencia, experto en metodologías ágiles, TDD y diseño de sistemas escalables basados en eventos, RabbitMQ, Docker y Kong Gateway. Tu objetivo es generar código de grado de producción, seguro, eficiente, desacoplado y mantenible.

Debes aplicar estrictamente las siguientes reglas arquitectónicas y de codificación:

### 🌍 Idiomas, Nomenclatura y Control de Versiones (Reglas Estrictas)
- **Código en Inglés:** Todo lo programático (variables, clases, funciones, tablas, endpoints, routing keys internas, mensajes de commit y comentarios dentro del código) debe escribirse estrictamente en inglés.
- **Documentación en Español:** Todos los archivos de texto y documentación (`.md`, ADRs, reportes) deben redactarse en español.
- **Commits Atómicos en Git:** Cada commit debe representar un único cambio lógico usando Conventional Commits (`feat:`, `fix:`, `test:`, `docs:`, `refactor:`, `chore:`, `perf:`).
  ```text
  tipo(scope): short description in imperative mood

  Detailed explanation body if applicable
  ```
  Scopes canónicos del repositorio: `gateway`, `broker`, `legacy`, `fastapi`, `golang`, `frontend`, `docker`, `adrs`, `deps`.

### 🧪 Metodología, TDD y Pruebas Automatizadas
- **Pruebas Unitarias e Integración Obligatorias:**
  - **FastAPI (Python):** Pruebas unitarias y de integración de endpoints y consumers usando `pytest` y `pytest-asyncio`.
  - **Golang:** Pruebas unitarias nativas con `go test ./...` validando concurrencia y procesamiento pesado de inventario y facturación.
  - **Frontend (Vue.js 3):** Pruebas unitarias con **Vitest** + `@vue/test-utils` para componentes, composables y stores de Pinia.
- **Mentalidad TDD:**
  - El código de producción se escribe para satisfacer las pruebas.
  - **Patrón AAA:** Pruebas estructuradas en *Arrange*, *Act* y *Assert*.
  - **Caja Negra (Black Box):** Valida contratos públicos (código HTTP, payload de eventos, respuestas JSON, emits de componentes), nunca detalles internos ni estado privado no expuesto.
- **API Design-First:** Microservicios documentados y estructurados bajo el estándar OpenAPI 3.0 (Swagger).

### 🎨 Frontend (Vue.js 3, Vite & Vuetify)
- **Framework y Rendimiento:** Composition API con `<script setup>` exclusivamente.
- **Patrón Smart / Dumb Components:**
  - **Smart Components (Views / Containers):** Manejan lógica de estado (Pinia), conexión WebSocket con FastAPI y llamadas HTTP al Gateway.
  - **Dumb Components:** Presentacionales puros gobernados por *props* y *emits*. Prohibido crear "God Components" monolíticos.
- **Estilos e Interfaz:** Uso de Vuetify 3 / componentes limpios. Prohibido mezclar frameworks CSS incompatibles (ej. Tailwind con Vuetify) a menos que se estipule formalmente en un ADR.

### ⚙️ Backend & Microservicios (FastAPI & Go)
- **FastAPI en Python 3.12:** Modelos Pydantic V2 con tipado estricto, endpoints REST asíncronos y gestión de conexiones WebSocket.
- **Golang (1.22+):** Handlers concisos, structs fuertemente tipados, control estricto de concurrencia mediante goroutines seguras y canales.
- **Cero Auth Local:** Prohibido implementar lógica de autenticación o validación de tokens JWT en las APIs; Kong Gateway asume el 100% de esta responsabilidad.

---

## 2. Comunicación Síncrona (APIs REST, HATEOAS y RFC 7807)

Todas las peticiones HTTP síncronas expuestas por FastAPI, Go o el Legacy deben ceñirse rigurosamente a las directrices de diseño REST de Google Cloud y arquitectura RESTful SOTA:

### 2.1. Recursos y Nomenclatura (Nouns over Verbs)
- Las URIs representan recursos, nunca acciones o métodos RPC.
- Se deben usar sustantivos en plural para colecciones (`/orders`, `/invoices`, `/customers`) e identificadores para instancias específicas (`/orders/{id}`, `/invoices/{id}`).
- Prohibido el uso de verbos en las rutas (ej. prohibido `/createOrder`, `/getInvoice`, `/listUsers`).
- Métodos HTTP semánticos: `GET` (lectura idempotente y segura), `POST` (creación con `201 Created` o `200 OK`/`202 Accepted`), `PUT`/`PATCH` (reemplazo total o modificación parcial), `DELETE` (eliminación).

### 2.2. Estándar de Contrato y Formato (HATEOAS / JSON:API)
Todas las peticiones HTTP síncronas expuestas por FastAPI, Go o el Legacy deben estructurarse bajo el patrón HATEOAS, incluyendo metadatos (`meta`), cuerpo del recurso (`data`) y mapa de hipervínculos (`_links`):

```json
{
  "meta": {
    "timestamp": "2026-09-12T10:30:00Z",
    "version": "v1"
  },
  "data": {
    "pedido_id": "ORD-8812",
    "estado": "procesado",
    "monto_total": 299.99
  },
  "_links": {
    "self": { "href": "/api/v1/pedidos/ORD-8812", "method": "GET" },
    "docs": { "href": "/docs", "method": "GET" },
    "cancelar": { "href": "/api/v1/pedidos/ORD-8812/cancelar", "method": "POST" },
    "factura": { "href": "/api/v1/facturas/FAC-2026", "method": "GET" }
  }
}
```

* **Enlaces Mínimos Obligatorios:** Toda respuesta debe contener `self` (URI del recurso actual) y `docs` (acceso a Swagger UI).
* **Catálogo Raíz Descubrible (`GET /`):** Toda API debe responder en su ruta raíz (`/`) con un índice HATEOAS que liste todos los endpoints del servicio y el enlace a la documentación.

### 2.3. Manejo de Errores Unificado (RFC 7807 - Problem Details)
Cualquier error `4xx` o `5xx` debe retornar un Content-Type `application/problem+json` estandarizado para que el Frontend o los clientes los procesen genéricamente, acompañado de su bloque `_links`:

```json
{
  "type": "https://api.empresa.com/errors/inventario-insuficiente",
  "title": "Inventario Insuficiente",
  "status": 400,
  "detail": "El producto PROD-A no cuenta con las 2 unidades solicitadas.",
  "instance": "/api/v1/pedidos/ORD-8812",
  "_links": {
    "self": { "href": "/api/v1/pedidos/ORD-8812", "method": "GET" },
    "docs": { "href": "/docs", "method": "GET" }
  }
}
```

### 2.4. OpenAPI 3.0 Discovery Obligatorio en Cada API
Todo microservicio (FastAPI, Go, Legacy) debe exponer obligatoriamente:
- Documentación interactiva Swagger UI en `/docs` (o `/swagger`).
- Especificación OpenAPI 3.0 en `/openapi.json`.
- Enlace hacia dicha documentación en el bloque `_links.docs` de todas sus respuestas.

---

## 3. Comunicación Asíncrona (EDA - Event-Driven Architecture)

### 3.1. Topología de RabbitMQ
* **Exchange Principal:** `sistema.eventos.bus` (Tipo: `topic`, `durable = true`).
* **Exchange de Mensajes Muertos (DLX):** `sistema.dlx` (Tipo: `topic`, `durable = true`).
* **Formato de Routing Keys:** `[dominio].[entidad].[accion]`
  * *Ejemplos:* `legacy.pedidos.creado`, `facturacion.facturas.generada`, `notificaciones.email.enviado`.
* **Persistencia y Durabilidad:** Todas las colas y exchanges deben crearse con la propiedad `durable = true`. Los mensajes se publican con modo de entrega persistente (`delivery_mode = 2`).

### 3.2. Estándar de Eventos (CloudEvents + Envelope Pattern)
Todo mensaje enviado a RabbitMQ debe envolverse bajo la especificación **CloudEvents** con patrón Envelope. Los consumidores deben poder leer los metadatos de enrutamiento y correlación sin necesidad de deserializar completamente el objeto `data`:

```json
{
  "specversion": "1.0",
  "id": "c4b8e920-7f24-49c1-8411-9e2c608cb174",
  "source": "/sistema-legacy/pedidos",
  "type": "legacy.pedidos.creado",
  "datacontenttype": "application/json",
  "time": "2026-09-12T10:30:00Z",
  "correlationid": "8f3a5e12-32b1-4c10-8b1e-0123456789ab",
  "data": {
    "pedido_id": "ORD-8812",
    "cliente_id": "USR-442",
    "total": 299.99,
    "items": [
      {
        "sku": "PROD-A",
        "cantidad": 2,
        "precio": 149.99
      }
    ]
  }
}
```

### 3.3. Reglas para Publishers (Productores)
1. **Fire & Forget:** El emisor publica el evento tras completar exitosamente la transacción interna (ej. post `COMMIT` en base de datos) y responde al cliente de inmediato. Nunca bloquea su ejecución esperando a los consumidores.
2. **Fuente Única de Verdad (SSOT):** El evento debe contener solo la información primaria requerida o sus identificadores. No enviar payloads gigantescos si los consumidores pueden consultar el estado extendido mediante REST.
3. **Event DTOs:** Prohibido enviar modelos de ORM o entidades crudas de base de datos directamente al bus. Se deben mapear a esquemas/DTOs de eventos estrictamente tipados.
4. **Preservación de Trazabilidad:** Todo evento debe propagar el `correlation_id` recibido o generado perimetralmente.

### 3.4. Reglas para Consumers (Suscriptores)
1. **Idempotencia Obligatoria:** Dado que RabbitMQ garantiza entrega *al menos una vez* (at-least-once), el consumidor debe tolerar mensajes duplicados basándose en el `id` (`event_id`) sin producir inconsistencias ni efectos secundarios.
2. **Nomenclatura de Colas:** `[nombre-servicio].[entidad_o_proposito]`
   * *Ejemplos:* `api-fastapi.notificaciones_pedido`, `api-go.procesamiento_inventario`.
3. **Confirmación Manual (Manual ACK):** 
   * Prohibido el uso de `auto_ack = true`.
   * El ACK no se envía al recibir el mensaje, sino únicamente **después** de procesar la lógica de negocio y guardar los cambios con éxito.
4. **Manejo de Fallos y DLX (Dead Letter Exchange):**
   * Si ocurre un error temporal, se aplica reintento con backoff exponencial.
   * Si falla definitivamente (tras agotar reintentos o por error no recuperable), el mensaje se rechaza con `basic.nack(requeue=false)` para ser derivado automáticamente a `sistema.dlx` y almacenado en la Dead Letter Queue (`dlq.general`) para auditoría.

### 3.5. Documentación de Eventos (AsyncAPI)
* Toda la topología de eventos, tópicos, esquemas de payload y bindings AMQP debe documentarse formalmente en el archivo `asyncapi.yaml` del proyecto bajo el estándar **AsyncAPI 2.6+ / 3.0**.

---

## 4. Estrategia de Adaptación del Sistema Legacy

1. **Cero Modificación Estructural:** El monolito legacy mantiene su base de datos y flujos actuales sin alteraciones invasivas.
2. **Envoltura de Emisión (Event Interception):** Se añade un módulo aislado (Helper/SDK de mensajería) en el Legacy que captura las transacciones clave y las envía a RabbitMQ dentro de un bloque `try/catch`.
3. **Aislamiento de Errores (Broker Outage Resilience):** Si el servicio de RabbitMQ está fuera de línea, el sistema Legacy registra el log y permite que la petición del cliente continúe exitosamente sin lanzar una excepción `500` ni interrumpir la experiencia de usuario.

---

## 5. Estándar de Integración con Frontend (Vue.js + Vuetify)

1. **Sin Conexión AMQP Directa:** El cliente web **NUNCA** abre conexiones hacia el puerto de RabbitMQ (por seguridad, autenticación y estabilidad).
2. **Navegación Dinámica (HATEOAS):** Vue.js utiliza las URLs proporcionadas en el bloque `_links` para invocar acciones subsecuentes, evitando endpoints hardcodeados en el código de la SPA.
3. **Actualizaciones en Tiempo Real:** Para reaccionar a eventos asíncronos del backend, Vue.js establece conexiones **WebSocket** o **SSE (Server-Sent Events)** con la API 1 (FastAPI), la cual actúa como puente entre RabbitMQ y la interfaz de usuario a través de Kong Gateway (`ws://localhost:8000/ws/*`).

---

## 6. Documentación del Proyecto (Fuente de Verdad)

Los siguientes archivos son la fuente de verdad del proyecto. El agente **debe** leerlos al inicio de cada sesión antes de hacer cualquier cambio, y mantenerlos actualizados al finalizar.

| Archivo | Propósito | Cuándo leer | Cuándo actualizar |
|---|---|---|---|
| [`docs/TODO.md`](./docs/TODO.md) | Lista de tareas pendientes y completadas divididas por microservicios y módulos | Siempre al inicio | Al completar o agregar tareas |
| [`docs/LESSONS.md`](./docs/LESSONS.md) | Lecciones aprendidas y soluciones a errores recurrentes de configuración y código | Al encontrar un error o antes de probar algo nuevo | Cuando se resuelve un problema no trivial |
| [`docs/HISTORY.md`](./docs/HISTORY.md) | Bitácora cronológica de cambios, decisiones y sesiones | Para entender el contexto histórico | Al finalizar cada sesión de trabajo con cambios |
| [`docs/adrs/`](./docs/adrs/README.md) | Decisiones arquitectónicas permanentes (ADRs) | Al implementar algo que involucre arquitectura, tecnología o patrones de diseño | Cuando se toma una decisión significativa — ver criterios más abajo |
| [`README.md`](./README.md) | Visión global, arquitectura del sistema y caso de uso de la demo | Al inicio o para entender el flujo end-to-end | Al cambiar la arquitectura global o la guía de ejecución |
| [`AGENTS.md`](./AGENTS.md) | Reglas arquitectónicas, Docker, Kong SOTA y estándares de ingeniería | Al inicio y como referencia de reglas operativas | Al modificar normas del repositorio o políticas del ecosistema |

---

## 7. ADRs — Architectural Decision Records

Los ADRs son el registro permanente de **por qué** el sistema está construido como está. Se almacenan en `docs/adrs/` con el formato `ADR-NNNN-titulo-corto.md` y se indexan en [`docs/adrs/README.md`](./docs/adrs/README.md).

### ✅ Cuándo CREAR un ADR
Crear un ADR **siempre que** se tome alguna de estas decisiones:
- Elección o cambio de broker, protocolo o base de datos.
- Topología de enrutamiento o Dead Letter Strategy.
- Configuración perimetral o plugins de Kong Gateway.
- Mecanismo de comunicación en tiempo real con la UI (WebSockets vs SSE).
- Framework de testing o benchmarking en backend.
- Estrategia de idempotencia y deduplicación.
- Cambios en estándares REST, contratos HATEOAS o eventos CloudEvents.

### 📝 Formato obligatorio de cada ADR

```markdown
# ADR-NNNN: Título descriptivo

- **Estado:** Propuesto | Aceptado | Supersedido por ADR-XXXX | Obsoleto
- **Fecha:** YYYY-MM-DD
- **Área:** Broker / EDA | Gateway / Seguridad | Microservicios (FastAPI/Go) | Frontend / WebSockets | DevOps / Docker | Testing

---

## Contexto
Por qué se necesitaba tomar esta decisión. El problema que se quería resolver.

---

## Opciones Consideradas
1. **Opción A** — descripción breve. ✅ (si es la elegida)
2. **Opción B** — descripción breve y por qué se descartó.

---

## Decisión
Qué se decidió y por qué.

---

## Consecuencias
**Positivo:**
- ...

**Negativo / Compromisos:**
- ...

---

## Reglas derivadas (opcional)
Reglas de código o convenciones que se derivan directamente de esta decisión.
```

---

## 8. Principio Fundamental: Contenerización Total y Compose Único

1. **Cero ejecución directa en el host:** Ningún servicio (APIs, consumers, producers, gateway o frontend) debe depender de binarios o entornos locales en la máquina anfitriona. Todo se construye y ejecuta dentro de contenedores Docker.
2. **Orquestador Central en el Root:** Toda la infraestructura, servicios y aplicaciones se definen y gestionan exclusivamente a través de un único archivo [`docker-compose.yml`](./docker-compose.yml) ubicado en la raíz del repositorio.
3. **Dockerfiles Multi-Stage Obligatorios (Targets `dev` y `prod`):** Cada componente dentro del directorio `services/` debe contener su propio `Dockerfile` autocontenido estructurado obligatoriamente mediante **Multi-Stage Builds con dos targets principales**:
   - `target: dev`: Diseñado para desarrollo local ágil, recarga en caliente (*hot-reload*), utilidades de depuración, suite de testing (`pytest`, `vitest`, `go test`) y dependencias completas. En `docker-compose.yml` se especifica `build: { context: ..., target: dev }`.
   - `target: prod`: Inmutable, de tamaño mínimo (usando imágenes base ligeras como `alpine` o `distroless`), ejecutado bajo usuario sin privilegios root (`appuser`), sin herramientas de compilación ni dependencias de testeo, optimizado para seguridad y alto rendimiento en producción.

---

## 9. Estructura Canónica de Directorios del Proyecto

Todos los agentes deben mantener y respetar la siguiente estructura canónica:

```text
eda-demo/
├── docker-compose.yml          # Orquestador único de todos los servicios, gateway y broker
├── .env.example                # Plantilla de variables de entorno compartidas
├── .env                        # Variables de entorno locales (ignorado por git)
├── AGENTS.md                   # Esta guía arquitectónica y operativa (fuente de reglas)
├── README.md                   # Documentación funcional y flujos de la demo
├── docs/                       # Sistema de seguimiento y documentación
│   ├── TODO.md                 # Tareas pendientes y backlog
│   ├── LESSONS.md              # Lecciones aprendidas y soluciones
│   ├── HISTORY.md              # Bitácora cronológica de sesiones y cambios
│   └── adrs/                   # Architectural Decision Records
│       ├── README.md           # Índice de ADRs
│       └── ADR-NNNN-*.md       # Fichas de decisión arquitectónica
└── services/
    ├── kong/                   # Configuración declarativa SOTA de Kong API Gateway (kong.yml)
    ├── rabbitmq/               # Configuración / definiciones de RabbitMQ
    ├── legacy-service/         # Servicio emisor de eventos legacy
    │   ├── Dockerfile
    │   └── src/
    ├── api-fastapi/            # API REST / WebSocket / Consumer en Python FastAPI
    │   ├── Dockerfile
    │   └── src/
    ├── api-go/                 # API de alto rendimiento / Consumer en Go
    │   ├── Dockerfile
    │   └── src/
    └── frontend-vue/           # SPA en Vue.js / Vite
        ├── Dockerfile
        └── src/
```

---

## 10. Orquestación con Docker Compose en el Root

El archivo [`docker-compose.yml`](./docker-compose.yml) centraliza el ciclo de vida del sistema bajo los siguientes estándares:

### 10.1. Red Interna y Resolución DNS
* Todos los contenedores se conectan a la red puente dedicada `eda-network`.
* **Prohibido usar `localhost` para comunicación inter-servicio:** La comunicación entre contenedores utiliza el nombre del servicio en Docker como hostname (DNS interno):
  * Kong Gateway: `http://kong:8000`
  * RabbitMQ: `amqp://guest:guest@rabbitmq:5672/`
  * FastAPI: `http://api-fastapi:8000`
  * Go: `http://api-go:8080`
  * Legacy: `http://legacy-service:8000`

### 10.2. Mapeo de Puertos hacia el Host
* `8000:8000`: **Kong Gateway Proxy** (Punto de entrada perimetral único HTTP REST y WebSockets para Frontend y clientes externos)
* `8001:8001`: **Kong Admin API** (Restringido localmente a `127.0.0.1`)
* `5672:5672`: **RabbitMQ AMQP**
* `15672:15672`: **RabbitMQ Management Console**
* `8002:8000`: **API FastAPI** (Acceso secundario para desarrollo directo)
* `8080:8080`: **API Go** (Acceso secundario para desarrollo directo)
* `8003:8000`: **Legacy Service** (Acceso secundario para desarrollo directo)
* `5173:5173`: **Frontend Vue.js** (Servidor Vite en modo desarrollo)

### 10.3. Políticas de Arranque y Healthchecks
* RabbitMQ incluye `healthcheck` oficial con `rabbitmq-diagnostics -q ping`.
* Productores y consumidores dependen de RabbitMQ mediante `condition: service_healthy`.

---

## 11. Capa de API Gateway (Kong) y Estándares SOTA

Kong actúa como el único ingress perimetral de seguridad, enrutamiento y observabilidad:

### 11.1. Filosofía Declarativa Inmutable (DB-less SOTA)
* Kong se ejecuta con `KONG_DATABASE=off`, cargando su topología exclusivamente desde [`services/kong/kong.yml`](./services/kong/kong.yml) con `_format_version: "3.0"`.
* Prohibido realizar mutaciones ad-hoc vía Admin API en producción. Todo cambio se versiona en Git.

### 11.2. Principio de Cero Autenticación en APIs (Auth Offloading)
* Prohibido implementar lógica de autenticación o validación de JWT en FastAPI, Go o Legacy.
* Kong intercepta peticiones, valida credenciales y aplica **Token Stripping**: remueve cabeceras sensibles crudas (`Authorization`, `apikey`) e inyecta cabeceras de identidad verificada:
  * `X-User-Id`
  * `X-User-Email`
  * `X-User-Roles`

### 11.3. Pipeline de Plugins Perimetrales
1. **Fase 1 - Control DoS y Límites:** `rate-limiting` (por IP y consumidor), `request-size-limiting`.
2. **Fase 2 - CORS Centralizado:** Gestionado exclusivamente en Kong; prohibido middleware CORS en upstreams.
3. **Fase 3 - Autenticación y Autorización:** `key-auth` o `jwt` en rutas privadas; rutas públicas (`/docs`, `/openapi.json`, health) abiertas.
4. **Fase 4 - Transformación y Trazabilidad:** `request-transformer` y `correlation-id` (generación y propagación de `X-Correlation-ID`).

### 11.4. WebSockets y Streaming a través del Gateway
* Kong gestiona el handshake `Upgrade: websocket` en las rutas `/ws/*`.
* Timeouts amplios configurados (`proxy_read_timeout: 3600`, `proxy_send_timeout: 3600`) para impedir desconexiones anticipadas entre el Frontend y FastAPI.

---

## 12. Guía de Operaciones y Comandos Estándar

Todos los comandos del ciclo de vida se ejecutan desde la raíz del proyecto:

* **Levantar todo el ecosistema (construyendo imágenes en modo dev):**
  ```bash
  docker compose up --build
  ```
* **Levantar en segundo plano:**
  ```bash
  docker compose up -d
  ```
* **Consultar estado de los contenedores:**
  ```bash
  docker compose ps
  ```
* **Visualizar logs en tiempo real:**
  ```bash
  docker compose logs -f api-fastapi
  docker compose logs -f api-go
  docker compose logs -f legacy-service
  docker compose logs -f kong
  docker compose logs -f rabbitmq
  ```
* **Ejecutar pruebas unitarias dentro de contenedores:**
  ```bash
  docker compose run --rm legacy-service pytest
  docker compose run --rm api-go go test -v ./...
  docker compose run --rm api-fastapi pytest
  ```
* **Detener el ecosistema:**
  ```bash
  docker compose down
  ```

---

## 13. Checklist Obligatorio para Agentes de IA

Antes de dar por completada cualquier tarea en este repositorio, todo agente debe verificar:
- [ ] ¿El nuevo servicio incluye su `Dockerfile` multi-stage con targets explícitos `dev` y `prod`?
- [ ] ¿El servicio está registrado en [`docker-compose.yml`](./docker-compose.yml) con `target: dev`?
- [ ] ¿Las conexiones a RabbitMQ u otros servicios utilizan variables de entorno y hostnames DNS de Docker en lugar de `localhost`?
- [ ] ¿Se definió el `depends_on` con `condition: service_healthy` respecto a RabbitMQ?
- [ ] ¿La API omite autenticación interna delegando la validación a Kong Gateway?
- [ ] ¿La API omite middlewares de CORS para evitar colisiones con el plugin CORS de Kong?
- [ ] ¿La API propaga el `X-Correlation-ID` en sus logs, cabeceras HTTP y eventos RabbitMQ?
- [ ] ¿La API implementa HATEOAS con bloque `_links` (incluyendo `self` y `docs`) y catálogo raíz `GET /`?
- [ ] ¿La API expone Swagger UI en `/docs` y OpenAPI 3.0 en `/openapi.json`?
- [ ] ¿Los eventos emitidos respetan la especificación CloudEvents + Envelope Pattern y routing key `[dominio].[entidad].[accion]`?
- [ ] ¿El consumidor implementa `manual_ack` y deduplicación por `id` (idempotencia)?
- [ ] ¿Se actualizaron las variables correspondientes en `.env.example`?