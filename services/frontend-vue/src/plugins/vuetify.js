import 'vuetify/styles';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

const customDarkTheme = {
  dark: true,
  colors: {
    background: '#0a0e17',
    surface: '#121826',
    'surface-variant': '#1e293b',
    primary: '#6366f1',
    'primary-darken-1': '#4f46e5',
    secondary: '#10b981',
    accent: '#8b5cf6',
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
