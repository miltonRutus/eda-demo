import 'vuetify/styles';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

const customDarkTheme = {
  dark: true,
  colors: {
    background: '#070a13',
    surface: '#0f172a',
    'surface-variant': '#141d33',
    'surface-bright': '#1e293b',
    'surface-container': '#1a2642',
    'surface-container-high': '#223254',
    'surface-container-highest': '#2a3e66',
    primary: '#6366f1',
    'primary-darken-1': '#4f46e5',
    'primary-container': '#232252',
    secondary: '#10b981',
    'secondary-container': '#064e3b',
    accent: '#8b5cf6',
    tertiary: '#ec4899',
    error: '#f43f5e',
    info: '#06b6d4',
    success: '#10b981',
    warning: '#f59e0b',
  },
};

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'customDarkTheme',
    themes: {
      customDarkTheme,
    },
  },
});
