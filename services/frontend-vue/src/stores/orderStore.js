import { defineStore } from 'pinia';
import axios from 'axios';

export const useOrderStore = defineStore('order', {
  state: () => ({
    catalog: [
      { sku: 'PROD-A', name: 'Teclado Mecánico RGB Custom', price: 149.99, icon: 'mdi-keyboard', description: 'Switches ópticos lineales y keycaps PBT' },
      { sku: 'PROD-B', name: 'Mouse Inalámbrico Ultraligero', price: 50.00, icon: 'mdi-mouse', description: 'Sensor 26K DPI y peso de 49g' },
      { sku: 'PROD-C', name: 'Monitor 4K IPS 144Hz', price: 89.99, icon: 'mdi-monitor', description: '1ms GTG, HDR600 y calibración de fábrica' },
      { sku: 'PROD-D', name: 'Auriculares Hi-Fi Wireless', price: 120.00, icon: 'mdi-headphones', description: 'Cancelación activa de ruido híbrida' },
    ],
    selectedQuantities: {
      'PROD-A': 1,
      'PROD-B': 0,
      'PROD-C': 0,
      'PROD-D': 0,
    },
    isLoading: false,
    lastOrder: null,
    error: null,
  }),

  getters: {
    cartItems: (state) => {
      return state.catalog
        .filter((item) => state.selectedQuantities[item.sku] > 0)
        .map((item) => ({
          sku: item.sku,
          name: item.name,
          precio: item.price,
          cantidad: state.selectedQuantities[item.sku],
          subtotal: Number((item.price * state.selectedQuantities[item.sku]).toFixed(2)),
        }));
    },
    totalAmount: (state) => {
      return Number(
        state.catalog
          .reduce((sum, item) => sum + item.price * (state.selectedQuantities[item.sku] || 0), 0)
          .toFixed(2)
      );
    },
  },

  actions: {
    setQuantity(sku, qty) {
      this.selectedQuantities[sku] = Math.max(0, qty);
    },
    async submitOrder(clientId = 'cli-442') {
      if (this.cartItems.length === 0) {
        throw new Error('Selecciona al menos un producto para realizar el pedido.');
      }

      this.isLoading = true;
      this.error = null;

      try {
        const payload = {
          cliente_id: clientId,
          items: this.cartItems.map((item) => ({
            sku: item.sku,
            cantidad: item.cantidad,
            precio: item.precio,
          })),
        };

        const correlationId = `web-${Date.now()}-${Math.random().toString(36).substring(2, 7)}`;

        const response = await axios.post('http://localhost:8000/api/v1/legacy/orders', payload, {
          headers: {
            'X-Correlation-ID': correlationId,
          },
        });

        this.lastOrder = response.data;
        return response.data;
      } catch (err) {
        this.error = err.response?.data?.detail || err.message || 'Error al enviar pedido a Legacy.';
        throw err;
      } finally {
        this.isLoading = false;
      }
    },
  },
});
