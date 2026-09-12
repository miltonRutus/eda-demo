import { defineStore } from 'pinia';
import axios from 'axios';

export const useLoyaltyStore = defineStore('loyalty', {
  state: () => ({
    clientId: 'cli-442',
    totalPoints: 0,
    tier: 'Standard',
    history: [],
    isLoading: false,
    error: null,
  }),

  getters: {
    nextTierTarget: (state) => {
      if (state.totalPoints >= 150) return { target: 150, nextTier: 'Gold', remaining: 0, progress: 100 };
      if (state.totalPoints >= 50) {
        return {
          target: 150,
          nextTier: 'Gold',
          remaining: 150 - state.totalPoints,
          progress: Math.min(100, Math.round(((state.totalPoints - 50) / 100) * 100)),
        };
      }
      return {
        target: 50,
        nextTier: 'Silver',
        remaining: 50 - state.totalPoints,
        progress: Math.min(100, Math.round((state.totalPoints / 50) * 100)),
      };
    },
    tierColor: (state) => {
      switch (state.tier) {
        case 'Gold':
          return '#f59e0b'; // Amber
        case 'Silver':
          return '#94a3b8'; // Slate Silver
        default:
          return '#6366f1'; // Indigo Standard
      }
    },
  },

  actions: {
    async fetchPoints(clientId = 'cli-442') {
      this.clientId = clientId;
      this.isLoading = true;
      this.error = null;

      try {
        const res = await axios.get(`http://localhost:8000/api/v1/fastapi/customers/${encodeURIComponent(clientId)}/points`);
        this.totalPoints = res.data.total_puntos || 0;
        this.tier = res.data.nivel || 'Standard';
        this.history = res.data.historial_reciente || [];
      } catch (err) {
        // If 404 (no orders yet), default to 0 points gracefully
        if (err.response?.status === 404) {
          this.totalPoints = 0;
          this.tier = 'Standard';
          this.history = [];
        } else {
          this.error = err.message;
        }
      } finally {
        this.isLoading = false;
      }
    },

    updateFromEvent(eventData) {
      if (eventData.total_puntos !== undefined) {
        this.totalPoints = eventData.total_puntos;
      }
      if (eventData.nivel) {
        this.tier = eventData.nivel;
      }
      if (eventData.puntos_obtenidos) {
        this.history.unshift({
          pedido_id: eventData.pedido_id,
          puntos_obtenidos: eventData.puntos_obtenidos,
          fecha: eventData.timestamp || new Date().toISOString(),
          correlation_id: eventData.correlation_id,
        });
      }
    },
  },
});
