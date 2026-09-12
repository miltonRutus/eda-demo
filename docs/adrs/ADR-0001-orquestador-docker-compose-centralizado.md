# ADR-0001: Orquestación Centralizada con Docker Compose en la Raíz

- **Estado:** Aceptado
- **Fecha:** 2026-09-12
- **Área:** DevOps / Docker

---

## Contexto
El ecosistema consta de múltiples microservicios heterogéneos (Python FastAPI, Go, Legacy monolito, Vue.js, Kong API Gateway y RabbitMQ). Ejecutar cada servicio en el entorno local del desarrollador generaba divergencias de versiones, dependencias conflictivas y dificultades para reproducir la comunicación inter-servicio y la topología de red.

---

## Opciones Consideradas
1. **Docker Compose Único en el Root** — Definición declarativa centralizada en `./docker-compose.yml` con red interna dedicada y Dockerfiles por servicio. ✅
2. **Docker Compose Múltiple por Subdirectorio** — Archivos compose separados en cada servicio con redes externas compartidas. (Descartado por fricción de arranque y dispersión de configuración).
3. **Ejecución Directa en el Host (Local Binaries)** — Cada desarrollador instala Python, Go, Node y RabbitMQ en su máquina. (Descartado por falta de reproducibilidad y problemas de dependencias).

---

## Decisión
Se decide adoptar un **único archivo `docker-compose.yml` en la raíz del repositorio** como el orquestador absoluto de todo el ciclo de vida del ecosistema. Todos los servicios están contenerizados y se comunican a través de la red puente interna `eda-network` mediante nombres de servicio DNS de Docker.

---

## Consecuencias
**Positivo:**
- Cero dependencias binarias en la máquina anfitriona (solo Docker y Docker Compose).
- Arranque y detención deterministas de todo el stack con un único comando (`docker compose up --build`).
- Resolución DNS inter-servicio aislada y homogénea idéntica a entornos de producción.

**Negativo / Compromisos:**
- Mayor consumo inicial de memoria RAM y almacenamiento por imágenes de contenedor.
- Requiere configurar bind-mounts explícitos para permitir recarga en caliente (*hot-reload*) durante el desarrollo local.

---

## Reglas derivadas
- Queda prohibida la ejecución de servicios directamente en el host para el ciclo de desarrollo o testing.
- Todo nuevo servicio debe incluir su propio `Dockerfile` y registrarse en `./docker-compose.yml`.
- Los servicios no deben utilizar `localhost` para comunicarse entre sí dentro de la red Docker.
