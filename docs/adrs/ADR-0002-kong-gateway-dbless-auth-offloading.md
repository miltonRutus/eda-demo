# ADR-0002: Adopción de Kong API Gateway en modo DB-less con Auth Offloading

- **Estado:** Aceptado
- **Fecha:** 2026-09-12
- **Área:** Gateway / Seguridad

---

## Contexto
Tanto FastAPI como Go y futuros microservicios requerían proteger sus endpoints y gestionar la autenticación de usuarios. Implementar validación de tokens JWT, control de sesiones y CORS en cada microservicio generaba duplicación de código, riesgo de inconsistencias de seguridad y mantenimiento pesado.

---

## Opciones Consideradas
1. **Kong API Gateway en modo DB-less con Auth Offloading y Token Stripping** — Kong centraliza la seguridad perimetral, valida credenciales e inyecta cabeceras de identidad de confianza (`X-User-Id`), eliminando tokens crudos sensibles. ✅
2. **Autenticación Descentralizada en cada Microservicio** — Cada API valida JWTs con su propia librería. (Descartado por acoplamiento, duplicación de lógica y dispersión de secretos).
3. **Kong con Base de Datos PostgreSQL** — Despliegue de Kong respaldado por una base de datos relacional. (Descartado por sobrecarga operativa innecesaria para la topología estática y versionada del proyecto).

---

## Decisión
Se adopta **Kong API Gateway en modo DB-less (`KONG_DATABASE=off`)** configurado declarativamente mediante `services/kong/kong.yml` bajo el formato `_format_version: "3.0"`. Las APIs downstream operan bajo el principio de **Cero Auth**: Kong intercepta, valida credenciales perimetralmente, remueve los headers sensibles (*token stripping*) e inyecta `X-User-Id` y roles asociados hacia los microservicios internos.

---

## Consecuencias
**Positivo:**
- Código de microservicios limpio y enfocado 100% en la lógica de negocio.
- Seguridad centralizada, mitigación DoS y CORS unificado en un solo punto perimetral.
- Configuración inmutable versionada en Git y reproducible en CI/CD.

**Negativo / Compromisos:**
- Las APIs internas dependen de la red interna Docker para estar aisladas del tráfico público y garantizar la veracidad de las cabeceras `X-User-*`.

---

## Reglas derivadas
- Las APIs downstream no deben implementar middlewares de autenticación ni de CORS.
- Todo endpoint accesible externamente debe enrutarse a través de Kong Gateway (`http://localhost:8000`).
- Kong debe eliminar las cabeceras `Authorization` y `apikey` antes de reenviar la petición hacia los upstreams.
