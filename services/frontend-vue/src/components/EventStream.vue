<script setup>
defineProps({
  events: {
    type: Array,
    required: true,
  },
});

const emit = defineEmits(['clear-events']);

function getEventColor(event) {
  if (event.event?.includes('loyalty') || event.type?.includes('loyalty')) return 'accent';
  if (event.event?.includes('facturas') || event.type?.includes('facturas')) return 'secondary';
  if (event.event?.includes('legacy') || event.type?.includes('legacy')) return 'warning';
  return 'info';
}

function getEventIcon(event) {
  if (event.event?.includes('loyalty')) return 'mdi-star-circle';
  if (event.event?.includes('facturas')) return 'mdi-file-document-check';
  if (event.event?.includes('connection')) return 'mdi-connection';
  return 'mdi-broadcast';
}
</script>

<template>
  <v-card class="pa-4 rounded-lg h-100" color="surface" elevation="2">
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-icon color="info" class="mr-2">mdi-console-line</v-icon>
        <span class="text-h6 font-weight-bold">Consola de Eventos en Tiempo Real (WebSocket)</span>
      </div>
      <div class="d-flex align-center ga-2">
        <v-chip size="small" color="primary" variant="outlined">{{ events.length }} eventos</v-chip>
        <v-btn
          icon="mdi-trash-can-outline"
          size="x-small"
          variant="text"
          color="grey"
          :disabled="events.length === 0"
          @click="emit('clear-events')"
        ></v-btn>
      </div>
    </div>

    <div v-if="events.length === 0" class="text-caption text-grey font-italic text-center py-10">
      Esperando mensajes del broker RabbitMQ a través de Kong Gateway...
    </div>

    <div v-else class="d-flex flex-column ga-2" style="max-height: 480px; overflow-y: auto;">
      <v-card
        v-for="evt in events"
        :key="evt.id"
        class="pa-3 rounded-lg bg-surface-variant event-card-enter"
        variant="flat"
      >
        <div class="d-flex justify-space-between align-center mb-1">
          <div class="d-flex align-center ga-2">
            <v-avatar :color="getEventColor(evt)" size="24" variant="tonal">
              <v-icon size="14">{{ getEventIcon(evt) }}</v-icon>
            </v-avatar>
            <span class="font-weight-bold text-body-2">{{ evt.title || evt.event }}</span>
          </div>
          <span class="text-caption text-grey">{{ evt.receivedAt }}</span>
        </div>

        <div class="text-caption text-grey-lighten-1 my-1">
          <span v-if="evt.puntos_obtenidos !== undefined">
            +{{ evt.puntos_obtenidos }} puntos acreditados al cliente {{ evt.cliente_id }} (Total: {{ evt.total_puntos }} pts - Nivel {{ evt.nivel }})
          </span>
          <span v-else-if="evt.factura_id">
            Factura <strong class="text-white">{{ evt.factura_id }}</strong> emitida por ${{ evt.total?.toFixed(2) }} (Almacén: {{ evt.almacen_asignado }})
          </span>
          <span v-else-if="evt.message">
            {{ evt.message }}
          </span>
        </div>

        <div class="d-flex justify-space-between align-center mt-2 pt-1 border-t-grey-darken-3 text-caption text-grey-darken-1">
          <span class="mono-font" style="font-size: 10px;">{{ evt.event || evt.type }}</span>
          <span v-if="evt.correlation_id" class="mono-font" style="font-size: 10px;">
            corr: {{ evt.correlation_id.substring(0, 16) }}...
          </span>
        </div>
      </v-card>
    </div>
  </v-card>
</template>
