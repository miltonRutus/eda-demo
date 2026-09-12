<script setup>
import { computed, ref } from 'vue';

const props = defineProps({
  events: {
    type: Array,
    required: true,
  },
});

const emit = defineEmits(['clear-events']);

const selectedFilter = ref('all');

const filteredEvents = computed(() => {
  if (selectedFilter.value === 'loyalty') {
    return props.events.filter(e => e.event?.includes('loyalty') || e.type?.includes('loyalty'));
  }
  if (selectedFilter.value === 'invoices') {
    return props.events.filter(e => e.event?.includes('facturas') || e.type?.includes('facturas'));
  }
  return props.events;
});

function getEventColor(event) {
  if (event.event?.includes('loyalty') || event.type?.includes('loyalty')) return 'accent';
  if (event.event?.includes('facturas') || event.type?.includes('facturas')) return 'secondary';
  if (event.event?.includes('legacy') || event.type?.includes('legacy')) return 'warning';
  return 'info';
}

function getEventIcon(event) {
  if (event.event?.includes('loyalty')) return 'mdi-star-shooting';
  if (event.event?.includes('facturas')) return 'mdi-receipt-text-check';
  if (event.event?.includes('connection')) return 'mdi-connection';
  return 'mdi-broadcast';
}
</script>

<template>
  <v-card class="pa-5 md3-card-expressive h-100 d-flex flex-column" elevation="0">
    <!-- Header -->
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-avatar
          size="40"
          class="mr-3"
          style="background: linear-gradient(135deg, rgba(6, 182, 212, 0.2) 0%, rgba(99, 102, 241, 0.3) 100%); border: 1px solid rgba(6, 182, 212, 0.3);"
        >
          <v-icon color="info" size="22">mdi-console-line</v-icon>
        </v-avatar>
        <div>
          <div class="text-h6 font-weight-bold text-white" style="letter-spacing: -0.01em;">
            Consola de Eventos
          </div>
          <div class="text-caption text-grey">
            Streaming WebSocket &middot; Kong Proxy Ingress
          </div>
        </div>
      </div>

      <div class="d-flex align-center ga-1">
        <v-chip size="small" color="primary" variant="tonal" class="font-weight-bold">
          {{ events.length }} eventos
        </v-chip>
        <v-btn
          icon="mdi-trash-can-outline"
          size="small"
          variant="text"
          color="grey"
          :disabled="events.length === 0"
          @click="emit('clear-events')"
        ></v-btn>
      </div>
    </div>

    <!-- Filter Pills -->
    <div class="d-flex ga-1 mb-3">
      <v-chip
        size="x-small"
        :color="selectedFilter === 'all' ? 'primary' : 'surface-variant'"
        :variant="selectedFilter === 'all' ? 'flat' : 'tonal'"
        class="font-weight-bold cursor-pointer"
        @click="selectedFilter = 'all'"
      >
        Todos ({{ events.length }})
      </v-chip>
      <v-chip
        size="x-small"
        :color="selectedFilter === 'loyalty' ? 'accent' : 'surface-variant'"
        :variant="selectedFilter === 'loyalty' ? 'flat' : 'tonal'"
        class="font-weight-bold cursor-pointer"
        @click="selectedFilter = 'loyalty'"
      >
        Fidelidad
      </v-chip>
      <v-chip
        size="x-small"
        :color="selectedFilter === 'invoices' ? 'secondary' : 'surface-variant'"
        :variant="selectedFilter === 'invoices' ? 'flat' : 'tonal'"
        class="font-weight-bold cursor-pointer"
        @click="selectedFilter = 'invoices'"
      >
        Facturas
      </v-chip>
    </div>

    <!-- Empty State -->
    <div
      v-if="events.length === 0"
      class="pa-8 text-center rounded-2xl my-auto d-flex flex-column align-center justify-center"
      style="background: rgba(20, 29, 51, 0.4); border: 1px dashed rgba(255, 255, 255, 0.1);"
    >
      <v-avatar size="48" color="surface-variant" class="mb-3">
        <v-icon size="26" color="info">mdi-broadcast</v-icon>
      </v-avatar>
      <div class="text-body-2 font-weight-bold text-grey-lighten-1">Esperando eventos en vivo...</div>
      <div class="text-caption text-grey mt-1">
        Los mensajes publicados en RabbitMQ se empujarán instantáneamente a través de Kong Gateway.
      </div>
    </div>

    <!-- Events Feed -->
    <div v-else class="d-flex flex-column ga-2 flex-grow-1" style="max-height: 520px; overflow-y: auto;">
      <div
        v-for="evt in filteredEvents"
        :key="evt.id"
        class="pa-3 rounded-xl md3-item-interactive event-card-enter"
        style="background: rgba(20, 29, 51, 0.6);"
      >
        <!-- Top Row: Icon + Title + Time -->
        <div class="d-flex justify-space-between align-center mb-1">
          <div class="d-flex align-center ga-2">
            <v-avatar :color="getEventColor(evt)" size="26" variant="tonal">
              <v-icon size="15">{{ getEventIcon(evt) }}</v-icon>
            </v-avatar>
            <span class="font-weight-bold text-body-2 text-white">{{ evt.title || evt.event }}</span>
          </div>
          <span class="text-caption text-grey mono-font" style="font-size: 11px;">{{ evt.receivedAt }}</span>
        </div>

        <!-- Event Body -->
        <div class="text-caption text-grey-lighten-2 my-2 pl-8">
          <div v-if="evt.puntos_obtenidos !== undefined" class="d-flex align-center ga-1">
            <span>Acreditados:</span>
            <v-chip size="x-small" color="accent" variant="flat" class="font-weight-bold">
              +{{ evt.puntos_obtenidos }} pts
            </v-chip>
            <span>&middot; Total: <strong class="text-white">{{ evt.total_puntos }} pts</strong> (Nivel {{ evt.nivel }})</span>
          </div>
          <div v-else-if="evt.factura_id" class="d-flex align-center ga-1">
            <span>Factura <strong class="text-white">{{ evt.factura_id }}</strong>:</span>
            <v-chip size="x-small" color="secondary" variant="flat" class="font-weight-bold">
              ${{ evt.total?.toFixed(2) }}
            </v-chip>
            <span>en <strong class="text-info">{{ evt.almacen_asignado }}</strong></span>
          </div>
          <div v-else-if="evt.message">
            {{ evt.message }}
          </div>
        </div>

        <!-- CloudEvents Metadata Footer -->
        <div class="d-flex justify-space-between align-center mt-2 pt-2 border-t-grey-darken-4 text-caption text-grey">
          <span class="mono-font" style="font-size: 10px; color: #818cf8;">
            {{ evt.event || evt.type }}
          </span>
          <span v-if="evt.correlation_id" class="mono-font" style="font-size: 10px;">
            corr: {{ evt.correlation_id.substring(0, 16) }}...
          </span>
        </div>
      </div>
    </div>
  </v-card>
</template>
