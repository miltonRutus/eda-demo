# ADR-0005: Estandarización de APIs RESTful, HATEOAS y OpenAPI Discovery

- **Estado:** Aceptado
- **Fecha:** 2026-09-12
- **Área:** Microservicios (FastAPI/Go) / API Design

---

## Contexto
El ecosistema de microservicios (Python FastAPI, Go y Legacy) debe interactuar de forma interoperable con clientes web, sistemas externos y el API Gateway (Kong). Basado en las directrices de diseño de APIs REST de Google Cloud y el Nivel 3 del Modelo de Madurez de Richardson, una API verdaderamente RESTful requiere una **Interfaz Uniforme** gobernada por hipermedios, documentación interactiva estandarizada y un desacoplamiento estricto entre cliente y servidor.

Previamente, las respuestas JSON carecían de hiperenlaces dinámicos de navegación y no todos los servicios exponían catálogos de descubrimiento de sus especificaciones OpenAPI.

---

## Opciones Consideradas
1. **REST Nivel 3 con HATEOAS y OpenAPI Discovery en Cada API:** ✅
   Toda respuesta incluye hipervínculos navegables (`_links` con `self`, `docs` y recursos vinculados). Cada API expone Swagger UI (`/docs`), especificación OpenAPI (`/openapi.json`) y un endpoint raíz (`/`) como catálogo descubrible.
2. **REST Nivel 2 Plano (Sin Hipermedios):**
   Respuestas con datos planos sin enlaces. Los clientes deben conocer y acoplarse rígidamente a patrones de URL predefinidos. Se descarta por violar el principio de uniformidad e hipermedios de REST.
3. **RPC sobre HTTP / gRPC:**
   Uso de endpoints de tipo comando (`/createOrder`, `/cancelInvoice`). Se descarta porque el ecosistema se integra sobre Kong Gateway mediante REST HTTP estándar y OpenAPI 3.0.

---

## Decisión
Se establece como estándar mandatorio que **todas las APIs del repositorio** deben cumplir con las siguientes directrices:
1. **Nomenclatura basada en sustantivos plurales:** `/orders`, `/invoices`, `/health` (prohibido verbos en URIs).
2. **HATEOAS en todas las respuestas:** Bloque canónico `_links` con al menos:
   - `self`: URI del recurso actual.
   - `docs`: Enlace a la documentación interactiva OpenAPI.
   - Enlaces a recursos o acciones contextuales (`collection`, `invoices`, etc.).
3. **Catálogo Raíz (`GET /`):** Retorna el mapa navegable de recursos y documentación del microservicio.
4. **OpenAPI 3.0 Discovery:** Toda API expone `/docs` (Swagger UI) y `/openapi.json` (OpenAPI spec 3.0).
5. **Manejo uniforme de errores:** Respuestas de error estructuradas con código de error, mensaje descriptivo y bloque `_links`.

---

## Consecuencias
**Positivo:**
- Clientes y Frontend navegan la API dinámicamente sin URLs hardcodeadas.
- Documentación viva y accesible de forma uniforme en cualquier microservicio.
- Mayor mantenibilidad, estandarización y compatibilidad con herramientas de prueba y gateways.

**Negativo / Compromisos:**
- Requiere estructurar modelos de respuesta y DTOs con el campo `_links` en Go y Python.

---

## Reglas derivadas
- Todo nuevo endpoint o modelo de respuesta debe contemplar el bloque `_links`.
- En todos los microservicios debe estar implementada la ruta `/docs` y `/openapi.json`.
