import { defineStore } from 'pinia';

export const useWebSocketStore = defineStore('websocket', {
  state: () => ({
    socket: null,
    isConnected: false,
    clientId: 'cli-442',
    events: [],
    reconnectTimer: null,
    lastActiveTimestamp: null,
  }),

  actions: {
    connect(clientId = 'cli-442') {
      this.clientId = clientId;

      // Close existing socket if different
      if (this.socket) {
        try {
          this.socket.close();
        } catch {
          // ignore
        }
      }

      const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      // In dev browser, connect to Kong Gateway port 8000
      const wsHost = window.location.hostname || 'localhost';
      const wsUrl = `${wsProtocol}//${wsHost}:8000/ws/${encodeURIComponent(clientId)}`;

      console.log(`[WebSocket] Connecting to ${wsUrl}...`);

      try {
        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
          console.log('[WebSocket] Connection open.');
          this.isConnected = true;
          this.lastActiveTimestamp = new Date().toISOString();
        };

        this.socket.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            console.log('[WebSocket] Message received:', data);
            this.lastActiveTimestamp = new Date().toISOString();

            // Add unique ID and local arrival time for UI
            const enrichedEvent = {
              id: `evt-${Date.now()}-${Math.random().toString(36).substring(2, 6)}`,
              receivedAt: new Date().toLocaleTimeString(),
              ...data,
            };

            this.events.unshift(enrichedEvent);

            // Limit memory to last 50 events
            if (this.events.length > 50) {
              this.events.pop();
            }
          } catch (err) {
            console.error('[WebSocket] Failed to parse message JSON:', err);
          }
        };

        this.socket.onclose = () => {
          console.log('[WebSocket] Disconnected. Will retry in 3s...');
          this.isConnected = false;
          this.scheduleReconnect();
        };

        this.socket.onerror = (err) => {
          console.warn('[WebSocket] Socket encountered error:', err);
          this.isConnected = false;
        };
      } catch (err) {
        console.error('[WebSocket] Connect exception:', err);
        this.isConnected = false;
        this.scheduleReconnect();
      }
    },

    scheduleReconnect() {
      if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
      this.reconnectTimer = setTimeout(() => {
        if (!this.isConnected) {
          this.connect(this.clientId);
        }
      }, 3000);
    },

    disconnect() {
      if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
      this.isConnected = false;
    },

    clearEvents() {
      this.events = [];
    },
  },
});
