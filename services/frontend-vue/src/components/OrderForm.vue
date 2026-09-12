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
  <v-card class="pa-4 rounded-lg" color="surface" elevation="2">
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-icon color="primary" class="mr-2">mdi-cart-arrow-right</v-icon>
        <span class="text-h6 font-weight-bold">Simulador de Compra (Legacy)</span>
      </div>
    </div>

    <!-- Client ID Selection -->
    <v-text-field
      v-model="clientId"
      label="ID de Cliente Simulado"
      prepend-inner-icon="mdi-account"
      variant="outlined"
      density="compact"
      class="mb-2"
      hide-details
    ></v-text-field>

    <!-- Product Catalog List -->
    <div class="text-caption text-grey mb-2">Selecciona productos del catálogo:</div>
    <v-list class="bg-transparent pa-0">
      <v-list-item
        v-for="product in catalog"
        :key="product.sku"
        class="pa-2 mb-2 rounded-lg bg-surface-variant"
      >
        <template #prepend>
          <v-avatar color="primary-darken-1" size="36" class="mr-2">
            <v-icon size="20">{{ product.icon }}</v-icon>
          </v-avatar>
        </template>

        <v-list-item-title class="font-weight-medium text-body-2">
          {{ product.name }}
        </v-list-item-title>
        <v-list-item-subtitle class="text-caption text-grey">
          ${{ product.price.toFixed(2) }} — {{ product.description }}
        </v-list-item-subtitle>

        <template #append>
          <div class="d-flex align-center ga-1">
            <v-btn
              icon="mdi-minus"
              size="x-small"
              variant="tonal"
              :disabled="isLoading || (selectedQuantities[product.sku] || 0) <= 0"
              @click="changeQuantity(product.sku, -1)"
            ></v-btn>
            <span class="text-body-2 font-weight-bold mx-2" style="min-width: 16px; text-align: center;">
              {{ selectedQuantities[product.sku] || 0 }}
            </span>
            <v-btn
              icon="mdi-plus"
              size="x-small"
              variant="tonal"
              color="primary"
              :disabled="isLoading"
              @click="changeQuantity(product.sku, 1)"
            ></v-btn>
          </div>
        </template>
      </v-list-item>
    </v-list>

    <!-- Total Summary & Submit -->
    <v-divider class="my-3"></v-divider>
    <div class="d-flex align-center justify-space-between mb-3">
      <span class="text-subtitle-2 text-grey">Total a Facturar:</span>
      <span class="text-h6 font-weight-bold text-success">${{ totalAmount.toFixed(2) }}</span>
    </div>

    <v-btn
      block
      color="primary"
      size="large"
      rounded="lg"
      prepend-icon="mdi-send-check"
      class="submit-order-btn"
      :loading="isLoading"
      :disabled="totalAmount <= 0 || isLoading"
      @click="handleOrderSubmit"
    >
      Disparar Compra a Monolito
    </v-btn>
  </v-card>
</template>
