import { defineStore } from 'pinia';
import axios from 'axios';

export const useInvoiceStore = defineStore('invoice', {
  state: () => ({
    invoices: [],
    isLoading: false,
    error: null,
  }),

  actions: {
    async fetchInvoices() {
      this.isLoading = true;
      this.error = null;

      try {
        const res = await axios.get('http://localhost:8000/api/v1/go/invoices');
        this.invoices = res.data.invoices || [];
      } catch (err) {
        this.error = err.message;
      } finally {
        this.isLoading = false;
      }
    },

    addFromEvent(eventData) {
      // Check if invoice already exists
      const exists = this.invoices.some((inv) => inv.factura_id === eventData.factura_id);
      if (!exists) {
        const newInvoice = {
          factura_id: eventData.factura_id,
          pedido_id: eventData.pedido_id,
          cliente_id: eventData.cliente_id,
          total: eventData.total,
          almacen_asignado: eventData.almacen_asignado,
          estado: eventData.estado || 'EMITIDA',
          correlation_id: eventData.correlation_id,
          fecha: eventData.timestamp || new Date().toISOString(),
        };
        this.invoices.unshift(newInvoice);
      }
    },
  },
});
