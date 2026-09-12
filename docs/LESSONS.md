# LESSONS.md — Lecciones Aprendidas y Soluciones Recurrentes

Este archivo documenta problemas encontrados durante el desarrollo, errores de configuración comunes y las soluciones efectivas para evitar repetirlos.

---

## 🏛️ Arquitectura & RabbitMQ

### 1. Inconsistencia entre Productor y Consumidor en el Flujo Asíncrono
* **Problema:** En arquitecturas EDA con WebSocket hacia la UI, se tiende erróneamente a intentar conectar el frontend directamente al broker o a hacer que servicios de cálculo (ej. Go) abran conexiones directas al frontend.
* **Lección:** Mantener la separación estricta: Go solo publica en RabbitMQ (`facturacion.facturas.generada`). La API de orquestación y notificaciones (FastAPI) consume dicho evento y empuja la notificación a través de su canal WebSocket gestionado por Kong.

### 2. Condición de Carrera en el Arranque de Microservicios vs Broker
* **Problema:** Los servicios en Python o Go intentan conectarse a AMQP mientras RabbitMQ aún está levantando sus plugins de management o inicializando el erlang node, provocando crash loops (`connection refused`).
* **Lección:** Usar siempre `depends_on` con `condition: service_healthy` en `docker-compose.yml`, respaldado por el healthcheck nativo de RabbitMQ (`rabbitmq-diagnostics -q ping`).

---

## 🛡️ Kong API Gateway

### 1. Colisión de Cabeceras CORS entre Kong y Upstreams
* **Problema:** Si tanto Kong como FastAPI/Go implementan middlewares de CORS, el navegador recibe cabeceras `Access-Control-Allow-Origin` duplicadas y bloquea la petición por seguridad.
* **Lección:** El CORS debe ser **exclusivo de Kong Gateway**. Los microservicios downstream no deben tener middleware CORS activado.

### 2. Timeouts de Conexiones WebSocket
* **Problema:** Por defecto, los proxies HTTP cierran conexiones inactivas tras 60 segundos, interrumpiendo el flujo push del frontend.
* **Lección:** Configurar `proxy_read_timeout` y `proxy_send_timeout` en Kong con valores amplios (ej. 3600 segundos) para las rutas que gestionan WebSockets (`/ws/*`).

### 3. Resolución DNS en Kong Gateway para Upstreams Creados Progresivamente
* **Problema:** Si Kong arranca antes de que los contenedores de los upstreams (`api-go`, `api-fastapi`) existan en la red Docker, el ring-balancer interno de Kong registra `DNS_ERROR (name error 3)` y responde `503 failure to get a peer from the ring-balancer` de forma transitoria.
* **Lección:** Al crear o levantar un nuevo upstream en Compose durante desarrollo incremental, reiniciar Kong (`docker compose restart kong`) para refrescar la tabla de resolución DNS interna de los upstreams declarativos.

### 4. Rutas Absolutas vs Relativas en Swagger UI tras un Reverse Proxy (Kong Gateway)
* **Problema:** Si el bundle HTML de Swagger UI usa una ruta absoluta fija como `url: '/openapi.json'`, cuando el usuario accede a través del Gateway (ej. `http://localhost:8000/api/v1/go/docs`), el navegador solicita `http://localhost:8000/openapi.json` en lugar de la ruta con prefijo (`/api/v1/go/openapi.json`), arrojando el error `Failed to load API definition: Fetch error Not Found /openapi.json`.
* **Lección:** Calcular la URL de la especificación dinámicamente en el cliente usando el pathname actual:
  ```javascript
  const basePath = window.location.pathname.replace(/\/docs\/?$/, '');
  const specUrl = (basePath || '') + '/openapi.json';
  window.ui = SwaggerUIBundle({ url: specUrl, ... });
  ```
  Esto garantiza portabilidad total tanto en acceso directo al microservicio como detrás de cualquier Gateway o prefijo inverso.

### 5. Generación Automática de OpenAPI a partir de Modelos (Go vs Python)
* **Python (FastAPI):** Genera OpenAPI automáticamente en tiempo de ejecución a partir de modelos Pydantic (`BaseModel`).
* **Golang:** Al ser compilado y estático, el estándar SOTA es **Swaggo (`swag`)**. Lee los `structs` de Go (`models.InvoiceData`, `models.OrderItem`, etc.) y las etiquetas JSON/Swaggo para emitir automáticamente `swagger.json` y `docs.go`, eliminando el mantenimiento manual de esquemas JSON.

---

## 🐳 Docker & Multi-Stage Builds

### 1. Separación de Entornos de Desarrollo y Producción
* **Problema:** Si se usa una sola imagen para desarrollo y producción, o bien la imagen de producción queda inflada con dependencias de testing (`pytest`, `vitest`), hot-reload y compiladores, o el entorno de desarrollo carece de herramientas de depuración rápida.
* **Lección:** Todos los `Dockerfile` del proyecto deben estructurarse con dos targets explícitos:
  - `target: dev`: Contiene dependencias completas, herramientas de testeo y configuración para recarga en caliente vía bind-mounts.
  - `target: prod`: Multi-stage mínimo, binarios compilados o assets estáticos limpios sobre imágenes base ligeras (`alpine`/`distroless`), sin herramientas de compilación y ejecutado con usuario sin privilegios root.

