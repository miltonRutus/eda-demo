# ADR-0003: Topología de Mensajería con Topic Exchange y Dead Letter Exchange en RabbitMQ

- **Estado:** Aceptado
- **Fecha:** 2026-09-12
- **Área:** Broker / EDA

---

## Contexto
El sistema Legacy monolítico y los microservicios modernos deben compartir eventos de negocio de manera no bloqueante. Inicialmente se consideró un esquema simple de fanout o colas punto a punto, pero resultaba insuficiente para permitir suscripciones selectivas, patrones de enrutamiento por dominio y manejo resiliente de fallos.

---

## Opciones Consideradas
1. **Topic Exchange (`sistema.eventos.bus`) con Dead Letter Exchange (`sistema.dlx`)** — Enrutamiento flexible por patrón `[dominio].[entidad].[accion]` con confirmación manual de consumo (`manual_ack`) y tolerancia a fallos. ✅
2. **Fanout Exchange Simple** — Difunde cada mensaje a todas las colas sin capacidad de filtrado por tópicos. (Descartado por falta de escalabilidad y sobrecarga innecesaria en consumidores no interesados).
3. **Punto a Punto Direct Exchange** — Requiere conocer de antemano el destino exacto de cada mensaje. (Descartado por acoplamiento implícito).

---

## Decisión
Se implementa un **Topic Exchange durable (`sistema.eventos.bus`)** como bus principal de eventos y un **Dead Letter Exchange (`sistema.dlx`)** para descarte seguro. Las routing keys siguen estrictamente el formato canónico `[dominio].[entidad].[accion]`. Los productores publican de forma persistente (`delivery_mode=2`) bajo el paradigma *Fire and Forget*, y los consumidores aplican confirmación manual (`ack`) e idempotencia por `event_id`.

---

## Consecuencias
**Positivo:**
- Desacoplamiento total: nuevos servicios pueden escuchar eventos existentes sin modificar el emisor.
- Resiliencia y persistencia ante caídas temporales de consumidores.
- Trazabilidad y recuperación de mensajes no procesables en la cola de mensajes muertos (DLQ).

**Negativo / Compromisos:**
- Los consumidores deben gestionar explícitamente la deduplicación de mensajes para asegurar idempotencia (garantía *at-least-once* de RabbitMQ).

---

## Reglas derivadas
- Todo evento debe respetar el formato JSON canónico con `event_id` (UUID), `correlation_id`, `timestamp`, `type` y `data`.
- Las colas deben configurarse obligatoriamente con el argumento `x-dead-letter-exchange: sistema.dlx`.
- Prohibido el uso de `auto_ack = true` en los consumidores.
