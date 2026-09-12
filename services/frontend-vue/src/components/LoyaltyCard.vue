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
      return 'blue-grey-lighten-1';
    default:
      return 'indigo-lighten-1';
  }
});

const tierGradient = computed(() => {
  switch (props.tier) {
    case 'Gold':
      return 'linear-gradient(135deg, rgba(245, 158, 11, 0.25) 0%, rgba(180, 83, 9, 0.35) 100%)';
    case 'Silver':
      return 'linear-gradient(135deg, rgba(148, 163, 184, 0.22) 0%, rgba(71, 85, 105, 0.35) 100%)';
    default:
      return 'linear-gradient(135deg, rgba(99, 102, 241, 0.2) 0%, rgba(139, 92, 246, 0.28) 100%)';
  }
});

const tierBorder = computed(() => {
  switch (props.tier) {
    case 'Gold':
      return 'rgba(245, 158, 11, 0.4)';
    case 'Silver':
      return 'rgba(148, 163, 184, 0.4)';
    default:
      return 'rgba(99, 102, 241, 0.4)';
  }
});
</script>

<template>
  <v-card class="pa-5 md3-card-expressive h-100 d-flex flex-column" elevation="0">
    <!-- Header -->
    <div class="d-flex align-center justify-space-between mb-3">
      <div class="d-flex align-center">
        <v-avatar
          size="40"
          class="mr-3"
          style="background: linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(236, 72, 153, 0.3) 100%); border: 1px solid rgba(139, 92, 246, 0.3);"
        >
          <v-icon color="accent" size="22">mdi-star-shooting-outline</v-icon>
        </v-avatar>
        <div>
          <div class="text-h6 font-weight-bold text-white" style="letter-spacing: -0.01em;">
            Programa de Fidelidad
          </div>
          <div class="text-caption text-grey">
            FastAPI &middot; Orquestación y Fidelidad
          </div>
        </div>
      </div>

      <!-- Expressive Tier Pill Badge -->
      <v-chip
        :color="tierBadgeColor"
        variant="elevated"
        size="small"
        class="font-weight-black px-3"
        style="box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);"
      >
        <v-icon start size="15">mdi-crown</v-icon>
        NIVEL {{ tier.toUpperCase() }}
      </v-chip>
    </div>

    <div class="text-caption text-grey-lighten-1 mb-3">
      Cliente: <span class="mono-font font-weight-bold text-white">{{ clientId }}</span>
    </div>

    <!-- Dynamic Tier Hero Banner -->
    <div
      class="pa-4 rounded-2xl mb-4 transition-swing"
      :style="{
        background: tierGradient,
        border: `1px solid ${tierBorder}`,
        boxShadow: '0 8px 24px -6px rgba(0, 0, 0, 0.4)'
      }"
    >
      <div class="d-flex justify-space-between align-center">
        <div>
          <div class="text-caption text-grey-lighten-2 font-weight-medium">Puntos Acumulados</div>
          <div class="text-h3 font-weight-black text-white mono-font my-1">
            {{ totalPoints }}
          </div>
          <div class="text-caption text-grey-lighten-2">
            Equivalente a ${{ (totalPoints * 0.1).toFixed(2) }} en descuentos
          </div>
        </div>

        <v-avatar size="56" color="surface" variant="flat" style="border: 2px solid rgba(255, 255, 255, 0.15);">
          <v-icon size="32" :color="tierBadgeColor">
            {{ tier === 'Gold' ? 'mdi-trophy-variant' : tier === 'Silver' ? 'mdi-shield-star' : 'mdi-medal' }}
          </v-icon>
        </v-avatar>
      </div>
    </div>

    <!-- Tier Milestones Progress -->
    <div class="mb-4">
      <div class="d-flex justify-space-between text-caption mb-1">
        <span class="font-weight-bold text-white">Progreso a Nivel {{ nextTierTarget.nextTier }}</span>
        <span class="font-weight-bold mono-font text-accent">{{ totalPoints }} / {{ nextTierTarget.target }} pts</span>
      </div>
      <v-progress-linear
        :model-value="nextTierTarget.progress"
        color="accent"
        height="10"
        rounded="pill"
        style="box-shadow: 0 0 10px rgba(139, 92, 246, 0.3);"
      ></v-progress-linear>

      <div v-if="nextTierTarget.remaining > 0" class="text-caption text-grey mt-1 text-right">
        Faltan {{ nextTierTarget.remaining }} puntos (1 pt por cada $10 gastados)
      </div>
      <div v-else class="text-caption text-success font-weight-bold mt-1 text-right">
        ¡Has alcanzado la categoría máxima!
      </div>

      <!-- Tier Ladder Badges -->
      <div class="d-flex justify-space-between align-center mt-3 pt-2 border-t-grey-darken-4">
        <v-chip size="x-small" :color="tier === 'Standard' ? 'primary' : 'surface-variant'" variant="tonal">
          Standard (0+)
        </v-chip>
        <v-icon size="14" color="grey-darken-2">mdi-chevron-right</v-icon>
        <v-chip size="x-small" :color="tier === 'Silver' ? 'blue-grey-lighten-1' : 'surface-variant'" variant="tonal">
          Silver (50+)
        </v-chip>
        <v-icon size="14" color="grey-darken-2">mdi-chevron-right</v-icon>
        <v-chip size="x-small" :color="tier === 'Gold' ? 'amber-darken-1' : 'surface-variant'" variant="tonal">
          Gold (150+)
        </v-chip>
      </div>
    </div>

    <!-- Recent History -->
    <div class="mt-auto">
      <div class="d-flex justify-space-between align-center mb-2">
        <span class="text-caption text-grey font-weight-bold text-uppercase" style="letter-spacing: 0.05em;">
          Acreditaciones Recientes
        </span>
        <v-chip size="x-small" color="accent" variant="outlined">{{ history.length }} registros</v-chip>
      </div>

      <div v-if="history.length === 0" class="pa-3 text-center rounded-xl bg-surface-variant text-caption text-grey font-italic">
        No hay registros previos. Realiza una compra para acumular puntos.
      </div>
      <div v-else class="d-flex flex-column ga-1" style="max-height: 120px; overflow-y: auto;">
        <div
          v-for="(item, idx) in history.slice(0, 3)"
          :key="idx"
          class="pa-2 px-3 rounded-xl md3-item-interactive d-flex justify-space-between align-center text-caption"
          style="background: rgba(20, 29, 51, 0.5);"
        >
          <div class="d-flex align-center">
            <v-icon size="16" color="accent" class="mr-2">mdi-plus-circle-outline</v-icon>
            <span class="mono-font text-white">{{ item.pedido_id }}</span>
          </div>
          <v-chip size="x-small" color="accent" variant="flat" class="font-weight-bold">
            +{{ item.puntos_ganados }} pts
          </v-chip>
        </div>
      </div>
    </div>
  </v-card>
</template>
