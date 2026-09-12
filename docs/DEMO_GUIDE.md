# 🎬 Guía de Demostración Ejecutiva y Técnica: Arquitectura EDA con Kong Gateway

Esta guía está diseñada para realizar una presentación en vivo clara, impactante y profesional ante el equipo de ingeniería, líderes técnicos y directivos. Explica paso a paso el funcionamiento del sistema, los beneficios de la arquitectura orientada a eventos (EDA) y cómo ejecutar cada escenario de prueba.

---

## 🎯 Objetivos de la Presentación

1. **Evidenciar el Desacoplamiento Real:** Mostrar cómo un sistema Legacy puede coexistir con microservicios modernos en Go y Python sin llamadas síncronas bloqueantes ni dependencias rígidas.
2. **Seguridad y Centralización con Kong Gateway (SOTA):** Demostrar el principio de *Auth Offloading* y *Cero Auth Local* en las APIs internas.
3. **Rendimiento Especializado:** Go procesando stock e impuestos a alta velocidad de forma concurrente.
4. **Experiencia de Usuario Reactiva:** Vue.js 3 actualizándose en tiempo real mediante WebSockets a través de Kong, eliminando el *polling*.
5. **Resiliencia y Tolerancia a Fallos:** Demostrar que el negocio sigue operando incluso si el broker de mensajería experimenta caídas temporales.
6. **Contenerización Total y TDD:** Sistema 100% dockerizado sin dependencias en el host con 26/26 pruebas unitarias en verde.

---

## 🗺️ Mapa de URLs del Ecosistema

| Componente | Rol en la Demo | URL de Acceso |
|---|---|---|
| **Frontend SPA (Vue 3 + Vuetify)** | Panel reactivo para disparar pedidos y ver eventos en vivo | [http://localhost:5173/](http://localhost:5173/) |
| **Kong API Gateway (Proxy)** | Ingress perimetral unificado y enrutador de WebSockets | [http://localhost:8000](http://localhost:8000) |
| **API 1 — FastAPI (Swagger Docs)** | Orquestación, programa de fidelidad y WebSockets | [http://localhost:8000/api/v1/fastapi/docs](http://localhost:8000/api/v1/fastapi/docs) |
| **API 2 — Go Service (Swagger Docs)** | Reserva de stock, cálculo fiscal y facturación | [http://localhost:8000/api/v1/go/docs](http://localhost:8000/api/v1/go/docs) |
| **API 3 — Legacy Service (Swagger Docs)**| Monolito de pedidos con emisión CloudEvents | [http://localhost:8000/api/v1/legacy/docs](http://localhost:8000/api/v1/legacy/docs) |
| **RabbitMQ Management Console** | Inspección visual de exchanges, colas y tasas de mensajes | [http://localhost:15672](http://localhost:15672) *(guest / guest)* |

---

## ⏱️ Guión de Demostración Paso a Paso (20-25 Minutos)

### Acto 0: Preparación y Arranque Limpio (1 minuto)

Antes de iniciar la presentación, levanta el ecosistema completo desde la raíz del proyecto:

```bash
docker compose up -d
```

Verifica que los 6 contenedores estén corriendo y en estado saludable:

```bash
docker compose ps
```

*Debes ver los 6 contenedores:* `eda-rabbitmq (healthy)`, `eda-kong (healthy)`, `eda-legacy-service`, `eda-api-go`, `eda-api-fastapi` y `eda-frontend-vue`.

---

### Acto 1: Introducción Arquitectónica y Perímetro Kong (3 minutos)

> **Narrativa para el equipo:**
> *"Hoy vamos a ver cómo desacoplar nuestro monolito de compras utilizando una Arquitectura Orientada a Eventos protegida por Kong Gateway bajo estándares SOTA. Ningún microservicio downstream gestiona autenticación ni CORS; Kong asume toda la seguridad perimetral e inyecta la identidad verificada y la trazabilidad distribuida mediante `X-Correlation-ID`."*

1. **Abre en el navegador:**
   - Swagger de Legacy: [http://localhost:8000/api/v1/legacy/docs](http://localhost:8000/api/v1/legacy/docs)
   - Swagger de Go: [http://localhost:8000/api/v1/go/docs](http://localhost:8000/api/v1/go/docs)
   - Swagger de FastAPI: [http://localhost:8000/api/v1/fastapi/docs](http://localhost:8000/api/v1/fastapi/docs)
2. **Puntos a destacar:**
   - Todas las APIs están descubiertas bajo el puerto unificado de Kong (`8000`).
   - Swagger UI genera dinámicamente sus rutas relativas (`/openapi.json`) adaptándose al reverse proxy.
   - En Go, la documentación OpenAPI 3.0 se genera **100% automáticamente** desde structs con `swaggo/swag`, sin archivos JSON manuales.
   - Todas las respuestas implementan **HATEOAS** (`_links.self`, `_links.docs`) y errores bajo **RFC 7807 (Problem Details)**.

---

### Acto 2: El Dashboard Reactivo y Conexión en Tiempo Real (2 minutos)

1. Abre el panel de control web: [http://localhost:5173/](http://localhost:5173/)
2. **Muestra el chip de estado:**
   - En la barra superior derecha, señala el chip verde: **`● Kong Proxy Conectado`**.
   - Explica que el navegador no abre conexiones al puerto AMQP de RabbitMQ (prohibido por seguridad); en su lugar, mantiene un túnel persistente **WebSocket con FastAPI a través de Kong Gateway** (`ws://localhost:8000/ws/cli-442`).
3. **Muestra la estructura de la pantalla:**
   - **Formulario de Compra (Legacy):** Catálogo de productos simulados.
   - **Programa de Fidelidad (FastAPI):** Saldo actual de puntos, categoría (*Standard*, *Silver*, *Gold*) y barra de progreso.
   - **Facturación & Stock (Go):** Historial de facturas fiscales emitidas y almacenes asignados.
   - **Consola de Eventos en Vivo:** Log de eventos asíncronos recibidos por streaming.
   - **Diagrama de Pipeline:** Nodos interactivos que parpadean cuando un servicio procesa un evento.

---

### Acto 3: La Demostración en Vivo del Ciclo EDA Completo (7 minutos)

> **Narrativa para el equipo:**
> *"Vamos a demostrar el ciclo completo bajo dos escenarios: primero, un usuario interactuando con el portal web; segundo, un sistema externo (B2B o ERP) enviando pedidos vía API REST a través de Kong Gateway. Observen cómo el sistema responde en milisegundos y los microservicios procesan la lógica asíncrona concurrentemente."*

#### 3.1. Escenario A: Compra desde el Dashboard Web (Vue 3)
1. En la tarjeta **"Simulador de Compra (Legacy Service)"**:
   - Selecciona productos (por ejemplo, **2 Teclados Mecánicos RGB** = $299.98).
   - Haz clic en el botón **"Disparar Compra a Monolito"**.
2. **Observa la reacción inmediata en pantalla (sin recargar):**
   - El botón responde en **~20 ms** confirmando el pedido.
   - **Tarjeta de Fidelidad (FastAPI):**
     - Los puntos suben automáticamente (+29 puntos acumulados).
     - Si superas los 50 puntos, la categoría asciende dinámicamente a **Silver**, o a **Gold** si superas los 150 puntos.
   - **Tarjeta de Facturas (Go):**
     - Aparece una nueva factura fiscal (ej. `INV-xxxxxx`) con el 19% de IVA calculado ($356.98).
     - Se asigna determinísticamente un centro de distribución (ej. `WH-SOUTH-02` o `WH-NORTH-01`).
   - **Consola de Eventos:**
     - Se registran dos eventos CloudEvents 1.0 en tiempo real:
       1. `legacy.pedidos.creado` ➔ Notificación de compra y actualización de puntos.
       2. `facturacion.facturas.generada` ➔ Factura fiscal emitida por Go.
   - **Diagrama de Arquitectura:** Los nodos de Legacy, RabbitMQ, Go y FastAPI emiten pulsos luminosos indicando el flujo del mensaje.

#### 3.2. Escenario B: Simulación de API Externa (Postman / cURL / ERP Partner)
> **Narrativa para el equipo:**
> *"¿Qué sucede si un partner B2B, un sistema de facturación externo o una app móvil crea una orden vía API REST sin interactuar con nuestra interfaz web? Kong Gateway recibe la petición perimetral, la enruta al Monolito Legacy, y nuestro Frontend en el navegador reacciona automáticamente en tiempo real gracias a los WebSockets."*

1. **Uso de la Colección Oficial de Postman:**
   - Importa en Postman el archivo [`docs/eda-demo.postman_collection.json`](./eda-demo.postman_collection.json) (también en la raíz como [`eda-demo.postman_collection.json`](../eda-demo.postman_collection.json)).
   - Contiene variables preconfiguradas: `base_url = http://localhost:8000` y `client_id = cli-442`.
   - Abre la carpeta **"1. Simulación Externa (B2B / ERP / Partner)"** y ejecuta:
     - **"Crear Compra Externa (B2B Standard)":** Simula una orden de $249.99 con cabeceras `X-Correlation-ID` y `X-User-Id`.
     - **"Crear Compra de Alto Valor (Trigger Nivel Gold)":** Simula una orden por $689.95 que acumula +68 puntos, disparando el ascenso inmediato a Nivel Gold.

2. **O ejecución directa vía cURL desde la terminal:**
   ```bash
   curl -i -X POST http://localhost:8000/api/v1/legacy/orders \
     -H "Content-Type: application/json" \
     -H "X-Correlation-ID: b2b-external-partner-9988" \
     -H "X-User-Id: external-erp-system" \
     -d '{
       "cliente_id": "cli-442",
       "items": [
         { "sku": "PROD-A", "cantidad": 1, "precio": 149.99 },
         { "sku": "PROD-B", "cantidad": 2, "precio": 50.00 }
       ]
     }'
   ```

3. **Puntos a destacar al equipo durante la llamada:**
   - **Respuesta síncrona en 15 ms:** La respuesta HTTP retorna `200 OK` con metadatos HATEOAS (`_links`) y `event_published: true`.
   - **Efecto en Vivo en la SPA abierta:** Sin tocar el teclado ni el ratón en la ventana del navegador (`http://localhost:5173`), el Frontend actualiza los puntos de fidelidad, añade la factura generada por Go y muestra los CloudEvents en la consola con el `correlation_id` externo.
   - **Trazabilidad Distribuida:** En la consola de eventos, señala cómo el `corr: b2b-external-partner-9988` viajó desde la petición externa por Kong, RabbitMQ, Go y FastAPI hasta llegar al WebSocket del navegador.

---

### Acto 4: Inspección en RabbitMQ Management Console (3 minutos)

1. Abre la consola de RabbitMQ: [http://localhost:15672](http://localhost:15672) (credenciales: `guest` / `guest`).
2. Ve a la pestaña **Exchanges**:
   - Muestra el Topic Exchange principal: `sistema.eventos.bus`.
   - Muestra el Dead Letter Exchange: `sistema.dlx`.
3. Ve a la pestaña **Queues**:
   - Muestra las dos colas activas suscritas al Topic Exchange:
     - `api-go.procesamiento_inventario` (Binding con `legacy.pedidos.creado`).
     - `api-fastapi.notificaciones_pedido` (Bindings con `legacy.pedidos.creado` y `facturacion.facturas.generada`).
   - Muestra la cola de mensajes muertos: `dlq.general`.
4. Explica el patrón **CloudEvents 1.0 + Envelope Pattern**:
   - Los consumidores no necesitan deserializar el `data` completo para conocer el origen (`source`), tipo de evento (`type`) o `correlationid`.

---

### Acto 5: La Prueba de Fuego: Resiliencia ante Caída del Broker (5 minutos)

> **Narrativa para el equipo:**
> *"Una de las mayores ventajas de la arquitectura EDA bien implementada es que la indisponibilidad de componentes no esenciales no interrumpe las ventas del negocio. Veamos qué pasa si RabbitMQ se cae por completo."*

1. **Detén el contenedor de RabbitMQ en vivo:**
   ```bash
   docker compose stop rabbitmq
   ```
2. **Regresa al Frontend:**
   - Intenta realizar una compra en el formulario.
   - **Resultado:** La compra **se completa exitosamente** (`200 OK / 201 Created`). El usuario no recibe un error `500 Internal Server Error` ni la aplicación colapsa.
3. **Muestra los logs del servicio Legacy:**
   ```bash
   docker compose logs --tail=20 legacy-service
   ```
   - Verás que el helper de mensajería capturó la indisponibilidad de RabbitMQ con un log de advertencia (`[WARN] Broker unavailable, transaction committed safely`), garantizando la continuidad de negocio.
4. **Restaura RabbitMQ:**
   ```bash
   docker compose start rabbitmq
   ```
   - El broker vuelve a estar disponible y el ecosistema recupera la comunicación asíncrona sin necesidad de reiniciar los demás servicios.

---

### Acto 6: Suite de Pruebas Automatizadas (TDD y Calidad) (3 minutos)

> **Narrativa para el equipo:**
> *"Todo el desarrollo se rigió bajo TDD con patrón AAA (Arrange, Act, Assert). Ningún desarrollador necesita tener instalado Python, Go ni Node en su máquina host; todas las pruebas corren dentro de los contenedores Docker."*

Ejecuta en una terminal los comandos de pruebas unitarias de los 4 servicios:

```bash
# 1. Pruebas Legacy (Python FastAPI / Pytest) -> 4/4 passing
docker compose run --rm legacy-service pytest

# 2. Pruebas API Go (Go test ./...) -> 8/8 passing
docker compose run --rm api-go go test -v ./...

# 3. Pruebas FastAPI (Python FastAPI / Pytest-asyncio) -> 9/9 passing
docker compose run --rm api-fastapi pytest

# 4. Pruebas Frontend (Vue 3 / Vitest) -> 5/5 passing
docker compose run --rm frontend-vue npm test
```

*Destaca el resultado:* **26/26 pruebas unitarias exitosas** cubriendo lógica de negocio, contratos HATEOAS, esquemas OpenAPI, cálculo de impuestos, tiers de lealtad y reactividad de componentes.

---

## 💡 Preguntas Frecuentes del Equipo (Cheat Sheet)

### P1: ¿Por qué el Frontend usa WebSockets hacia FastAPI y no Web STOMP directo a RabbitMQ?
* **R:** Conectar clientes web directamente a RabbitMQ expone credenciales AMQP, dificulta la autenticación perimetral y satura el broker con conexiones efímeras. FastAPI actúa como *BFF (Backend for Frontend)* y orquestador push a través de Kong Gateway.

### P2: ¿Cómo se maneja la idempotencia si RabbitMQ entrega un mensaje duplicado?
* **R:** Siguiendo la especificación CloudEvents 1.0, cada evento incluye un `id` único global (`UUIDv4`). Los consumidores registran los IDs procesados; si reciben un `id` previamente completado, confirman con `ACK` inmediatamente sin repetir transacciones.

### P3: ¿Cómo escalamos los consumidores si el volumen de compras crece 10x?
* **R:** Docker Compose o Kubernetes permite escalar réplicas de consumidores (`docker compose up --scale api-go=4 -d`). RabbitMQ distribuye automáticamente los mensajes de la cola entre las réplicas mediante el patrón *Competing Consumers* con `prefetch_count=1`.

---

## 🏁 Conclusión del Pitch

Con esta arquitectura logramos:
1. **Riesgo Cero para el Monolito:** El sistema Legacy continúa operando sin refactorizaciones invasivas.
2. **Escalabilidad Elástica:** Go y FastAPI escalan horizontalmente según su propia demanda de CPU/memoria.
3. **Observabilidad Perimetral:** Kong centraliza el ingreso, la autenticación y la trazabilidad de extremo a extremo.
4. **Experiencia de Usuario de Última Generación:** Notificaciones reactivas en tiempo real sin recargar la pantalla.
