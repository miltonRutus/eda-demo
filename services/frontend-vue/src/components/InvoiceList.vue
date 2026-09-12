<script setup>
defineProps({
  invoices: {
    type: Array,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});
</script>

<template>
  <v-card class="pa-4 rounded-lg h-100" color="surface" elevation="2">
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-icon color="secondary" class="mr-2">mdi-receipt-text-check-outline</v-icon>
        <span class="text-h6 font-weight-bold">Facturación Concurrente (Go)</span>
      </div>
      <v-chip size="small" color="secondary" variant="flat">
        {{ invoices.length }} emitidas
      </v-chip>
    </div>

    <div class="text-caption text-grey mb-2">
      Facturas fiscales generadas de manera asíncrona tras reservar stock:
    </div>

    <div v-if="invoices.length === 0" class="text-caption text-grey font-italic text-center py-6">
      Aún no se han generado facturas. Dispara un pedido para iniciar la concurrencia de Go.
    </div>

    <div v-else class="d-flex flex-column ga-2" style="max-height: 240px; overflow-y: auto;">
      <div
        v-for="inv in invoices"
        :key="inv.factura_id"
        class="pa-3 rounded-lg bg-surface-variant event-card-enter"
      >
        <div class="d-flex justify-space-between align-center mb-1">
          <div class="d-flex align-center ga-1">
            <span class="mono-font font-weight-bold text-body-2">{{ inv.factura_id }}</span>
            <v-chip size="x-small" color="secondary" variant="tonal">{{ inv.estado }}</v-chip>
          </div>
          <span class="text-subtitle-2 font-weight-bold text-secondary">${{ inv.total?.toFixed(2) }}</span>
        </div>

        <div class="d-flex justify-space-between text-caption text-grey">
          <span>Pedido: <strong class="mono-font text-grey-lighten-2">{{ inv.pedido_id }}</strong></span>
          <span>Almacén: <v-chip size="x-small" color="info" variant="outlined">{{ inv.almacen_asignado }}</v-chip></span>
        </div>
      </div>
    </div>
  </v-card>
</template>
