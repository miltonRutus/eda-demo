<script setup>
import { computed } from 'vue';

const props = defineProps({
  clientId: {
    type: String,
    required: true,
  },
  totalPoints: {
    type: Number,
    required: true,
  },
  tier: {
    type: String,
    required: true,
  },
  nextTierTarget: {
    type: Object,
    required: true,
  },
  history: {
    type: Array,
    default: () => [],
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const tierBadgeColor = computed(() => {
  switch (props.tier) {
    case 'Gold':
      return 'amber-darken-1';
    case 'Silver':
      return 'blue-grey-lighten-2';
    default:
      return 'indigo-lighten-1';
  }
});
</script>

<template>
  <v-card class="pa-4 rounded-lg h-100" color="surface" elevation="2">
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-icon color="accent" class="mr-2">mdi-star-shooting-outline</v-icon>
        <span class="text-h6 font-weight-bold">Programa de Fidelidad (FastAPI)</span>
      </div>
      <v-chip :color="tierBadgeColor" variant="flat" size="small" class="font-weight-bold">
        <v-icon start size="16">mdi-crown</v-icon>
        NIVEL {{ tier.toUpperCase() }}
      </v-chip>
    </div>

    <div class="text-caption text-grey mb-1">Cliente: {{ clientId }}</div>

    <!-- Points Metric Banner -->
    <div class="pa-4 rounded-lg bg-surface-variant d-flex align-center justify-space-between my-2">
      <div>
        <div class="text-caption text-grey">Puntos Acumulados</div>
        <div class="text-h4 font-weight-black text-primary">{{ totalPoints }}</div>
      </div>
      <v-avatar color="accent" size="48" variant="tonal">
        <v-icon size="28">mdi-medal-outline</v-icon>
      </v-avatar>
    </div>

    <!-- Tier Progress -->
    <div class="my-3">
      <div class="d-flex justify-space-between text-caption mb-1">
        <span>Progreso a Nivel {{ nextTierTarget.nextTier }}</span>
        <span class="font-weight-medium text-grey">{{ totalPoints }} / {{ nextTierTarget.target }} pts</span>
      </div>
      <v-progress-linear
        :model-value="nextTierTarget.progress"
        color="accent"
        height="8"
        rounded
      ></v-progress-linear>
      <div v-if="nextTierTarget.remaining > 0" class="text-caption text-grey mt-1 text-right">
        Faltan {{ nextTierTarget.remaining }} puntos (1 pt por cada $10 gastados)
      </div>
      <div v-else class="text-caption text-success mt-1 text-right">
        ¡Has alcanzado la categoría máxima!
      </div>
    </div>

    <!-- Recent History -->
    <v-divider class="my-2"></v-divider>
    <div class="text-caption text-grey mb-2">Acreditaciones Recientes:</div>
    <div v-if="history.length === 0" class="text-caption text-grey font-italic py-2">
      No hay registros previos. Realiza una compra para acumular puntos.
    </div>
    <div v-else class="d-flex flex-column ga-1" style="max-height: 120px; overflow-y: auto;">
      <div
        v-for="(item, idx) in history.slice(0, 3)"
        :key="idx"
        class="d-flex justify-space-between align-center pa-2 rounded bg-surface-variant text-caption"
      >
        <span class="mono-font">{{ item.pedido_id }}</span>
        <v-chip size="x-small" color="success" variant="tonal">+{{ item.puntos_obtenidos }} pts</v-chip>
      </div>
    </div>
  </v-card>
</template>
