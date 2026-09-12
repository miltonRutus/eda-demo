# ==============================================================================
# Development Stage (Target: dev)
# - Servidor Vite con recarga en caliente (Hot-Reload)
# - Suite de pruebas unitarias (Vitest)
# ==============================================================================
FROM node:20-alpine AS dev

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]

# ==============================================================================
# Build Stage: Generación de assets estáticos
# ==============================================================================
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .
RUN npm run build

# ==============================================================================
# Production Stage (Target: prod)
# - Inmutable, mínimo y seguro con Nginx Alpine
# ==============================================================================
FROM nginx:alpine AS prod

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
