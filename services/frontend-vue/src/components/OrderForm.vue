<script setup>
import { ref } from 'vue';

const props = defineProps({
  catalog: {
    type: Array,
    required: true,
  },
  selectedQuantities: {
    type: Object,
    required: true,
  },
  totalAmount: {
    type: Number,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update-quantity', 'submit-order']);

const clientId = ref('cli-442');

function changeQuantity(sku, delta) {
  const current = props.selectedQuantities[sku] || 0;
  const updated = Math.max(0, current + delta);
  emit('update-quantity', { sku, quantity: updated });
}

function handleOrderSubmit() {
  emit('submit-order', clientId.value);
}
</script>

<template>
  <v-card class="pa-5 md3-card-expressive h-100 d-flex flex-column" elevation="0">
    <!-- Header -->
    <div class="d-flex align-center justify-space-between mb-4">
      <div class="d-flex align-center">
        <v-avatar
          size="40"
          class="mr-3"
          style="background: linear-gradient(135deg, rgba(99, 102, 241, 0.2) 0%, rgba(139, 92, 246, 0.3) 100%); border: 1px solid rgba(99, 102, 241, 0.3);"
        >
          <v-icon color="primary" size="22">mdi-cart-arrow-right</v-icon>
        </v-avatar>
        <div>
          <div class="text-h6 font-weight-bold text-white" style="letter-spacing: -0.01em;">
            Simulador de Compra
          </div>
          <div class="text-caption text-grey">
            Monolito Legacy (Python) &middot; Emisor CloudEvents
          </div>
        </div>
      </div>
      <v-chip size="small" color="primary" variant="tonal">
        REST Ingress
      </v-chip>
    </div>

    <!-- Client ID Input -->
    <div class="pa-3 mb-4 rounded-xl" style="background: rgba(20, 29, 51, 0.5); border: 1px solid rgba(255, 255, 255, 0.05);">
      <div class="text-caption text-grey mb-1 font-weight-medium">Identificador de Usuario (Multi-tenant):</div>
      <v-text-field
        v-model="clientId"
        prepend-inner-icon="mdi-account-circle-outline"
        variant="solo-filled"
        density="compact"
        flat
        rounded="lg"
        bg-color="surface-variant"
        class="mono-font"
        hide-details
      ></v-text-field>
    </div>

    <!-- Catalog Header -->
    <div class="d-flex justify-space-between align-center mb-2">
      <span class="text-caption text-grey font-weight-bold text-uppercase" style="letter-spacing: 0.05em;">
        Catálogo de Productos
      </span>
      <span class="text-caption text-grey">Selecciona unidades</span>
    </div>

    <!-- Product Catalog List -->
    <div class="d-flex flex-column ga-2 flex-grow-1 mb-3" style="max-height: 260px; overflow-y: auto;">
      <div
        v-for="product in catalog"
        :key="product.sku"
        class="pa-3 rounded-xl md3-item-interactive d-flex align-center justify-space-between"
        :class="{ 'md3-item-active': (selectedQuantities[product.sku] || 0) > 0 }"
        :style="{
          background: (selectedQuantities[product.sku] || 0) > 0 ? 'rgba(99, 102, 241, 0.1)' : 'rgba(20, 29, 51, 0.5)'
        }"
      >
        <div class="d-flex align-center mr-2">
          <v-avatar
            size="40"
            class="mr-3"
            :color="(selectedQuantities[product.sku] || 0) > 0 ? 'primary' : 'surface-variant'"
            variant="tonal"
          >
            <v-icon size="22" :color="(selectedQuantities[product.sku] || 0) > 0 ? 'primary' : 'grey-lighten-1'">
              {{ product.icon }}
            </v-icon>
          </v-avatar>
          <div>
            <div class="font-weight-bold text-body-2 text-white">{{ product.name }}</div>
            <div class="text-caption text-grey">{{ product.description }}</div>
            <div class="text-caption font-weight-bold text-primary mt-1">
              ${{ product.price.toFixed(2) }}
            </div>
          </div>
        </div>

        <!-- Pill Stepper -->
        <div
          class="d-flex align-center rounded-pill px-2 py-1"
          style="background: rgba(15, 23, 42, 0.8); border: 1px solid rgba(255, 255, 255, 0.08);"
        >
          <v-btn
            icon="mdi-minus"
            size="x-small"
            variant="text"
            color="grey-lighten-1"
            :disabled="isLoading || (selectedQuantities[product.sku] || 0) <= 0"
            @click="changeQuantity(product.sku, -1)"
          ></v-btn>
          <span
            class="text-body-2 font-weight-black mx-2 mono-font"
            :class="(selectedQuantities[product.sku] || 0) > 0 ? 'text-primary' : 'text-grey'"
            style="min-width: 18px; text-align: center;"
          >
            {{ selectedQuantities[product.sku] || 0 }}
          </span>
          <v-btn
            icon="mdi-plus"
            size="x-small"
            variant="text"
            color="primary"
            :disabled="isLoading"
            @click="changeQuantity(product.sku, 1)"
          ></v-btn>
        </div>
      </div>
    </div>

    <!-- Sticky Checkout Bottom Bar -->
    <div
      class="pa-4 rounded-2xl mt-auto"
      style="background: rgba(15, 23, 42, 0.9); border: 1px solid rgba(255, 255, 255, 0.08);"
    >
      <div class="d-flex justify-space-between align-center mb-3">
        <div>
          <div class="text-caption text-grey">Total de la Orden</div>
          <div class="text-caption text-grey-lighten-2" style="font-size: 11px;">(Sincrónico + Asíncrono)</div>
        </div>
        <div class="text-h4 font-weight-black text-white mono-font">
          ${{ totalAmount.toFixed(2) }}
        </div>
      </div>

      <v-btn
        block
        size="large"
        class="submit-order-btn"
        :loading="isLoading"
        :disabled="isLoading || totalAmount <= 0"
        @click="handleOrderSubmit"
      >
        <v-icon start size="20">mdi-lightning-bolt</v-icon>
        Disparar Compra a Monolito
      </v-btn>
    </div>
  </v-card>
</template>
