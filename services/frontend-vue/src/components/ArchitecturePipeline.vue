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
const isFastAPIActive = computed(() => props.lastEventType?.includes('loyalty') || props.lastEventType?.includes('connection') || props.lastEventType?.includes('facturas'));
const isKongActive = computed(() => props.isConnected && Boolean(props.lastEventType));
</script>

<template>
  <v-card class="pa-5 mb-6 md3-card-expressive" elevation="0">
    <!-- Header: Title & Connection Status -->
    <div class="d-flex flex-wrap align-center justify-space-between ga-3 mb-4">
      <div class="d-flex align-center">
        <v-avatar
          size="38"
          class="mr-3"
          style="background: linear-gradient(135deg, rgba(99, 102, 241, 0.2) 0%, rgba(139, 92, 246, 0.3) 100%); border: 1px solid rgba(99, 102, 241, 0.3);"
        >
          <v-icon color="primary" size="22">mdi-transit-connection-variant</v-icon>
        </v-avatar>
        <div>
          <div class="text-subtitle-1 font-weight-bold text-white" style="letter-spacing: -0.01em;">
            Pipeline EDA en Tiempo Real
          </div>
          <div class="text-caption text-grey">
            Flujo reactivo distribuido: Compras ➔ AMQP Broker ➔ Motores Concurrentes ➔ WebSocket Push
          </div>
        </div>
      </div>

      <!-- Live WebSocket Indicator Pill -->
      <v-chip
        size="small"
        :color="isConnected ? 'success' : 'error'"
        variant="tonal"
        class="px-3 py-2 font-weight-bold"
        style="border: 1px solid rgba(16, 185, 129, 0.3);"
      >
        <span
          class="mr-2"
          :style="{
            display: 'inline-block',
            width: '8px',
            height: '8px',
            borderRadius: '50%',
            backgroundColor: isConnected ? '#10b981' : '#f43f5e',
            boxShadow: isConnected ? '0 0 8px #10b981' : 'none'
          }"
        ></span>
        {{ isConnected ? 'Kong Gateway Conectado (/ws)' : 'Desconectado' }}
      </v-chip>
    </div>

    <!-- Expressive Architecture Nodes Grid -->
    <v-row dense align="center" justify="space-between" class="py-1">
      <!-- 1. Monolito Legacy -->
      <v-col cols="12" sm="6" md="2">
        <div
          class="pa-3 rounded-xl transition-swing text-center md3-item-interactive"
          :class="{
            'node-pulsing': isLegacyActive,
            'md3-item-active': isLegacyActive
          }"
          :style="{
            background: isLegacyActive ? 'linear-gradient(135deg, rgba(99, 102, 241, 0.25) 0%, rgba(79, 70, 229, 0.35) 100%)' : 'rgba(20, 29, 51, 0.6)',
            border: isLegacyActive ? '1px solid #6366f1' : '1px solid rgba(255, 255, 255, 0.06)'
          }"
        >
          <v-avatar size="36" color="primary" variant="tonal" class="mb-1">
            <v-icon size="20">mdi-storefront-outline</v-icon>
          </v-avatar>
          <div class="text-caption font-weight-bold text-white">1. Monolito Legacy</div>
          <v-chip size="x-small" color="primary" variant="outlined" class="mt-1">Fire & Forget</v-chip>
          <div class="text-caption text-grey-lighten-1 mt-1 mono-font" style="font-size: 9px;">Python 3.12</div>
        </div>
      </v-col>

      <!-- Flow Connector Arrow -->
      <v-col cols="12" sm="1" class="d-none d-md-flex justify-center">
        <v-icon :color="isBrokerActive ? 'primary' : 'grey-darken-2'" size="22">
          {{ isBrokerActive ? 'mdi-motion-play-outline' : 'mdi-chevron-right' }}
        </v-icon>
      </v-col>

      <!-- 2. RabbitMQ Message Broker -->
      <v-col cols="12" sm="6" md="2">
        <div
          class="pa-3 rounded-xl transition-swing text-center md3-item-interactive"
          :class="{
            'node-pulsing': isBrokerActive,
            'md3-item-active': isBrokerActive
          }"
          :style="{
            background: isBrokerActive ? 'linear-gradient(135deg, rgba(249, 115, 22, 0.25) 0%, rgba(234, 88, 12, 0.35) 100%)' : 'rgba(20, 29, 51, 0.6)',
            border: isBrokerActive ? '1px solid #f97316' : '1px solid rgba(255, 255, 255, 0.06)'
          }"
        >
          <v-avatar size="36" color="warning" variant="tonal" class="mb-1">
            <v-icon size="20">mdi-message-processing-outline</v-icon>
          </v-avatar>
          <div class="text-caption font-weight-bold text-white">2. RabbitMQ Broker</div>
          <v-chip size="x-small" color="warning" variant="outlined" class="mt-1">Topic + DLX</v-chip>
          <div class="text-caption text-grey-lighten-1 mt-1 mono-font" style="font-size: 9px;">CloudEvents 1.0</div>
        </div>
      </v-col>

      <!-- Flow Connector Arrow -->
      <v-col cols="12" sm="1" class="d-none d-md-flex justify-center">
        <v-icon :color="isGoActive || isFastAPIActive ? 'secondary' : 'grey-darken-2'" size="22">
          {{ isGoActive || isFastAPIActive ? 'mdi-motion-play-outline' : 'mdi-chevron-right' }}
        </v-icon>
      </v-col>

      <!-- 3. Concurrent Engine (Go & FastAPI) -->
      <v-col cols="12" sm="6" md="3">
        <div class="d-flex flex-column ga-2">
          <!-- 3a: Go Billing Engine -->
          <div
            class="pa-2 px-3 rounded-xl transition-swing d-flex align-center justify-space-between md3-item-interactive"
            :class="{
              'node-pulsing': isGoActive,
              'md3-item-active': isGoActive
            }"
            :style="{
              background: isGoActive ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.25) 0%, rgba(5, 150, 105, 0.35) 100%)' : 'rgba(20, 29, 51, 0.6)',
              border: isGoActive ? '1px solid #10b981' : '1px solid rgba(255, 255, 255, 0.06)'
            }"
          >
            <div class="d-flex align-center">
              <v-avatar size="28" color="secondary" variant="tonal" class="mr-2">
                <v-icon size="16">mdi-language-go</v-icon>
              </v-avatar>
              <div class="text-left">
                <div class="text-caption font-weight-bold text-white">3a. Go Billing Engine</div>
                <div class="text-caption text-grey" style="font-size: 10px;">Stock Atómico & 19% IVA</div>
              </div>
            </div>
            <v-chip size="x-small" color="secondary" variant="flat">Manual ACK</v-chip>
          </div>

          <!-- 3b: FastAPI Loyalty & WS -->
          <div
            class="pa-2 px-3 rounded-xl transition-swing d-flex align-center justify-space-between md3-item-interactive"
            :class="{
              'node-pulsing': isFastAPIActive,
              'md3-item-active': isFastAPIActive
            }"
            :style="{
              background: isFastAPIActive ? 'linear-gradient(135deg, rgba(139, 92, 246, 0.25) 0%, rgba(124, 58, 237, 0.35) 100%)' : 'rgba(20, 29, 51, 0.6)',
              border: isFastAPIActive ? '1px solid #8b5cf6' : '1px solid rgba(255, 255, 255, 0.06)'
            }"
          >
            <div class="d-flex align-center">
              <v-avatar size="28" color="accent" variant="tonal" class="mr-2">
                <v-icon size="16">mdi-lightning-bolt</v-icon>
              </v-avatar>
              <div class="text-left">
                <div class="text-caption font-weight-bold text-white">3b. FastAPI Orchestrator</div>
                <div class="text-caption text-grey" style="font-size: 10px;">Puntos & Servidor WS</div>
              </div>
            </div>
            <v-chip size="x-small" color="accent" variant="flat">aio-pika</v-chip>
          </div>
        </div>
      </v-col>

      <!-- Flow Connector Arrow -->
      <v-col cols="12" sm="1" class="d-none d-md-flex justify-center">
        <v-icon :color="isKongActive ? 'info' : 'grey-darken-2'" size="22">
          {{ isKongActive ? 'mdi-motion-play-outline' : 'mdi-chevron-right' }}
        </v-icon>
      </v-col>

      <!-- 4. Kong API Gateway & Push Client -->
      <v-col cols="12" sm="6" md="2">
        <div
          class="pa-3 rounded-xl transition-swing text-center md3-item-interactive"
          :class="{
            'node-pulsing': isKongActive,
            'md3-item-active': isKongActive
          }"
          :style="{
            background: isKongActive ? 'linear-gradient(135deg, rgba(6, 182, 212, 0.25) 0%, rgba(14, 116, 144, 0.35) 100%)' : 'rgba(20, 29, 51, 0.6)',
            border: isKongActive ? '1px solid #06b6d4' : '1px solid rgba(255, 255, 255, 0.06)'
          }"
        >
          <v-avatar size="36" color="info" variant="tonal" class="mb-1">
            <v-icon size="20">mdi-shield-check</v-icon>
          </v-avatar>
          <div class="text-caption font-weight-bold text-white">4. Kong Gateway</div>
          <v-chip size="x-small" color="info" variant="outlined" class="mt-1">WS Streaming</v-chip>
          <div class="text-caption text-grey-lighten-1 mt-1 mono-font" style="font-size: 9px;">Port 8000 Proxy</div>
        </div>
      </v-col>
    </v-row>
  </v-card>
</template>
