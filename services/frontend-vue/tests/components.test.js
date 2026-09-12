import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

import OrderForm from '../src/components/OrderForm.vue';
import LoyaltyCard from '../src/components/LoyaltyCard.vue';

const vuetify = createVuetify({ components, directives });

describe('Dumb Components', () => {
  it('renders OrderForm and emits submit-order', async () => {
    const catalog = [
      { sku: 'PROD-A', name: 'Teclado RGB', price: 100, icon: 'mdi-keyboard', description: 'desc' },
    ];
    const selectedQuantities = { 'PROD-A': 2 };

    const wrapper = mount(OrderForm, {
      global: { plugins: [vuetify] },
      props: {
        catalog,
        selectedQuantities,
        totalAmount: 200,
        isLoading: false,
      },
    });

    expect(wrapper.text()).toContain('Teclado RGB');
    expect(wrapper.text()).toContain('$200.00');

    // Click submit button specifically
    const btn = wrapper.find('.submit-order-btn');
    await btn.trigger('click');

    expect(wrapper.emitted('submit-order')).toBeTruthy();
    expect(wrapper.emitted('submit-order')[0]).toEqual(['cli-442']);
  });

  it('renders LoyaltyCard with correct tier and points', () => {
    const wrapper = mount(LoyaltyCard, {
      global: { plugins: [vuetify] },
      props: {
        clientId: 'cli-442',
        totalPoints: 120,
        tier: 'Silver',
        nextTierTarget: { target: 150, nextTier: 'Gold', remaining: 30, progress: 70 },
        history: [],
        isLoading: false,
      },
    });

    expect(wrapper.text()).toContain('cli-442');
    expect(wrapper.text()).toContain('120');
    expect(wrapper.text()).toContain('NIVEL SILVER');
    expect(wrapper.text()).toContain('Faltan 30 puntos');
  });
});
