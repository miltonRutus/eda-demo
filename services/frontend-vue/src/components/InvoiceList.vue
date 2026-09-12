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
  <v-card class="pa-5 md3-card-expressive d-flex flex-column" elevation="0">
    <!-- Header -->
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-avatar
          size="40"
          class="mr-3"
          style="background: linear-gradient(135deg, rgba(16, 185, 129, 0.2) 0%, rgba(6, 182, 212, 0.3) 100%); border: 1px solid rgba(16, 185, 129, 0.3);"
        >
          <v-icon color="secondary" size="22">mdi-receipt-text-check-outline</v-icon>
        </v-avatar>
        <div>
          <div class="text-h6 font-weight-bold text-white" style="letter-spacing: -0.01em;">
            Facturación Concurrente
          </div>
          <div class="text-caption text-grey">
            API 2: Golang &middot; Goroutines & Impuestos Fiscales
          </div>
        </div>
      </div>

      <v-chip size="small" color="secondary" variant="tonal" class="font-weight-bold">
        {{ invoices.length }} emitidas
      </v-chip>
    </div>

    <div class="text-caption text-grey mb-3">
      Facturas fiscales generadas de manera asíncrona tras reservar stock:
    </div>

    <!-- Empty State -->
    <div
      v-if="invoices.length === 0"
      class="pa-6 text-center rounded-2xl my-auto d-flex flex-column align-center justify-center"
      style="background: rgba(20, 29, 51, 0.4); border: 1px dashed rgba(255, 255, 255, 0.1);"
    >
      <v-avatar size="44" color="surface-variant" class="mb-2">
        <v-icon size="24" color="grey">mdi-receipt-text-outline</v-icon>
      </v-avatar>
      <div class="text-body-2 font-weight-bold text-grey-lighten-1">Sin facturas emitidas</div>
      <div class="text-caption text-grey mt-1">
        Dispara un pedido para iniciar el procesamiento concurrente de Go.
      </div>
    </div>

    <!-- Invoices List -->
    <div v-else class="d-flex flex-column ga-2 flex-grow-1" style="max-height: 250px; overflow-y: auto;">
      <div
        v-for="inv in invoices"
        :key="inv.factura_id"
        class="pa-3 rounded-xl md3-item-interactive event-card-enter"
        style="background: rgba(20, 29, 51, 0.6);"
      >
        <div class="d-flex justify-space-between align-center mb-1">
          <div class="d-flex align-center ga-2">
            <v-icon size="16" color="secondary">mdi-check-decagram</v-icon>
            <span class="mono-font font-weight-black text-body-2 text-white">{{ inv.factura_id }}</span>
            <v-chip size="x-small" color="secondary" variant="tonal" class="font-weight-bold">{{ inv.estado }}</v-chip>
          </div>
          <span class="text-subtitle-1 font-weight-black text-secondary mono-font">
            ${{ inv.total?.toFixed(2) }}
          </span>
        </div>

        <div class="d-flex justify-space-between align-center text-caption text-grey mt-2 pt-2 border-t-grey-darken-4">
          <span>Pedido: <strong class="mono-font text-white">{{ inv.pedido_id }}</strong></span>
          <v-chip size="x-small" color="info" variant="outlined" class="font-weight-bold">
            <v-icon start size="12">mdi-map-marker-outline</v-icon>
            {{ inv.almacen_asignado }}
          </v-chip>
        </div>
      </div>
    </div>
  </v-card>
</template>
