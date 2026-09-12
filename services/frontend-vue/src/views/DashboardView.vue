<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import { useOrderStore } from '../stores/orderStore';
import { useWebSocketStore } from '../stores/webSocketStore';
import { useLoyaltyStore } from '../stores/loyaltyStore';
import { useInvoiceStore } from '../stores/invoiceStore';

import ArchitecturePipeline from '../components/ArchitecturePipeline.vue';
import OrderForm from '../components/OrderForm.vue';
import LoyaltyCard from '../components/LoyaltyCard.vue';
import InvoiceList from '../components/InvoiceList.vue';
import EventStream from '../components/EventStream.vue';

const orderStore = useOrderStore();
const wsStore = useWebSocketStore();
const loyaltyStore = useLoyaltyStore();
const invoiceStore = useInvoiceStore();

const activeClientId = ref('cli-442');
const lastEventKey = ref(null);
const snackbar = ref({ show: false, text: '', color: 'success' });

// Listen for incoming WebSocket events to update reactive domain stores
watch(
  () => wsStore.events[0],
  (newEvent) => {
    if (!newEvent) return;

    lastEventKey.value = newEvent.event || newEvent.type;

    if (newEvent.event === 'loyalty.points.updated') {
      loyaltyStore.updateFromEvent(newEvent);
    } else if (newEvent.event === 'facturacion.facturas.generada') {
      invoiceStore.addFromEvent(newEvent);
    }

    // Reset pulse animation after 3s
    setTimeout(() => {
      lastEventKey.value = null;
    }, 3000);
  }
);

async function handleOrderSubmit(clientId) {
  activeClientId.value = clientId;
  lastEventKey.value = 'order_submitted';

  try {
    const res = await orderStore.submitOrder(clientId);
    snackbar.value = {
      show: true,
      text: `¡Pedido ${res.pedido_id} enviado exitosamente al Monolito Legacy!`,
      color: 'success',
    };
  } catch (err) {
    snackbar.value = {
      show: true,
      text: err.message || 'Error al procesar el pedido.',
      color: 'error',
    };
  }
}

function handleQuantityUpdate({ sku, quantity }) {
  orderStore.setQuantity(sku, quantity);
}

onMounted(async () => {
  // 1. Establish WebSocket streaming through Kong Gateway
  wsStore.connect(activeClientId.value);

  // 2. Fetch initial domain state
  await loyaltyStore.fetchPoints(activeClientId.value);
  await invoiceStore.fetchInvoices();
});

onUnmounted(() => {
  wsStore.disconnect();
});
</script>

<template>
  <v-container fluid class="pa-4 pa-md-6">
    <!-- Architecture Flow Pipeline Indicator -->
    <ArchitecturePipeline
      :last-event-type="lastEventKey"
      :is-connected="wsStore.isConnected"
    />

    <v-row>
      <!-- Left Column: Order Simulator (Legacy Trigger) -->
      <v-col cols="12" md="4">
        <OrderForm
          :catalog="orderStore.catalog"
          :selected-quantities="orderStore.selectedQuantities"
          :total-amount="orderStore.totalAmount"
          :is-loading="orderStore.isLoading"
          @update-quantity="handleQuantityUpdate"
          @submit-order="handleOrderSubmit"
        />
      </v-col>

      <!-- Center Column: Microservices Reactive State (Loyalty & Invoices) -->
      <v-col cols="12" md="4">
        <div class="d-flex flex-column ga-4 h-100">
          <!-- FastAPI Loyalty Points -->
          <LoyaltyCard
            :client-id="activeClientId"
            :total-points="loyaltyStore.totalPoints"
            :tier="loyaltyStore.tier"
            :next-tier-target="loyaltyStore.nextTierTarget"
            :history="loyaltyStore.history"
            :is-loading="loyaltyStore.isLoading"
          />

          <!-- Go Fiscal Billing & Invoices -->
          <InvoiceList
            :invoices="invoiceStore.invoices"
            :is-loading="invoiceStore.isLoading"
          />
        </div>
      </v-col>

      <!-- Right Column: Real-Time Event Stream (Kong WebSockets) -->
      <v-col cols="12" md="4">
        <EventStream
          :events="wsStore.events"
          @clear-events="wsStore.clearEvents"
        />
      </v-col>
    </v-row>

    <!-- Feedback Snackbar (MD3 Expressive Pill) -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      timeout="4000"
      location="bottom right"
      rounded="pill"
      elevation="6"
      class="mb-4 mr-4"
    >
      <div class="d-flex align-center">
        <v-icon start size="18" class="mr-2">
          {{ snackbar.color === 'success' ? 'mdi-check-circle-outline' : 'mdi-alert-circle-outline' }}
        </v-icon>
        <span class="font-weight-medium text-body-2">{{ snackbar.text }}</span>
      </div>
      <template #actions>
        <v-btn variant="text" size="small" class="font-weight-bold" @click="snackbar.show = false">
          Cerrar
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>
