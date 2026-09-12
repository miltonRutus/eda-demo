# ADR-0004: Adopción de Dockerfiles Multi-Stage con Targets dev y prod

- **Estado:** Aceptado
- **Fecha:** 2026-09-12
- **Área:** DevOps / Docker

---

## Contexto
Los microservicios del ecosistema (Python FastAPI, Golang y Vue.js) requieren características diametralmente opuestas según el ciclo de vida:
- **Desarrollo local:** Requiere herramientas de compilación, recarga en caliente (*hot-reload*), utilidades de depuración y suites de pruebas automatizadas (`pytest`, `vitest`, `go test`).
- **Producción:** Exige imágenes inmutables, con una superficie de ataque mínima, tamaños reducidos (usando Alpine, Distroless o Scratch), sin herramientas de compilación y ejecutadas estrictamente bajo usuarios sin privilegios *root*.

Un Dockerfile plano o monolítico forzaría a empaquetar herramientas de desarrollo en imágenes de producción o dificultaría la ejecución ágil de pruebas dentro de contenedores.

---

## Opciones Consideradas
1. **Dockerfiles Multi-Stage con targets canónicos `dev` y `prod` en un solo archivo:** ✅
   Aprovecha el cache de capas de Docker, centraliza la definición en un único `Dockerfile` por microservicio y permite a Compose y CI/CD seleccionar el target apropiado.
2. **Archivos separados (`Dockerfile.dev` y `Dockerfile.prod`):**
   Duplica la declaración de dependencias base y genera deriva de configuración entre ambos entornos.
3. **Imagen única para ambos entornos:**
   Genera imágenes sobredimensionadas con vulnerabilidades potenciales en producción o elimina capacidades de hot-reload y testing local.

---

## Decisión
Se establece como estándar obligatorio que **cada microservicio en `services/` debe implementar un único `Dockerfile` multi-stage con dos targets explícitos**:
1. **Target `dev`:**
   - Incluye código completo o bind-mounts para desarrollo.
   - Herramientas completas para ejecución de pruebas unitarias (`pytest`, `go test`, `vitest`).
   - Comando por defecto orientado a recarga en caliente (*hot-reload* o watch mode).
   - Especificado en `./docker-compose.yml` mediante `build: { context: ..., target: dev }`.
2. **Target `prod`:**
   - Compilación estática u optimizada (ej. `CGO_ENABLED=0` en Go, `npm run build` con Nginx en Vue, dependencias sin tests en Python).
   - Ejecución bajo usuario no-root (`appuser`, UID 1000).
   - Imagen final ligera (Alpine o Distroless).

---

## Consecuencias
**Positivo:**
- Imágenes de producción seguras, reducidas y de inicio ultra-rápido.
- Entorno de desarrollo local consistente donde los comandos de test se ejecutan de manera estándar (`docker compose run --rm <service> <test-command>`).
- Mantenimiento centralizado en un único archivo por servicio.

**Negativo / Compromisos:**
- Requiere mayor disciplina en la sintaxis de Dockerfiles para no romper las capas de caché de `base`.

---

## Reglas derivadas
- Todo `Dockerfile` nuevo en `services/` debe definir los targets `dev` y `prod`.
- En `./docker-compose.yml` todo servicio debe indicar `target: dev` en su sección de `build`.
- Queda prohibido ejecutar contenedores en producción bajo el usuario `root`.
