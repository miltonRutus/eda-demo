import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach } from 'vitest';
import { useOrderStore } from '../src/stores/orderStore';
import { useLoyaltyStore } from '../src/stores/loyaltyStore';
import { useWebSocketStore } from '../src/stores/webSocketStore';

describe('Pinia Domain Stores', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('calculates order store cart totals correctly', () => {
    const orderStore = useOrderStore();

    // Default PROD-A has 1 unit ($149.99)
    expect(orderStore.totalAmount).toBe(149.99);

    // Set 2 units of PROD-B ($50.00 each)
    orderStore.setQuantity('PROD-B', 2);
    expect(orderStore.totalAmount).toBe(249.99);

    // Set PROD-A to 0
    orderStore.setQuantity('PROD-A', 0);
    expect(orderStore.totalAmount).toBe(100.00);
  });

  it('updates loyalty store reactively from WebSocket events', () => {
    const loyaltyStore = useLoyaltyStore();

    expect(loyaltyStore.totalPoints).toBe(0);
    expect(loyaltyStore.tier).toBe('Standard');

    // Simulate loyalty event received from WebSocket
    loyaltyStore.updateFromEvent({
      event: 'loyalty.points.updated',
      pedido_id: 'ORD-TEST-123',
      puntos_obtenidos: 55,
      total_puntos: 55,
      nivel: 'Silver',
      correlation_id: 'corr-123',
    });

    expect(loyaltyStore.totalPoints).toBe(55);
    expect(loyaltyStore.tier).toBe('Silver');
    expect(loyaltyStore.history.length).toBe(1);
    expect(loyaltyStore.history[0].pedido_id).toBe('ORD-TEST-123');
  });

  it('manages events queue in webSocketStore', () => {
    const wsStore = useWebSocketStore();

    expect(wsStore.events.length).toBe(0);

    // Add manual event
    wsStore.events.push({ id: '1', event: 'test.event' });
    expect(wsStore.events.length).toBe(1);

    wsStore.clearEvents();
    expect(wsStore.events.length).toBe(0);
  });
});
