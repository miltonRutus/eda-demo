<script setup>
import { computed } from 'vue';

const props = defineProps({
  lastEventType: {
    type: String,
    default: null,
  },
  isConnected: {
    type: Boolean,
    default: false,
  },
});

const isLegacyActive = computed(() => props.lastEventType === 'order_submitted' || props.lastEventType?.includes('legacy'));
const isBrokerActive = computed(() => Boolean(props.lastEventType));
const isGoActive = computed(() => props.lastEventType?.includes('facturas'));
const isFastAPIActive = computed(() => props.lastEventType?.includes('loyalty') || props.lastEventType?.includes('connection'));
</script>

<template>
  <v-card class="pa-4 mb-6 rounded-lg" color="#0f172a" elevation="2">
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-icon color="primary" class="mr-2">mdi-transit-connection-variant</v-icon>
        <span class="text-subtitle-1 font-weight-bold">Pipeline EDA en Tiempo Real</span>
      </div>
      <v-chip size="small" :color="isConnected ? 'success' : 'error'" variant="flat">
        <v-icon start size="small">{{ isConnected ? 'mdi-wifi-check' : 'mdi-wifi-off' }}</v-icon>
        {{ isConnected ? 'Kong Proxy Conectado' : 'Desconectado' }}
      </v-chip>
    </div>

    <!-- Pipeline Visual Nodes -->
    <v-row dense align="center" justify="space-between" class="text-center py-2">
      <!-- Node 1: Monolito Legacy -->
      <v-col cols="12" sm="2">
        <v-sheet
          class="pa-3 rounded-lg transition-swing"
          :color="isLegacyActive ? 'primary-darken-1' : '#1e293b'"
          :class="{ 'node-pulsing': isLegacyActive }"
        >
          <v-icon size="28" :color="isLegacyActive ? 'white' : 'grey-lighten-1'">mdi-storefront-outline</v-icon>
          <div class="text-caption font-weight-bold mt-1">1. Monolito Legacy</div>
          <div class="text-caption text-grey-lighten-2" style="font-size: 10px;">Python Fire & Forget</div>
        </v-sheet>
      </v-col>

      <!-- Arrow -->
      <v-col cols="12" sm="1" class="d-none d-sm-flex justify-center">
        <v-icon color="grey">mdi-arrow-right-bold</v-icon>
      </v-col>

      <!-- Node 2: RabbitMQ -->
      <v-col cols="12" sm="2">
        <v-sheet
          class="pa-3 rounded-lg transition-swing"
          :color="isBrokerActive ? '#f97316' : '#1e293b'"
          :class="{ 'node-pulsing': isBrokerActive }"
        >
          <v-icon size="28" :color="isBrokerActive ? 'white' : 'grey-lighten-1'">mdi-message-processing-outline</v-icon>
          <div class="text-caption font-weight-bold mt-1">2. RabbitMQ Broker</div>
          <div class="text-caption text-grey-lighten-2" style="font-size: 10px;">Topic & CloudEvents</div>
        </v-sheet>
      </v-col>

      <!-- Arrow -->
      <v-col cols="12" sm="1" class="d-none d-sm-flex justify-center">
        <v-icon color="grey">mdi-arrow-right-bold</v-icon>
      </v-col>

      <!-- Node 3: Concurrent Consumers (Go & FastAPI) -->
      <v-col cols="12" sm="3">
        <div class="d-flex flex-column ga-2">
          <!-- Go Microservice -->
          <v-sheet
            class="pa-2 rounded-lg transition-swing"
            :color="isGoActive ? 'secondary' : '#1e293b'"
            :class="{ 'node-pulsing': isGoActive }"
          >
            <div class="d-flex align-center justify-center">
              <v-icon size="20" class="mr-1">mdi-language-go</v-icon>
              <span class="text-caption font-weight-bold">3a. Go Billing Engine</span>
            </div>
            <div class="text-caption text-grey-lighten-2" style="font-size: 10px;">Inventario & Factura Fiscal</div>
          </v-sheet>

          <!-- FastAPI Microservice -->
          <v-sheet
            class="pa-2 rounded-lg transition-swing"
            :color="isFastAPIActive ? 'accent' : '#1e293b'"
            :class="{ 'node-pulsing': isFastAPIActive }"
          >
            <div class="d-flex align-center justify-center">
              <v-icon size="20" class="mr-1">mdi-lightning-bolt</v-icon>
              <span class="text-caption font-weight-bold">3b. FastAPI Orchestrator</span>
            </div>
            <div class="text-caption text-grey-lighten-2" style="font-size: 10px;">Fidelidad & Servidor WS</div>
          </v-sheet>
        </div>
      </v-col>

      <!-- Arrow -->
      <v-col cols="12" sm="1" class="d-none d-sm-flex justify-center">
        <v-icon color="grey">mdi-arrow-right-bold</v-icon>
      </v-col>

      <!-- Node 4: Kong Gateway & Frontend UI -->
      <v-col cols="12" sm="2">
        <v-sheet class="pa-3 rounded-lg" color="#0284c7">
          <v-icon size="28" color="white">mdi-shield-check</v-icon>
          <div class="text-caption font-weight-bold mt-1">4. Kong Gateway</div>
          <div class="text-caption text-grey-lighten-2" style="font-size: 10px;">Auth Offload & WS Push</div>
        </v-sheet>
      </v-col>
    </v-row>
  </v-card>
</template>
